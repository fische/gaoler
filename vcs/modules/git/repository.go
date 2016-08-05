package git

import (
	"os"
	"os/exec"

	"github.com/fische/vcs"
)

//Repository represents a git repository
type Repository struct {
	Path string
}

const cmd = "git"

//InitRepository inits a new Git repository at `path`
func InitRepository(path string, bare bool) (vcs.Repository, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	} else if err = os.Chdir(path); err != nil {
		return nil, err
	}
	args := []string{"init"}
	if bare {
		args = append(args, "--bare")
	}
	if err = exec.Command(cmd, args...).Run(); err != nil {
		return nil, err
	} else if err = os.Chdir(dir); err != nil {
		return nil, err
	}
	return &Repository{
		Path: path,
	}, nil
}

//OpenRepository opens a new git repository using `path`
func OpenRepository(path string) (vcs.Repository, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	} else if err = os.Chdir(path); err != nil {
		return nil, err
	}
	var p []byte
	if p, err = exec.Command(cmd, "rev-parse", "--show-toplevel").Output(); err != nil {
		return nil, err
	} else if err = os.Chdir(dir); err != nil {
		return nil, err
	}
	return &Repository{
		Path: string(p[:len(p)-1]),
	}, nil
}

//GetRevision returns the current revision of the repository
func (r Repository) GetRevision() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	} else if err = os.Chdir(r.Path); err != nil {
		return "", err
	}
	var rev []byte
	if rev, err = exec.Command(cmd, "describe", "--always").Output(); err != nil {
		return "", err
	} else if err = os.Chdir(dir); err != nil {
		return "", err
	}
	return string(rev[:len(rev)-1]), nil
}

//Fetch fetches repository from `remote`
func (r Repository) Fetch(remote string) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	} else if err = os.Chdir(r.Path); err != nil {
		return err
	} else if err = exec.Command(cmd, "fetch", remote).Run(); err != nil {
		return err
	} else if err = os.Chdir(dir); err != nil {
		return err
	}
	return nil
}

//Checkout checks repository out to `revision`
func (r Repository) Checkout(revision string) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	} else if err = os.Chdir(r.Path); err != nil {
		return err
	} else if err = exec.Command(cmd, "checkout", revision).Run(); err != nil {
		return err
	} else if err = os.Chdir(dir); err != nil {
		return err
	}
	return nil
}
