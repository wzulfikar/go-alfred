package teamup

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
	"github.com/wzulfikar/go-alfred/contracts"
	"github.com/wzulfikar/go-alfred/util"
	validator "gopkg.in/go-playground/validator.v9"
)

const finderName = "teamup_v1"
const baseUrl = "https://api.teamup.com"
const logoUrl = "http://logobucket.surge.sh/services/teamup-logo-sm.png"

type TeamupFinder struct {
	ApiKey      string
	CalendarKey string
}

// func NewFinder(apiKey, calendarKey string) *TeamupFinder {
//   validate := validator.New()
//   errs := validate.Struct(config)
//   if errs != nil {
//     log.Println("invalid configuration detected.", errs)
//     log.Fatal("to continue, please fix the configuration and restart the program")
//   }
// }

func (finder *TeamupFinder) Init() error {
	validate := validator.New()
	errs := validate.Struct(finder)
	if errs != nil {
		return errors.Wrap(errors.New(fmt.Sprintf("%s", errs)), finderName)
	}

	return nil
}

func (finder *TeamupFinder) FinderName() string {
	return finderName
}

func (finder *TeamupFinder) Find(query string) (*[]contracts.Result, error) {
	now := time.Now()

	// today -30d
	startDate := formatDate(now.Add(time.Hour * 24 * 30 * -1))

	// today +30d
	endDate := formatDate(now.Add(time.Hour * 24 * 30))

	endpoint := fmt.Sprintf("%s/%s/events?query=%s&startDate=%s&endDate=%s",
		baseUrl,
		finder.CalendarKey,
		url.QueryEscape(query),
		url.QueryEscape(startDate),
		url.QueryEscape(endDate))

	req, _ := http.NewRequest("GET", endpoint, nil)
	req.Header.Add("teamup-token", finder.ApiKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, finderName)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, finderName)
	}

	searchResult := &SearchResult{}
	if err := json.Unmarshal(body, searchResult); err != nil {
		log.Println("failed to unmarshal body:", string(body))
		return nil, errors.Wrap(err, finderName)
	}

	result := []contracts.Result{}
	for _, event := range searchResult.Events {
		who := event.Who
		if who == "" {
			who = "(Not set)"
		}

		where := event.Location
		if where == "" {
			where = "(Not set)"
		}

		var when string
		startDt := event.StartDt.Time
		endDt := event.EndDt.Time
		startDateHuman := formatDateHuman(startDt)
		endDateHuman := formatDateHuman(endDt)
		if event.AllDay {
			if startDateHuman == endDateHuman {
				when = fmt.Sprintf("%s (all-day)", startDateHuman)
			} else {
				when = fmt.Sprintf("%s - %s (all-day)", startDateHuman, endDateHuman)
			}
		} else {
			if startDateHuman == endDateHuman {
				when = fmt.Sprintf("%s. %s - %s", formatDateHuman(startDt), formatTime(startDt), formatTime(endDt))
			} else {
				when = fmt.Sprintf("%s - %s", formatDtHuman(startDt), formatDtHuman(endDt))
			}
		}

		desc := fmt.Sprintf("Who: %s\nWhen: %s\nWhere: %s", who, when, where)

		eventUrl := fmt.Sprintf("https://teamup.com/%s/events/%s",
			finder.CalendarKey,
			event.ID)

		r := contracts.Result{
			ID:          fmt.Sprintf("%s::%s", finderName, event.ID),
			Title:       event.Title,
			Description: desc,
			URL:         eventUrl,
			ThumbURL:    logoUrl,
			FinderName:  finderName,
		}

		r.Text = fmt.Sprintf("`Teamup Event`\n*%s*\n%s\n\n––\nOpen in browser:\n%s",
			r.Title,
			util.Truncate(util.EscapeMarkdown(r.Description), "...\\[redacted]"),
			r.URL)

		result = append(result, r)
	}

	return &result, nil
}

func formatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

func formatTime(t time.Time) string {
	return t.Format("15:04")
}

func formatDateHuman(t time.Time) string {
	return t.Format("Mon, Jan 2")
}

func formatDtHuman(t time.Time) string {
	return t.Format("Mon, Jan 2, 15:04")
}
