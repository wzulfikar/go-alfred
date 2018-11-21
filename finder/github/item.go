package github

type Item struct {
	ID        int64   `json:"id"`
	Title     string  `json:"title"`
	Body      string  `json:"body"`
	URL       string  `json:"url"`
	HTMLURL   string  `json:"html_url"`
	User      User    `json:"user"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	ClosedAt  *string `json:"closed_at"`
}
