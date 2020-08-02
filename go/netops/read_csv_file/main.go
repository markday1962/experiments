package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

type Row struct {
	ColA, ColB, ColC string
}

func main() {
	file, _ := os.Open("./input.csv")

	reader := csv.NewReader(file)
	rows := []Row{}

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}

		rows = append(rows, Row{
			ColA: row[0], ColB: row[1], ColC: row[2],
		})
	}

	for _, row := range rows {
		fmt.Printf("%s -- %s -- %s\n", row.ColA, row.ColB, row.ColC)
	}
}
