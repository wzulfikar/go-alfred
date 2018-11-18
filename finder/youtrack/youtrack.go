package youtrack

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

const finderName = "youtrack_v1"
const logo = "http://logobucket.surge.sh/services/youtrack-logo-md.png"

type YoutrackFinder struct {
	BaseUrl string
	Token   string
	Fields  string
}

const defaultFields = "id,summary,description,updated,created,votes,numberInProject,project(shortName),tags(name)"

func (finder *YoutrackFinder) FinderName() string {
	return finderName
}

func (finder *YoutrackFinder) issueUrl(projectShortName string, numberInProject int) string {
	return fmt.Sprintf("%s/issue/%s-%d", finder.BaseUrl, projectShortName, numberInProject)
}

func (finder *YoutrackFinder) Find(query string) (*[]contracts.Result, error) {
	log.Println("fetching issue from youtrack..")

	fields := finder.Fields
	if fields == "" {
		fields = defaultFields
	}

	endpoint := fmt.Sprintf("%s/api/issues/?query=%s&fields=%s",
		finder.BaseUrl,
		url.QueryEscape(query),
		fields)

	req, _ := http.NewRequest("GET", endpoint, nil)
	req.Header.Add("Authorization", "Bearer "+finder.Token)

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

	alfredResult := []contracts.Result{}
	for _, issue := range *searchResult {
		title := fmt.Sprintf("[%s-%d] %s",
			issue.Project.ShortName,
			issue.NumberInProject,
			issue.Summary)

		r := contracts.Result{
			ID:          issue.ID,
			Title:       title,
			Description: issue.Description,
			URL:         finder.issueUrl(issue.Project.ShortName, issue.NumberInProject),
			ThumbURL:    logo,
			FinderName:  finder.FinderName(),
		}

		// add book icon to signify knowledge cards.
		// TODO: move result decorator to consumer's domain.
		for _, tag := range issue.Tags {
			if tag.Name == "knowledge" {
				r.Title = fmt.Sprintf("%s 📖", r.Title)
			}
		}

		r.Text = fmt.Sprintf("*%s*\n%s\n\n––\nView in YouTrack:\n%s",
			r.Title,
			util.Truncate(util.EscapeMarkdown(r.Description), "...\\[redacted]"),
			r.URL)

		alfredResult = append(alfredResult, r)
	}

	return &alfredResult, nil
}
