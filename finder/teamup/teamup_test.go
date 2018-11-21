package teamup

import (
	"log"
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestTeamupFinder(t *testing.T) {
	finder := &TeamupFinder{
		ApiKey:      os.Getenv("TEAMUP_API_KEY"),
		CalendarKey: os.Getenv("TEAMUP_CALENDAR_KEY"),
	}

	if err := finder.Init(); err != nil {
		log.Fatal(err)
	}

	query := "event"
	result, err := finder.Find(query)
	if err != nil {
		t.Fatal(err)
	}

	spew.Dump(result)
}
