package main

import "os"

/* Create a file, write to it and close */
func main() {
	data := "The example demonstrates the opening a file and writing to it"
	file, _ := os.OpenFile("file.txt",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	file.Write([]byte(data))
	file.Close()
}
