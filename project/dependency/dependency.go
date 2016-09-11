package dependency

import (
	"github.com/fische/gaoler/project/dependency/pkg"
	"github.com/fische/gaoler/vcs"
	"github.com/fische/gaoler/vcs/modules"
)

type Dependency struct {
	RootPackage string `json:"-" yaml:"-"`

	repository vcs.Repository
	VCS        string
	Revision   string
	Remote     string
	Branch     string `json:",omitempty" yaml:",omitempty"`

	Packages []*pkg.Package
}

func New(p *pkg.Package) (*Dependency, error) {
	dep := &Dependency{
		Packages: []*pkg.Package{p},
	}
	if !p.IsVendored() {
		dep.SetRepository(p)
	} else {
		dep.RootPackage = p.Path
	}
	return dep, nil
}

func (d *Dependency) Add(p *pkg.Package) (added bool) {
	if d.HasPackage(p.Path) {
		return
	}
	added = true
	d.Packages = append(d.Packages, p)
	if d.repository == nil && !p.IsVendored() {
		d.SetRepository(p)
	}
	return
}

func (d Dependency) HasPackage(packagePath string) bool {
	for _, pkg := range d.Packages {
		if pkg.Path == packagePath {
			return true
		}
	}
	return false
}

func (d Dependency) IsVendored() bool {
	for _, p := range d.Packages {
		if !p.IsVendored() {
			return false
		}
	}
	return true
}

func (d *Dependency) SetRepository(p *pkg.Package) error {
	var (
		err  error
		path string
	)
	if d.repository, err = modules.OpenRepository(p.Dir); err != nil {
		return err
	} else if path, err = d.repository.GetPath(); err != nil {
		return err
	} else if d.RootPackage, err = pkg.GetPackagePath(path); err != nil {
		return err
	} else if d.Remote, err = d.repository.GetRemote(); err != nil {
		return err
	} else if d.Revision, err = d.repository.GetRevision(); err != nil {
		return err
	}
	d.VCS = d.repository.GetVCSName()
	return nil
}
