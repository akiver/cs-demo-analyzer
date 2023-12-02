package csv

import (
	"encoding/csv"
	"log"
	"os"
)

func WriteLinesIntoCsvFile(csvFilePath string, lines [][]string) {
	file, err := os.Create(csvFilePath)
	if err != nil {
		log.Fatal("Error creating csv file", csvFilePath, err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	if err := writer.WriteAll(lines); err != nil {
		log.Fatal("Cannot write CSV file", csvFilePath, err)
	}
}
