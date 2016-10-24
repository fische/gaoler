package dependency

import (
	"errors"

	"github.com/fische/gaoler/vcs"
)

type State struct {
	vcs      string
	revision string
	remote   string
	branch   string
}

func (d *Dependency) LockCurrentState() error {
	if d.repository == nil {
		return errors.New("Package repository has not been opened.")
	}
	state, err := newState(d.repository)
	if err != nil {
		return err
	}
	d.State = state
	return nil
}

func newState(r vcs.Repository) (*State, error) {
	var (
		err      error
		remote   string
		revision string
		branch   string
	)
	if remote, err = r.GetRemote(); err != nil {
		return nil, err
	} else if revision, err = r.GetRevision(); err != nil {
		return nil, err
	} else if branch, err = r.GetBranch(); err != nil {
		return nil, err
	}
	return &State{
		vcs:      r.GetVCSName(),
		remote:   remote,
		revision: revision,
		branch:   branch,
	}, nil
}

func (s State) VCS() string {
	return s.vcs
}

func (s State) Revision() string {
	return s.revision
}

func (s State) Remote() string {
	return s.remote
}

func (s State) Branch() string {
	return s.branch
}
