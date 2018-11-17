package trello

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/wzulfikar/alfred/contracts"
	"github.com/wzulfikar/alfred/util"
)

const logoUrl = "http://logobucket.surge.sh/services/trello-logo-md.png"
const baseUrl = "https://api.trello.com"

type TrelloFinder struct {
	Key   string
	Token string
}

func (finder *TrelloFinder) FinderName() string {
	return "trello"
}

func (finder *TrelloFinder) Find(query string) (*[]contracts.Result, error) {
	log.Println("fetching cards from trello..")

	var client = &http.Client{}
	endpoint := fmt.Sprintf("%s/1/search?query=%s&key=%s&token=%s",
		baseUrl,
		url.QueryEscape(query),
		finder.Key,
		finder.Token)

	req, err := http.NewRequest("GET", endpoint, nil)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result := &Result{}
	if err := json.Unmarshal(body, result); err != nil {
		log.Println("failed to unmarshal body:", string(body))
		return nil, err
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
