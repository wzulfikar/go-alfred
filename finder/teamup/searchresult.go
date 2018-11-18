package teamup

type SearchResult struct {
	Events    []Event `json:"events"`
	Timestamp int64   `json:"timestamp"`
}
