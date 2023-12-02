package api

import (
	"encoding/json"
	"errors"
	"os"
)

func buildOutputFilePath(match *Match, outputPath string) (string, error) {
	if outputPath == "" {
		return match.DemoFilePath + ".json", nil
	}

	stat, err := os.Stat(outputPath)
	if err != nil {
		return "", errors.New("invalid output provided, make sure the path exists and you have write access")
	}

	if stat.IsDir() {
		return outputPath + string(os.PathSeparator) + match.DemoFileName + ".json", nil
	}

	return outputPath, nil
}

func exportMatchToJSON(match *Match, outputPath string, minify bool) error {
	var err error
	outputFilePath, err := buildOutputFilePath(match, outputPath)
	if err != nil {
		return err
	}

	var jsonString []byte
	if minify {
		jsonString, err = json.Marshal(match)
	} else {
		jsonString, err = json.MarshalIndent(match, "", "  ")
	}

	if err != nil {
		return err
	}

	err = os.WriteFile(outputFilePath, jsonString, os.ModePerm)

	return err
}
