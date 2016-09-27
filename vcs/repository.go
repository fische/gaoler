package vcs

type Repository interface {
	GetRevision() (string, error)
	GetRemote() (string, error)
	GetVCSName() string
	GetPath() (string, error)
	GetBranch() (string, error)

	AddRemote(remote string) error
	Fetch() error
	Checkout(revision string) error
	CheckoutBranch(branch string) error
}
