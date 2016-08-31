package dependency

import (
	"github.com/fische/gaoler/vcs"
	"github.com/fische/gaoler/vcs/modules"
)

type Dependency struct {
	RootPackage string
	Repository  vcs.Repository
	Packages    []*Package
}

func New(p *Package) (*Dependency, error) {
	repo, err := modules.OpenRepository(p.Path())
	if err != nil {
		return nil, err
	}
	path, err := repo.GetPath()
	if err != nil {
		return nil, err
	}
	pkgPath, err := GetPackagePathFromPath(path)
	if err != nil {
		return nil, err
	}
	return &Dependency{
		Repository:  repo,
		Packages:    []*Package{p},
		RootPackage: pkgPath,
	}, nil
}

func (d *Dependency) Add(p *Package) (added bool) {
	if d.HasPackage(p.Name()) {
		return
	}
	added = true
	d.Packages = append(d.Packages, p)
	return
}

func (d Dependency) HasPackage(packagePath string) bool {
	for _, pkg := range d.Packages {
		if pkg.Name() == packagePath {
			return true
		}
	}
	return false
}
