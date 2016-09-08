package dependency

import (
	"github.com/fische/gaoler/project/dependency/pkg"
	"github.com/fische/gaoler/vcs"
	"github.com/fische/gaoler/vcs/modules"
)

type Dependency struct {
	RootPackage string `json:"-"`
	Repository  vcs.Repository
	Packages    []*pkg.Package
}

func New(p *pkg.Package) (*Dependency, error) {
	var (
		pkgPath string
		path    string
		repo    vcs.Repository
		err     error
	)
	if repo, err = modules.OpenRepository(p.Dir); err != nil {
		return nil, err
	} else if path, err = repo.GetPath(); err != nil {
		return nil, err
	} else if pkgPath, err = pkg.GetPackagePath(path); err != nil {
		return nil, err
	}
	return &Dependency{
		RootPackage: pkgPath,
		Repository:  repo,
		Packages:    []*pkg.Package{p},
	}, nil
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
