package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

const (
	arangoDir = "/mnt/data/arango/"
	arangoBackupDir = "/mnt/data/arango/backup/"
	arangoSharedDir = "/mnt/data/arango/shared/"
	arangoLibDir = "/mnt/data/arango/arangodb3/"
	redisDir = "/mnt/data/redis/"
	redisDataDir = "/mnt/data/redis/data/"
	configDir = "/mnt/data/config/"
	dataBucket = "aistemos-data-backups/"
	configBucket = "aistemos-cloud-init/"
)

var region = "eu-west-2"
var sess, _ = session.NewSession(&aws.Config{Region: aws.String(region)})
var svc = s3.New(sess)

type CipherResponse struct {
	FEVersion   string `json:"Aistemos-FrontendService-Version"`
	Version     string `json:"Aistemos-Software-Version"`
	DataVersion int64  `json:"Aistemos-Data-Version"`
	Host        string `json:"Aistemos-Application-Id"`
}

//Function to display errors and exit.
func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func getLiveCipher() (string, error) {
	/* Return's the current live cipher host */
	resp, err := http.Get("http://app.cipher.ai/version")
	if err != nil {
		fmt.Printf("Error response too http request %v\n", err)
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf(err.Error())
		return "", err
	}
	x := body[1 : len(body)-1]
	c := new(CipherResponse)
	err = json.Unmarshal([]byte(x), c)
	if err != nil {
		fmt.Printf(err.Error())
		return "", err
	}
	return c.Host, nil
}

// Create download folders
func makeDir() error {
	os.Mkdir(redisDir, 0777)
	os.Mkdir(redisDataDir, 0777)
	os.Mkdir(arangoDir, 0777)
	os.Mkdir(arangoBackupDir, 0777)
	os.Mkdir(arangoSharedDir, 0777)
	os.Mkdir(arangoLibDir, 0777)
	os.Mkdir(configDir, 0777)
	return nil
}

func downloadRedisBackup(lc string) {

	df := "dump.rdb"
	pf := "redis" + "/" + lc + "-pfcache/"

	params := &s3.ListObjectsInput{
		Bucket: aws.String(dataBucket),
		Prefix: aws.String(pf),
	}

	resp, err := svc.ListObjects(params)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, key := range resp.Contents {
		if strings.Contains(*key.Key, df) {
			fmt.Println(strings.Contains(*key.Key, df))

			os.Chdir(redisDataDir)
			file, err := os.Create(df)
			if err != nil {
				exitErrorf("Unable to create file %v, %v", df, err)
			}

			defer file.Close()

			downloader := s3manager.NewDownloader(sess)
			_, err = downloader.Download(file,
				&s3.GetObjectInput{
					Bucket: aws.String(dataBucket),
					Key:    aws.String(pf + "/" + df),
				})
			if err != nil {
				exitErrorf("Unable to download item %q, %v", df, err)
			}
			//fmt.Println("Downloaded", file.Name, numBytes(replace _), "bytes")
		}
	}
}

func downloadArangoBackup() {
	/* Download the most current arango db backup files that has been uploaded to aws s3 */
	pf := "arango"

	params := &s3.ListObjectsInput{
		Bucket: aws.String(dataBucket),
		Prefix: aws.String(pf),
	}

	resp, err := svc.ListObjects(params)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	bd := getDateOffset(0)
	for _, key := range resp.Contents {
		s := strings.Split(*key.Key, "/")
		if s[1] == bd {
			break
		} else {
			for i := 1; i < 32; i++ {
				ofd := getDateOffset(int(^uint(int(i) - 1)))
				if s[1] == ofd {
					fmt.Printf("%v\n", ofd)
					bd = ofd
				}
			}
		}
	}

	for _, key := range resp.Contents {
		fs := strings.Split(*key.Key, "/")
		os.Chdir(arangoBackupDir)
		if len(fs) > 2 {
			fname := fs[len(fs)-1]
			file, err := os.Create(fname)
			if err != nil {
				exitErrorf("Unable to create file %v, %v", fname, err)
			}

			defer file.Close()

			downloader := s3manager.NewDownloader(sess)
			_, err = downloader.Download(file,
				&s3.GetObjectInput{
					Bucket: aws.String(dataBucket),
					Key:    aws.String(pf + "/" + bd + "/" + fname),
				})
			if err != nil {
				exitErrorf("Unable to download item %q, %v", fname, err)
			}
		}
	}

	// Helper file to show backup date
	os.Chdir(arangoBackupDir)
	_, err = os.Create(bd + ".txt")
	if err != nil {
		exitErrorf("Unable to create file %q, %v", bd + ".txt", err)
	}
}

func getDateOffset(of int) string {
	/* returns a string of the current date with offset applied */
	dt := time.Now().AddDate(0, 0, of)
	return dt.Format("20060102")
}

func startDockerServices() {
	/* Download compose file and start services */

	sess2, _ := session.NewSession(&aws.Config{Region: aws.String("eu-west-1")})

	pf := "backup-checker"
	df := "docker-compose.yml"

	os.Chdir(configDir)
	file, err := os.Create(df)
	if err != nil {
		exitErrorf("Unable to create file %q, %v", df, err)
	}

	defer file.Close()

	downloader := s3manager.NewDownloader(sess2)
	_, err = downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(configBucket),
			Key:    aws.String(pf + "/" + df),
		})
	if err != nil {
		exitErrorf("Unable to download item %q, %v", df, err)
	}

	cmd := exec.Command("docker-compose", "up")
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

func getArangoBackupScript(){
	/* Download helper script to restore arango db */

	sess2, _ := session.NewSession(&aws.Config{Region: aws.String("eu-west-1")})

	pf := "backup-checker"
	df := "arango-restore.sh"

	os.Chdir(arangoSharedDir)
	file, err := os.Create(df)
	if err != nil {
		exitErrorf("Unable to create file %q, %v", df, err)
	}

	defer file.Close()

	downloader := s3manager.NewDownloader(sess2)
	_, err = downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(configBucket),
			Key:    aws.String(pf + "/" + df),
		})
	if err != nil {
		exitErrorf("Unable to download item %q, %v", df, err)
	}
}


func main() {
	lc, err := getLiveCipher()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = makeDir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	downloadRedisBackup(lc)
	downloadArangoBackup()
	startDockerServices()
	getArangoBackupScript()

}
