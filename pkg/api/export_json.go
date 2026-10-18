package api

import (
	"bufio"
	"encoding/json"
	"encoding/json/jsontext"
	jsonv2 "encoding/json/v2"
	"errors"
	"io"
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
	outputFilePath, err := buildOutputFilePath(match, outputPath)
	if err != nil {
		return err
	}

	// The JSON is streamed to a temporary file renamed once complete, so that an error (or a crash) half-way
	// doesn't leave a truncated file where the previous export was.
	temporaryFilePath := outputFilePath + ".tmp"

	file, err := os.OpenFile(temporaryFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}

	writer := bufio.NewWriterSize(file, 1<<20)

	err = writeMatchJSON(writer, match, minify)
	if err == nil {
		err = writer.Flush()
	}

	if closeErr := file.Close(); err == nil {
		err = closeErr
	}

	if err != nil {
		os.Remove(temporaryFilePath)

		return err
	}

	return os.Rename(temporaryFilePath, outputFilePath)
}

func writeMatchJSON(w io.Writer, match *Match, minify bool) error {
	options := []jsonv2.Options{json.DefaultOptionsV1()}
	if !minify {
		options = append(options, jsontext.WithIndent("  "))
	}

	// The MatchJSON value is marshalled rather than the match so that Match.MarshalJSON, which marshals the whole
	// document in memory, isn't called.
	return jsonv2.MarshalWrite(w, match.toJSON(), options...)
}
