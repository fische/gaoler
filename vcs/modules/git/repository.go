package git

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/fische/gaoler/vcs"
)

type Repository struct {
	Path string
}

const (
	cmd               = "git"
	defaultRemoteName = "origin"
)

var (
	remoteRegexp = regexp.MustCompile("(.+?)\\s+(.+)\\s+\\((push|fetch)\\)")
)

func InitRepository(path string, bare bool) (vcs.Repository, error) {
	var err error
	if path, err = filepath.Abs(path); err != nil {
		return nil, err
	}
	args := []string{"init"}
	if bare {
		args = append(args, "--bare")
	}
	args = append(args, path)
	if o, err := exec.Command(cmd, args...).CombinedOutput(); err != nil {
		return nil, errors.New(string(o))
	}
	return &Repository{
		Path: path,
	}, nil
}

func OpenRepository(path string) (vcs.Repository, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	} else if err = os.Chdir(path); err != nil {
		return nil, err
	}
	var p []byte
	if p, err = exec.Command(cmd, "rev-parse", "--show-toplevel").CombinedOutput(); err != nil {
		return nil, errors.New(string(p))
	} else if err = os.Chdir(dir); err != nil {
		return nil, err
	} else if len(p) == 0 {
		// TODO: Implement support for bare repositories
		return nil, errors.New("Do not support bare repository")
	}
	return &Repository{
		Path: string(p[:len(p)-1]),
	}, nil
}

func (r Repository) GetRevision() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	} else if err = os.Chdir(r.Path); err != nil {
		return "", err
	}
	var rev []byte
	if rev, err = exec.Command(cmd, "describe", "--always").CombinedOutput(); err != nil {
		return "", errors.New(string(rev))
	} else if err = os.Chdir(dir); err != nil {
		return "", err
	}
	return string(rev[:len(rev)-1]), nil
}

func (r Repository) GetRemote() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	} else if err = os.Chdir(r.Path); err != nil {
		return "", err
	}
	var list []byte
	if list, err = exec.Command(cmd, "remote", "-v").CombinedOutput(); err != nil {
		return "", errors.New(string(list))
	} else if err = os.Chdir(dir); err != nil {
		return "", err
	}
	lines := remoteRegexp.FindAllString(string(list), -1)
	for _, line := range lines {
		res := remoteRegexp.FindStringSubmatch(line)
		if res[1] == defaultRemoteName {
			return res[2], nil
		}
	}
	return "", fmt.Errorf("Could not find remote : %s", defaultRemoteName)
}

func (r Repository) AddRemote(remote string) error {
	if dir, err := os.Getwd(); err != nil {
		return err
	} else if err = os.Chdir(r.Path); err != nil {
		return err
	} else if o, err := exec.Command(cmd, "remote", "add", defaultRemoteName, remote).CombinedOutput(); err != nil {
		return errors.New(string(o))
	} else if err = os.Chdir(dir); err != nil {
		return err
	}
	return nil
}

func (r Repository) Fetch() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	} else if err = os.Chdir(r.Path); err != nil {
		return err
	} else if o, err := exec.Command(cmd, "fetch", "--all").CombinedOutput(); err != nil {
		return errors.New(string(o))
	} else if err = os.Chdir(dir); err != nil {
		return err
	}
	return nil
}

func (r Repository) Checkout(revision string) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	} else if err = os.Chdir(r.Path); err != nil {
		return err
	} else if o, err := exec.Command(cmd, "checkout", revision).CombinedOutput(); err != nil {
		return errors.New(string(o))
	} else if err = os.Chdir(dir); err != nil {
		return err
	}
	return nil
}

func (r Repository) CheckoutBranch(branch string) error {
	return r.Checkout(branch)
}

func (r Repository) GetPath() (string, error) {
	return r.Path, nil
}

func (r Repository) GetBranch() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	} else if err = os.Chdir(r.Path); err != nil {
		return "", err
	}
	var path []byte
	if path, err = exec.Command(cmd, "rev-parse", "--abbrev-ref", "HEAD").CombinedOutput(); err != nil {
		return "", errors.New(string(path))
	} else if err = os.Chdir(dir); err != nil {
		return "", err
	} else if len(path) <= 1 {
		return "", errors.New("No branch found")
	}
	return string(path[:len(path)-1]), nil
}

func (r Repository) GetVCSName() string {
	return vcsName
}
