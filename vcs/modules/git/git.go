package git

import "github.com/fische/gaoler/vcs"

//Git implements `vcs.VCS` for git
type Git struct{}

//New creates a new Repository at `path`
func (g Git) New(path string) (vcs.Repository, error) {
	return InitRepository(path, false)
}

//Open opens an existing repository at `path`
func (g Git) Open(path string) (vcs.Repository, error) {
	return OpenRepository(path)
}
