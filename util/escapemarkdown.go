package util

import (
	"strings"
)

func EscapeMarkdown(mdtext string) string {
	escaped := strings.Replace(mdtext, "*", "", -1)
	escaped = strings.Replace(escaped, "_", "", -1)
	escaped = strings.Replace(escaped, "~", "", -1)
	escaped = strings.Replace(escaped, "`", "", -1)

	return escaped
}
