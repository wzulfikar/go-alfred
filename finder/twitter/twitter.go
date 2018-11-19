package twitter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"github.com/wzulfikar/alfred/contracts"
	"github.com/wzulfikar/alfred/util"
)

const finderName = "twitter_v1"
const baseUrl = "https://api.trello.com"
const logoUrl = "http://logobucket.surge.sh/services/trello-logo-md.png"

type TrelloFinder struct {
	Key   string
	Token string
}

func (finder *TrelloFinder) FinderName() string {
	return finderName
}

func (finder *TrelloFinder) Find(query string) (*[]contracts.Result, error) {
	endpoint := fmt.Sprintf("%s/1/search?query=%s&key=%s&token=%s",
		baseUrl,
		url.QueryEscape(query),
		finder.Key,
		finder.Token)

	req, _ := http.NewRequest("GET", endpoint, nil)
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

	alfredResult := []contracts.Result{}
	for _, card := range result.Cards {
		r := contracts.Result{
			ID:          card.ID,
			Title:       card.Name,
			Description: card.Desc,
			URL:         card.ShortURL,
			ThumbURL:    logoUrl,
			FinderName:  finder.FinderName(),
		}

		r.Text = fmt.Sprintf("*%s*\n%s\n\n––\nView in Trello:\n%s",
			r.Title,
			util.Truncate(util.EscapeMarkdown(r.Description), "...\\[redacted]"),
			r.URL)

		alfredResult = append(alfredResult, r)
	}

	return &alfredResult, nil
}