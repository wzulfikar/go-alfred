package github

type SearchResult struct {
	TotalCount        int64  `json:"total_count"`
	IncompleteResults bool   `json:"incomplete_results"`
	Items             []Item `json:"items"`
}
