package filepath

import (
	"log"
	"path/filepath"
	"strings"
)

func GetFileNameWithoutExtension(filePath string) string {
	fileNameWithExtension := filepath.Base(filePath)
	fileNameWithoutExtension := strings.TrimSuffix(fileNameWithExtension, filepath.Ext(fileNameWithExtension))
	return fileNameWithoutExtension
}

func GetAbsoluteFilePath(filePath string) string {
	absolutePath, err := filepath.Abs(filePath)
	if err != nil {
		log.Fatal(err)
	}
	return absolutePath
}
