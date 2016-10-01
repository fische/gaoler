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
	d.Repository = repo
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
	if path, err = d.Repository.GetPath(); err != nil {
		return err
	} else if rootPackage, err = pkg.GetPackagePath(path); err != nil {
		return err
	} else if remote, err = d.Repository.GetRemote(); err != nil {
		return err
	} else if revision, err = d.Repository.GetRevision(); err != nil {
		return err
	} else if branch, err = d.Repository.GetBranch(); err != nil {
		return err
	}
	d.VCS = d.Repository.GetVCSName()
	d.RootPackage = rootPackage
	d.Remote = remote
	d.Revision = revision
	d.Branch = branch
	return nil
}
