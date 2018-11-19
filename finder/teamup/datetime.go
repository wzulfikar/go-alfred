package teamup

import (
	"encoding/json"
	"strings"
	"time"
)

type DateTime struct {
	time.Time
}

// converts RFC3339 time to Go time.
func (t *DateTime) UnmarshalJSON(data []byte) error {
	var dt string
	if err := json.Unmarshal(data, &dt); err != nil {
		return err
	}

	if !strings.Contains(dt, "+") {
		dt += "+00:00"
	}

	parsedTime, err := time.Parse(time.RFC3339, dt)
	if err != nil {
		return err
	}
	t.Time = parsedTime

	return nil
}
