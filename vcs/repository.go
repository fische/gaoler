package vcs

//Repository provides methods to interact with its Version Control System
type Repository interface {
	//GetRevision returns the current revision of the repository
	GetRevision() (string, error)
	//Fetch fetches repository from `remote`
	Fetch(remote string) error
	//Checkout checks repository out to `revision`
	Checkout(revision string) error
}
