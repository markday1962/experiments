package main

import (
	"fmt"
	"io/ioutil"
)

func main() {
	files, _ := ioutil.ReadDir("/Users/markday/")

	for _, file := range files {
		fmt.Println(file.Name(), file.ModTime())
	}
}
