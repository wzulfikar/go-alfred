package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type Issue struct {
	Description string `json:"description"`
	Created     int64  `json:"created"`
	Summary     string `json:"summary"`
	Votes       int64  `json:"votes"`
	Updated     int64  `json:"updated"`
	ID          string `json:"id"`
	Type        string `json:"$type"`
}

const baseUrl = "https://yt.luxtag.io"
const token = "Bearer perm:d2lsZGFu.UHl0aG9uIENsaWVudA==.t8E2hTxVJvghsR61oFy9fqngoPobq3"

const query = "tag: knowledge state: Resolved sort by: {resolved date} resolved date: Today"
const fields = "id,summary,description,updated,created,votes"

func GetLink(issueId string) string {
	return fmt.Sprintf("%s/issue/%s", baseUrl, issueId)
}

func SearchIssue(query string) (*[]Issue, error) {
	log.Println("searching issue..")
	endpoint := "/api/issues/?query=%s&fields=%s", url.QueryEscape(query), fields
	return Fetch(endpoint)
}

func FetchKnowledge() (*[]Issue, error) {
	log.Println("fetching knowledge..")
	query := fmt.Sprintf("query=%s&fields=%s", url.QueryEscape(query), fields)
	return FetchIssue(query)
}

func FetchIssue(query string) (*[]Issue, error) {
	log.Println("fetching issue from youtrack..")

	var client = &http.Client{}

	const endpoint = baseUrl + "/api/issues/?" + query
	req, err := http.NewRequest("GET", endpoint, nil)
	req.Header.Add("Authorization", token)

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
		return nil, err
	}

	return issues, nil
}
