package trello

type Card struct {
	ID       string `json:"id"`
	Desc     string `json:"desc"`
	Name     string `json:"name"`
	URL      string `json:"url"`
	ShortURL string `json:"shortUrl"`
}
