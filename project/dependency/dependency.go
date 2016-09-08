package dependency

import (
	"github.com/fische/gaoler/project/dependency/pkg"
	"github.com/fische/gaoler/vcs"
	"github.com/fische/gaoler/vcs/modules"
)

type Dependency struct {
	RootPackage string `json:"-"`

	repository vcs.Repository
	VCS        string
	Revision   string
	Remote     string
	Branch     string `json:",omitempty"`

	Packages []*pkg.Package
}

func New(p *pkg.Package) (*Dependency, error) {
	var (
		path string
		err  error

		dep = &Dependency{
			Packages: []*pkg.Package{p},
		}
	)
	if dep.repository, err = modules.OpenRepository(p.Dir); err != nil {
		return nil, err
	} else if path, err = dep.repository.GetPath(); err != nil {
		return nil, err
	} else if dep.RootPackage, err = pkg.GetPackagePath(path); err != nil {
		return nil, err
	} else if dep.Remote, err = dep.repository.GetRemote(); err != nil {
		return nil, err
	} else if dep.Revision, err = dep.repository.GetRevision(); err != nil {
		return nil, err
	}
	dep.VCS = dep.repository.GetVCSName()
	return dep, nil
}

func (d *Dependency) Add(p *pkg.Package) (added bool) {
	if d.HasPackage(p.Path) {
		return
	}
	added = true
	d.Packages = append(d.Packages, p)
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
