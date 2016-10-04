package modules

import (
	"errors"

	"github.com/fische/gaoler/vcs"
	"github.com/fische/gaoler/vcs/modules/git"
)

var impl = make(map[string]vcs.VCS)

func Register(v vcs.VCS) {
	impl[v.GetName()] = v
}

func OpenRepository(path string) (repo vcs.Repository, err error) {
	for _, v := range impl {
		repo, err = v.Open(path)
		if err == nil {
			return
		}
	}
	return nil, errors.New("Could not open repository")
}

func GetVCS(vcsName string) vcs.VCS {
	return impl[vcsName]
}

func VCS() []string {
	ret := make([]string, len(impl))
	idx := 0
	for v := range impl {
		ret[idx] = v
		idx++
	}
	return ret
}

func init() {
	// Register modules here
	Register(&git.Git{})
}
