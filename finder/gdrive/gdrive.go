package gdrive

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"github.com/wzulfikar/alfred/contracts"
)

const finderName = "gdrive_v1"
const baseUrl = "https://api.teamup.com"

type GDriveFinder struct {
	ApiKey      string
	CalendarKey string
}

func (finder *GDriveFinder) FinderName() string {
	return finderName
}

func (finder *GDriveFinder) Find(query string) (*[]contracts.Result, error) {

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
