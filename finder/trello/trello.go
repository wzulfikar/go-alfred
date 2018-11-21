package trello

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
	validator "gopkg.in/go-playground/validator.v9"
)

const finderName = "trello_v1"
const baseUrl = "https://api.trello.com"
const logoUrl = "http://logobucket.surge.sh/services/trello-logo-md.png"

type TrelloFinder struct {
	Key   string `validate:"required"`
	Token string `validate:"required"`
}

func (finder *TrelloFinder) Init() error {
	validate := validator.New()
	errs := validate.Struct(finder)
	if errs != nil {
		return errors.Wrap(errors.New(fmt.Sprintf("%s", errs)), finderName)
	}

	return nil
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

	searchResult := &SearchResult{}
	if err := json.Unmarshal(body, searchResult); err != nil {
		log.Println("failed to unmarshal body:", string(body))
		return nil, errors.Wrap(err, finderName)
	}

	result := []contracts.Result{}
	for _, card := range searchResult.Cards {
		r := contracts.Result{
			ID:          card.ID,
			Title:       card.Name,
			Description: card.Desc,
			URL:         card.ShortURL,
			ThumbURL:    logoUrl,
			FinderName:  finderName,
		}

		if r.Description == "" {
			r.Description = "(No description)"
		}

		r.Text = fmt.Sprintf("`Trello Card`\n*%s*\n%s\n\n––\nOpen in browser:\n%s",
			r.Title,
			util.Truncate(util.EscapeMarkdown(r.Description), "...\\[redacted]"),
			r.URL)

		result = append(result, r)
	}

	return &result, nil
}
