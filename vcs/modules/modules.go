package modules

import (
	"github.com/fische/vcs"
	"github.com/fische/vcs/errors"
	"github.com/fische/vcs/modules/git"
)

var impl = make(map[string]vcs.VCS)

//Register registers `v` at `name`
func Register(name string, v vcs.VCS) {
	impl[name] = v
}

//GetRepository opens an existing `Repository` at `path`
//using the appropriate registered `vcs.VCS`
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

//GetVCS returns the `vcs.VCS` registered using `vcsName`
func GetVCS(vcsName string) (vcs.VCS, bool) {
	f, ok := impl[vcsName]
	return f, ok
}

func init() {
	//Register modules here
	Register("git", &git.Git{})
}
