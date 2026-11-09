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

// FileWriter writes a CSV file line by line, so that the lines don't have to be all built in memory before writing
// the file. Errors are fatal like WriteLinesIntoCsvFile.
type FileWriter struct {
	path   string
	file   *os.File
	writer *csv.Writer
}

func NewFileWriter(csvFilePath string) *FileWriter {
	file, err := os.Create(csvFilePath)
	if err != nil {
		log.Fatal("Error creating csv file", csvFilePath, err)
	}

	return &FileWriter{
		path:   csvFilePath,
		file:   file,
		writer: csv.NewWriter(file),
	}
}

func (w *FileWriter) Write(line []string) {
	if err := w.writer.Write(line); err != nil {
		log.Fatal("Cannot write CSV file", w.path, err)
	}
}

func (w *FileWriter) Close() {
	w.writer.Flush()
	if err := w.writer.Error(); err != nil {
		log.Fatal("Cannot write CSV file", w.path, err)
	}
	if err := w.file.Close(); err != nil {
		log.Fatal("Cannot write CSV file", w.path, err)
	}
}
