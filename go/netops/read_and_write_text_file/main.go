package main

import (
	"bufio"
	"os"
)

func main() {
	infile, _ := os.Open("./input.txt")
	scanner := bufio.NewScanner(infile)
	scanner.Split(bufio.ScanLines)

	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	infile.Close()

	outfile, _ := os.OpenFile("./output.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	writer := bufio.NewWriter(outfile)

	for _, line := range lines {
		_, _ = writer.WriteString(line + "\n")
	}
	writer.Flush()
	outfile.Close()
}
