package main

import (
	"fmt"
	"strings"
	"time"
)

func diffTime(t1, t2 time.Time) string {
	diff := t1.Sub(t2)

	if diff.Hours() < 24 {
		return diff.String()
	}

	days := int(diff.Hours() / 24)
	hrs := int(diff.Hours()) % (24 * days)
	mins := int(diff.Minutes()) - (days * 24 * 60) - (hrs * 60)
	return fmt.Sprintf("%dd %dh %dm", days, hrs, mins)
}

func shorten(text string) string {
	if len(text) < 160 {
		return text
	}

	trim := text[0:160]
	split := strings.Split(trim, " ")
	newText := strings.Join(split[0:len(split)-1], " ")

	return fmt.Sprintf("%s...", newText)
}
