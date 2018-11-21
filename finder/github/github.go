package github

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
	"github.com/wzulfikar/alfred/contracts"
	"github.com/wzulfikar/alfred/util"
	validator "gopkg.in/go-playground/validator.v9"
)

const finderName = "github_v1"
const baseUrl = "https://api.github.com"
const logoUrl = "http://logobucket.surge.sh/services/github-logo-sm.png"

type GithubFinder struct {
	Token string `validate:"required"`
}

func (finder *GithubFinder) Init() error {
	validate := validator.New()
	errs := validate.Struct(finder)
	if errs != nil {
		return errors.Wrap(errors.New(fmt.Sprintf("%s", errs)), finderName)
	}

	return nil
}

func (finder *GithubFinder) FinderName() string {
	return finderName
}

func (finder *GithubFinder) Find(query string) (*[]contracts.Result, error) {
	// search issue endpoint
	endpoint := fmt.Sprintf("%s/search/issues?q=%s",
		baseUrl,
		url.QueryEscape(query))

	req, _ := http.NewRequest("GET", endpoint, nil)
	req.SetBasicAuth(finder.Token, "x-oauth-basic")
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
	for _, item := range searchResult.Items {
		r := contracts.Result{
			ID:          fmt.Sprintf("%s::%s", finderName, strconv.Itoa(int(item.ID))),
			Title:       item.Title,
			Description: item.Body,
			URL:         item.HTMLURL,
			ThumbURL:    logoUrl,
			FinderName:  finderName,
		}

		r.Text = fmt.Sprintf("`Github Issue (submitted by %s)`\n*%s*\n%s\n\n––\nOpen issue in browser:\n%s",
			item.User.Login,
			r.Title,
			util.Truncate(util.EscapeMarkdown(r.Description), "...\\[redacted]"),
			r.URL)

		result = append(result, r)
	}

	return &result, nil
}
