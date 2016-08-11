package vcs

//Repository provides methods to interact with its Version Control System
type Repository interface {
	//GetRevision returns the current revision of the repository
	GetRevision() (string, error)
	//GetRemotes returns repository remotes
	GetRemotes() ([]Remote, error)

	//AddRemote adds given remote to the repository
	AddRemote(remote Remote) error
	//Fetch fetches repository from `remote`
	Fetch() error
	//Checkout checks repository out to `revision`
	Checkout(revision string) error
}
