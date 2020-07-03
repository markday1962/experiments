package main

// A simple app to save redis databases
import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

func main() {
	host := ""
	if len(os.Args) > 1 {
		host = os.Args[1]
	} else {
		host = "localhost"
	}
	backupRedis(host)
}

func backupRedis(host string) {

	// Create a new redis client
	rdb := redis.NewClient(&redis.Options{
		Addr: host + ":6379",
	})

	defer rdb.Close()

	// Test redis connection with PING
	pong, err := rdb.Ping().Result()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	msg := "Response to PING: " + pong
	fmt.Printf("%v\n", msg)

	// Get last save in unix time
	ls := rdb.LastSave()
	lst, err := ls.Result()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	msg = "Current LastSave value: " + strconv.FormatInt(lst, 10)
	fmt.Printf("%v\n", msg)

	// start background save
	resp := rdb.BgSave()
	fmt.Println(resp)
	time.Sleep(2 * time.Second)
	nls := rdb.LastSave()
	nlst, err := nls.Result()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for {
		if nlst > lst {
			msg = "New LastSave value: " + strconv.FormatInt(nlst, 10)
			fmt.Println(msg)
			fmt.Println("Redis BGSAVE has completed.")
			break
		}
		time.Sleep(30 * time.Second)
		nls = rdb.LastSave()
		nlst, err = nls.Result()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
