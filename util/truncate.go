package util

import (
	"fmt"
	"strings"
)

// the maxChars is optimized for good reading experience;
// not so long, and not so short.
const maxChars = 150

func TruncateWords(text string, suffix string, length ...int) string {
	var limit = maxChars
	if len(length) > 0 {
		limit = length[0]
	}

	if len(text) < limit {
		return text
	}

	trim := text[0:limit]
	split := strings.Split(trim, " ")
	newText := strings.Join(split[0:len(split)-1], " ")

	return fmt.Sprintf("%s%s", newText, suffix)
}

func TruncateString(text string, suffix string, length ...int) string {
	var limit = maxChars
	if len(length) > 0 {
		limit = length[0]
	}

	if len(text) < limit {
		return text
	}

	return fmt.Sprintf("%s%s", text[0:limit], suffix)
}
