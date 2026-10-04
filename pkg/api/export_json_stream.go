package api

import (
	"bufio"
	"encoding/json"
	"os"
	"reflect"
	"strings"
)

// exportMatchToJSONStreaming writes the minified Match JSON without materializing the whole document in memory.
//
// json.Marshal(match) serializes the entire document several times: Match.MarshalJSON marshals a MatchJSON wrapper
// into a buffer and returns a copy, then the outer json.Marshal compacts those bytes into a second buffer and
// returns another copy for os.WriteFile. When positions are included the document can exceed 1 GB, so several
// document-sized buffers are alive at the peak.
//
// Here the top-level object is written key by key in struct declaration order (the order encoding/json uses) and
// every slice field is streamed element by element through a buffered writer. Each element is still encoded by
// json.Marshal on the same structs, so the output bytes are identical to json.Marshal(match), only the
// document-sized buffers are gone.
//
// Keep in sync with MatchJSON in match.go: the trailing "gameModeStr" key mirrors the wrapper struct.
func exportMatchToJSONStreaming(match *Match, outputFilePath string) (err error) {
	file, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := file.Close(); err == nil {
			err = closeErr
		}
	}()

	// Write errors are sticky on a bufio.Writer, the final Flush reports the first one.
	writer := bufio.NewWriterSize(file, 1<<20)

	// Reflect over MatchAlias rather than Match so the custom Match.MarshalJSON is bypassed at the top level.
	// Nested types still go through their own MarshalJSON.
	value := reflect.ValueOf((*MatchAlias)(match)).Elem()
	valueType := value.Type()

	writer.WriteByte('{')
	isFirstKey := true
	for i := 0; i < valueType.NumField(); i++ {
		field := valueType.Field(i)
		if !field.IsExported() {
			continue
		}
		tag := field.Tag.Get("json")
		if tag == "-" {
			continue
		}
		name, _, _ := strings.Cut(tag, ",")
		if name == "" {
			name = field.Name
		}

		key, err := json.Marshal(name)
		if err != nil {
			return err
		}
		if !isFirstKey {
			writer.WriteByte(',')
		}
		isFirstKey = false
		writer.Write(key)
		writer.WriteByte(':')

		fieldValue := value.Field(i)
		if fieldValue.Kind() != reflect.Slice {
			encoded, err := json.Marshal(fieldValue.Interface())
			if err != nil {
				return err
			}
			writer.Write(encoded)
			continue
		}

		// encoding/json emits null for a nil slice.
		if fieldValue.IsNil() {
			writer.WriteString("null")
			continue
		}

		writer.WriteByte('[')
		for j := 0; j < fieldValue.Len(); j++ {
			if j > 0 {
				writer.WriteByte(',')
			}
			encoded, err := json.Marshal(fieldValue.Index(j).Interface())
			if err != nil {
				return err
			}
			if _, err := writer.Write(encoded); err != nil {
				return err
			}
		}
		writer.WriteByte(']')
	}

	gameModeStr, err := json.Marshal(match.GameModeStr().String())
	if err != nil {
		return err
	}
	if !isFirstKey {
		writer.WriteByte(',')
	}
	writer.WriteString(`"gameModeStr":`)
	writer.Write(gameModeStr)
	writer.WriteByte('}')

	return writer.Flush()
}
