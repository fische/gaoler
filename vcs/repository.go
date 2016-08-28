package vcs

type Repository interface {
	GetRevision() (string, error)
	GetRemotes() ([]Remote, error)

	AddRemote(remote Remote) error
	Fetch() error
	Checkout(revision string) error
}
