package vcs

import "os"

type VCS interface {
	New(path string) (Repository, error)
	Open(path string) (Repository, error)
	GetName() string
}

func CloneAtRevision(v VCS, path, revision string, remote string) (Repository, error) {
	r, err := v.New(path)
	if err != nil {
		return nil, err
	} else if err = r.AddRemote(remote); err != nil {
		return nil, err
	} else if err = r.Fetch(); err != nil {
		return nil, err
	} else if err = r.Checkout(revision); err != nil {
		return nil, err
	}
	return r, nil
}

func CloneRepository(v VCS, path string, repo Repository) (Repository, error) {
	err := os.MkdirAll(path, 0775)
	if err != nil {
		return nil, err
	}
	rev, err := repo.GetRevision()
	if err != nil {
		return nil, err
	}
	remote, err := repo.GetRemote()
	if err != nil {
		return nil, err
	}
	ret, err := CloneAtRevision(v, path, rev, remote)
	if err != nil {
		return nil, err
	}
	return ret, nil
}
