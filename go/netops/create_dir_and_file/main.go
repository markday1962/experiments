package main

import (
	"fmt"
	"os"
	"path/filepath"
)

/* Create a directory and file */
func main() {
	dpath := "./newdir"
	fname := "file.txt"

	_ = os.MkdirAll(dpath, 0777)
	fpath := filepath.Join(dpath, fname)
	file, err := os.Create(fpath)
	if err != nil {
		fmt.Printf("Error caught %v", err)
	}

	file.Close()
}
