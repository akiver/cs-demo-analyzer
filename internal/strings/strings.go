package strings

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

func RemoveInvisibleChars(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsGraphic(r) {
			return r
		}

		return -1
	}, s)
}

func RemoveInvalidUTF8Sequences(s string) string {
	return strings.ToValidUTF8(s, "")
}

// Replaces invalid UTF-8 byte sequences from a string by their corresponding Unicode replacement character.
// For example: "🚩Test\xf6Test界" will become "🚩TestöTest界".
func ReplaceUTF8ByteSequences(s string) string {
	if utf8.ValidString(s) {
		return s
	}

	var b strings.Builder
	for i, r := range s {
		r, _ := utf8.DecodeRuneInString(string(r))
		if r == utf8.RuneError {
			b.WriteString(string(s[i]))
		} else {
			b.WriteRune(r)
		}
	}

	return b.String()
}
