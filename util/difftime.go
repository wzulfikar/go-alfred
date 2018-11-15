package util

import (
	"fmt"
	"time"
)

func DiffTime(t1, t2 time.Time) string {
	diff := t1.Sub(t2)

	if diff.Hours() < 24 {
		return diff.String()
	}

	days := int(diff.Hours() / 24)
	hrs := int(diff.Hours()) % (24 * days)
	mins := int(diff.Minutes()) - (days * 24 * 60) - (hrs * 60)
	return fmt.Sprintf("%dd %dh %dm", days, hrs, mins)
}
