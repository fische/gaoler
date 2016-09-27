package dependency

import (
	"github.com/fische/gaoler/pkg"
	"github.com/fische/gaoler/vcs"
	"github.com/fische/gaoler/vcs/modules"
)

func (d *Dependency) OpenRepository(dir string) error {
	var (
		err  error
		path string

		repo        vcs.Repository
		rootPackage string
		remote      string
		revision    string
	)
	if repo, err = modules.OpenRepository(dir); err != nil {
		return err
	} else if path, err = repo.GetPath(); err != nil {
		return err
	} else if rootPackage, err = pkg.GetPackagePath(path); err != nil {
		return err
	} else if remote, err = repo.GetRemote(); err != nil {
		return err
	} else if revision, err = repo.GetRevision(); err != nil {
		return err
	}
	d.VCS = repo.GetVCSName()
	d.repository = repo
	d.rootPackage = rootPackage
	d.Remote = remote
	d.Revision = revision
	return nil
}

func (dep *Dependency) HasOpenedRepository() bool {
	return dep.repository != nil
}
