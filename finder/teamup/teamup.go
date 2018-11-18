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
	"github.com/wzulfikar/alfred/contracts"
)

type TeamupFinder struct {
	ApiKey      string
	CalendarKey string
}

const finderName = "teamup_v1"
const baseUrl = "https://api.teamup.com"

func (finder *TeamupFinder) FinderName() string {
	return finderName
}

func (finder *TeamupFinder) Find(query string) (*[]contracts.Result, error) {
	now := time.Now()

	// today -30d
	startDate := now.Add(time.Hour * 24 * 30)

	// today +30d
	endDate := now.Add(time.Hour * 24 * 30 * -1)

	endpoint := fmt.Sprintf("%s/%s/events?query=%s&startDate=%s&endDate=%s",
		baseUrl,
		finder.CalendarKey,
		url.QueryEscape(query),
		startDate,
		endDate)

	req, _ := http.NewRequest("GET", url, nil)
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

	result := &SearchResult{}
	if err := json.Unmarshal(body, result); err != nil {
		log.Println("failed to unmarshal body:", string(body))
		return nil, errors.Wrap(err, finderName)
	}

	return searchResult, nil
}
