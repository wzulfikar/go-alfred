package contracts

type Finder interface {
	Find(query string) (*[]Result, error)

	// returns finder name (identification purpose)
	// eg. trello_v1, youtrack_v1
	FinderName() string
}
