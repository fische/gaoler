package dependency

import (
	"github.com/fische/gaoler/pkg"
	"github.com/fische/gaoler/vcs"
	"github.com/fische/gaoler/vcs/modules"
)

func (d *Dependency) OpenRepository(dir string) error {
	var (
		repo vcs.Repository
		err  error
	)
	if repo, err = modules.OpenRepository(dir); err != nil {
		return err
	}
	d.repository = repo
	return nil
}

func (d *Dependency) LockCurrentState() error {
	var (
		path string
		err  error

		rootPackage string
		remote      string
		revision    string
		branch      string
	)
	if path, err = d.repository.GetPath(); err != nil {
		return err
	} else if rootPackage, err = pkg.GetPackagePath(path); err != nil {
		return err
	} else if remote, err = d.repository.GetRemote(); err != nil {
		return err
	} else if revision, err = d.repository.GetRevision(); err != nil {
		return err
	} else if branch, err = d.repository.GetBranch(); err != nil {
		return err
	}
	d.VCS = d.repository.GetVCSName()
	d.rootPackage = rootPackage
	d.Remote = remote
	d.Revision = revision
	d.Branch = branch
	return nil
}

func (dep Dependency) HasOpenedRepository() bool {
	return dep.repository != nil
}
