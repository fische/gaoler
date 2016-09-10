package vcs

import "os"

type VCS interface {
	New(path string) (Repository, error)
	Open(path string) (Repository, error)
	GetName() string
}

func CloneAtRevision(v VCS, remote, revision, path string) (Repository, error) {
	var ret Repository
	if err := os.MkdirAll(path, 0775); err != nil {
		return nil, err
	} else if ret, err = v.New(path); err != nil {
		return nil, err
	} else if err = ret.AddRemote(remote); err != nil {
		return nil, err
	} else if err = ret.Fetch(); err != nil {
		return nil, err
	} else if err = ret.Checkout(revision); err != nil {
		return nil, err
	}
	return ret, nil
}

func CloneRepository(v VCS, repo Repository, path string) (Repository, error) {
	var ret Repository
	if rev, err := repo.GetRevision(); err != nil {
		return nil, err
	} else if remote, err := repo.GetRemote(); err != nil {
		return nil, err
	} else if ret, err = CloneAtRevision(v, path, rev, remote); err != nil {
		return nil, err
	}
	return ret, nil
}
