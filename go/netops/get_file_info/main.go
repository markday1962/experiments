package main

import (
	"fmt"
	"os"
)

/* Gets the stats info off a file */
func main() {
	file, _ := os.Stat("main.go")

	fmt.Println("File Name: ", file.Name())
	fmt.Println("Size in Bytes: ", file.Size())
	fmt.Println("Last Modified: ", file.ModTime())
	fmt.Println("Is a directory: ", file.IsDir())

	fmt.Printf("Permissions 9-bit: %s\n", file.Mode())
	fmt.Printf("Permissions 3-digit: %o\n", file.Mode())
	fmt.Printf("Permissions 4-digit: %04o\n", file.Mode())
}
