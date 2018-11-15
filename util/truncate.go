package util

import (
	"fmt"
	"strings"
)

// the maxChars is optimized for good reading experience;
// not so long, and not so short.
const maxChars = 150

func Truncate(text string, suffix string) string {
	if len(text) < maxChars {
		return text
	}

	trim := text[0:maxChars]
	split := strings.Split(trim, " ")
	newText := strings.Join(split[0:len(split)-1], " ")
	return fmt.Sprintf("%s%s", newText, suffix)
}
