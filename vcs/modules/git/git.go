package git

import "github.com/fische/gaoler/vcs"

type Git struct{}

func (g Git) New(path string) (vcs.Repository, error) {
	return InitRepository(path, false)
}

func (g Git) Open(path string) (vcs.Repository, error) {
	return OpenRepository(path)
}
