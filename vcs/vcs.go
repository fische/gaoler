package vcs

//VCS represents a Version Control System
type VCS interface {
	//New creates a new repository at `path`
	New(path string) (Repository, error)
	//Open opens an existing repository at `path`
	Open(path string) (Repository, error)
}

//CloneAtRevision clones repository from `remote`
func CloneAtRevision(v VCS, path, revision string, remotes []Remote) (Repository, error) {
	r, err := v.New(path)
	if err != nil {
		return nil, err
	}
	for _, remote := range remotes {
		if err = r.AddRemote(remote); err != nil {
			return nil, err
		}
	}
	if err = r.Fetch(); err != nil {
		return nil, err
	} else if err = r.Checkout(revision); err != nil {
		return nil, err
	}
	return r, nil
}
