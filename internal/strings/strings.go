package strings

import (
	"strings"
	"unicode"
)

func RemoveInvisibleChars(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsGraphic(r) {
			return r
		}

		return -1
	}, s)
}
