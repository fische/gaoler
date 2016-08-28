package vcs

type VCS interface {
	New(path string) (Repository, error)
	Open(path string) (Repository, error)
}

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
