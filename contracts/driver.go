package contracts

type Driver interface {
	Find(query string) (*[]Result, error)
	DriverName() string
}
