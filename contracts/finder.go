package contracts

type Finder interface {
	Find(query string) (*[]Result, error)

	// returns finder name (identification purpose)
	// eg. trello_v1, youtrack_v1
	FinderName() string

	// Init initializes finder and verify if the finder
	// contains any error (eg. empty config, etc.).
	// will be called inside `alfred.New()`
	Init() error
}
