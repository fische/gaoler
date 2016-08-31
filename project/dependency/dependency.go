package dependency

import (
	"errors"
	"go/build"
	"path/filepath"
	"strings"

	"github.com/fische/gaoler/vcs"
	"github.com/fische/gaoler/vcs/modules"
)

type Dependency struct {
	RootPackage string
	Repository  vcs.Repository
	Packages    []*Package
}

var (
	srcDirs = build.Default.SrcDirs()
)

func New(p *Package) (*Dependency, error) {
	repo, err := modules.OpenRepository(p.Path())
	if err != nil {
		return nil, err
	}
	for _, dir := range srcDirs {
		path, err := repo.GetPath()
		if err != nil {
			return nil, err
		}
		if strings.HasPrefix(path, dir) {
			rel, err := filepath.Rel(dir, path)
			if err != nil {
				return nil, err
			}
			return &Dependency{
				Repository:  repo,
				Packages:    []*Package{p},
				RootPackage: rel,
			}, nil
		}
	}
	return nil, errors.New("Could not find package in src directories")
}

func (d *Dependency) Add(p *Package) (added bool) {
	if d.HasPackage(p.Path()) {
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
