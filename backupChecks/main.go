package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

const (
	//arangoDir = "/mnt/data/arango/"
	arangoBackupDir = "/mnt/data/arango/backup/"
	arangoSharedDir = "/mnt/data/arango/shared/"
	arangoLibDir = "/mnt/data/arango/arangodb3/"
	elasticsearchDir = "/mnt/data/elasticsearch"
	elasticsearchDataDir = "/mnt/data/elasticsearch/data"
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

	bd := ""
	dates := []string{}
	for _, key := range resp.Contents {
		s := strings.Split(*key.Key, "/")
		if len(s) > 2 {
			if s[1] == getDate() {
				bd = s[1]
				break
			} else {
				if ! contains(dates, s[1]) {
					dates = append(dates, s[1])
				}
			}
		}
	}
	if len(bd) != 8 {
		bd = dates[len(dates)-1]
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

func contains(s []string, e string) bool {
	// Helper function to see if a slice contains a value
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func getDate() string {
	// Helper function to return current date in a format used by the S3 prekey
	dt := time.Now().AddDate(0, 0, 0)
	return dt.Format("20060102")
}

func main() {
	_, err := getLiveCipher()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	downloadArangoBackup()

	os.Exit(0)
}
