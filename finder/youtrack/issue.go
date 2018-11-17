package youtrack

type Issue struct {
	ID              string  `json:"id"`
	Summary         string  `json:"summary"`
	Description     string  `json:"description"`
	Tags            []Tag   `json:"tags"`
	Created         int64   `json:"created"`
	Votes           int64   `json:"votes"`
	Updated         int64   `json:"updated"`
	NumberInProject int     `json:"numberInProject"`
	Project         Project `json:"project"`
}
