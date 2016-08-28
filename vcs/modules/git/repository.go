package git

import (
	"os"
	"os/exec"
	"regexp"

	"github.com/fische/gaoler/errors"
	"github.com/fische/gaoler/vcs"
)

type Repository struct {
	Path string
}

const (
	cmd = "git"
)

var (
	remoteRegexp = regexp.MustCompile("(.+?)\\s+(.+)\\s+\\((push|fetch)\\)")
)

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

func (r Repository) GetRemotes() ([]vcs.Remote, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	} else if err = os.Chdir(r.Path); err != nil {
		return nil, err
	}
	var list []byte
	if list, err = exec.Command(cmd, "remote", "-v").Output(); err != nil {
		return nil, err
	} else if err = os.Chdir(dir); err != nil {
		return nil, err
	}
	lines := remoteRegexp.FindAllString(string(list), -1)
	remotes := make(map[string]vcs.Remote, len(lines))
	for _, line := range lines {
		res := remoteRegexp.FindStringSubmatch(line)
		var (
			remote vcs.Remote
			ok     bool
		)
		if remote, ok = remotes[res[1]]; !ok {
			remotes[res[1]] = &Remote{
				Name: res[1],
				URL:  res[2],
			}
			remote = remotes[res[1]]
		}
		if res[3] == "fetch" {
			remote.(*Remote).Type &= fetch
		} else if res[3] == "push" {
			remote.(*Remote).Type &= push
		}
	}
	var ret []vcs.Remote
	for _, v := range remotes {
		ret = append(ret, v)
	}
	return ret, nil
}

func (r Repository) AddRemote(remote vcs.Remote) error {
	gitRemote, ok := remote.(*Remote)
	if !ok {
		return errors.ErrNotValidRemote
	}
	if dir, err := os.Getwd(); err != nil {
		return err
	} else if err = os.Chdir(r.Path); err != nil {
		return err
	} else if err = exec.Command(cmd, "remote", "add", gitRemote.Name, gitRemote.URL).Run(); err != nil {
		return err
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
	} else if err = exec.Command(cmd, "fetch", "--all").Run(); err != nil {
		return err
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
	} else if err = exec.Command(cmd, "checkout", revision).Run(); err != nil {
		return err
	} else if err = os.Chdir(dir); err != nil {
		return err
	}
	return nil
}
