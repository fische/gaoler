package vcs

//Repository provides methods to interact with its Version Control System
type Repository interface {
	//GetRevision returns the current revision of the repository
	GetRevision() (string, error)

	Fetch(url string) error

	Checkout(revision string) error
}
