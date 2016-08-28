package modules

import (
	"github.com/fische/gaoler/errors"
	"github.com/fische/gaoler/vcs"
	"github.com/fische/gaoler/vcs/modules/git"
)

var impl = make(map[string]vcs.VCS)

func Register(name string, v vcs.VCS) {
	impl[name] = v
}

func GetRepository(path string) (vcsName string, repo vcs.Repository, err error) {
	var v vcs.VCS
	for vcsName, v = range impl {
		repo, err = v.Open(path)
		if err == nil {
			return
		}
	}
	return "", nil, errors.ErrNotValidRepository
}

func GetVCS(vcsName string) (vcs.VCS, bool) {
	f, ok := impl[vcsName]
	return f, ok
}

func init() {
	// Register modules here
	Register("git", &git.Git{})
}
