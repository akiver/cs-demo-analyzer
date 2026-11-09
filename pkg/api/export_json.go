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
	options := []jsonv2.Options{json.DefaultOptionsV1(), jsonv2.WithMarshalers(positionsMarshalers(match))}
	if !minify {
		options = append(options, jsontext.WithIndent("  "))
	}

	// The MatchJSON value is marshalled rather than the match so that Match.MarshalJSON, which marshals the whole
	// document in memory, isn't called.
	return jsonv2.MarshalWrite(w, match.toJSON(), options...)
}

// positionsMarshalers returns marshalers of the match's position slices that stream the positions a round at a time,
// including the ones spooled to disk during the analysis which aren't in the slices, see Match.eachPlayerPositions().
func positionsMarshalers(match *Match) *jsonv2.Marshalers {
	return jsonv2.JoinMarshalers(
		jsonv2.MarshalToFunc(func(encoder *jsontext.Encoder, _ []*PlayerPosition) error {
			return marshalPositions(encoder, match.eachPlayerPositions)
		}),
		jsonv2.MarshalToFunc(func(encoder *jsontext.Encoder, _ []*GrenadePosition) error {
			return marshalPositions(encoder, match.eachGrenadePositions)
		}),
		jsonv2.MarshalToFunc(func(encoder *jsontext.Encoder, _ []*InfernoPosition) error {
			return marshalPositions(encoder, match.eachInfernoPositions)
		}),
		jsonv2.MarshalToFunc(func(encoder *jsontext.Encoder, _ []*ChickenPosition) error {
			return marshalPositions(encoder, match.eachChickenPositions)
		}),
		jsonv2.MarshalToFunc(func(encoder *jsontext.Encoder, _ []*HostagePosition) error {
			return marshalPositions(encoder, match.eachHostagePositions)
		}),
	)
}

// marshalPositions writes the positions yielded by each as a JSON array.
func marshalPositions[T any](encoder *jsontext.Encoder, each func(fn func(positions []T) error) error) error {
	if err := encoder.WriteToken(jsontext.BeginArray); err != nil {
		return err
	}

	err := each(func(positions []T) error {
		for _, position := range positions {
			if err := jsonv2.MarshalEncode(encoder, position); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	return encoder.WriteToken(jsontext.EndArray)
}
