package youtrack

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

const logo = "http://logobucket.surge.sh/services/youtrack-logo-md.png"

type YoutrackDriver struct {
	BaseUrl string
	Token   string
	Fields  string
}

const defaultFields = "id,summary,description,updated,created,votes,numberInProject,project(shortName),tags(name)"

func (driver *YoutrackDriver) DriverName() string {
	return "youtrack"
}

func (driver *YoutrackDriver) issueUrl(projectShortName string, numberInProject int) string {
	return fmt.Sprintf("%s/issue/%s-%d", driver.BaseUrl, projectShortName, numberInProject)
}

func (driver *YoutrackDriver) Find(query string) (*[]contracts.Result, error) {
	log.Println("fetching issue from youtrack..")

	var client = &http.Client{}

	fields := driver.Fields
	if fields == "" {
		fields = defaultFields
	}

	endpoint := fmt.Sprintf("%s/api/issues/?query=%s&fields=%s", driver.BaseUrl, url.QueryEscape(query), fields)
	req, err := http.NewRequest("GET", endpoint, nil)
	req.Header.Add("Authorization", "Bearer "+driver.Token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	issues := &[]Issue{}
	if err := json.Unmarshal(body, issues); err != nil {
		log.Println("failed to unmarshal body:", string(body))
		return nil, err
	}

	alfredResult := []contracts.Result{}
	for _, issue := range *issues {
		title := fmt.Sprintf("[%s-%d] %s",
			issue.Project.ShortName,
			issue.NumberInProject,
			issue.Summary)

		r := contracts.Result{
			ID:          issue.ID,
			Title:       title,
			Description: issue.Description,
			URL:         driver.issueUrl(issue.Project.ShortName, issue.NumberInProject),
			ThumbURL:    logo,
			DriverName:  driver.DriverName(),
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
