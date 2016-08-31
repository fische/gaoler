package dependency

import (
	"fmt"
	"os"
	"path/filepath"

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

func (d Dependency) Vendor(vendorRoot string) error {
	v, _ := modules.GetVCS(d.Repository.GetVCSName())
	path := filepath.Clean(fmt.Sprintf("%s/%s/", vendorRoot, d.RootPackage))
	_, err := vcs.CloneRepository(v, path, d.Repository)
	if err != nil {
		return err
	}
	return d.CleanVendor(vendorRoot)
}

func (d Dependency) CleanVendor(vendorRoot string) error {
	root := filepath.Clean(fmt.Sprintf("%s/%s/", vendorRoot, d.RootPackage))
	removeRootFiles := !d.HasPackage(d.RootPackage)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			var rel string
			rel, err = filepath.Rel(vendorRoot, path)
			if err != nil {
				return err
			}
			if rel != d.RootPackage && !d.HasPackage(rel) {
				err = os.RemoveAll(path)
				if err != nil {
					return err
				}
				return filepath.SkipDir
			}
		} else if removeRootFiles && path == filepath.Clean(fmt.Sprintf("%s/%s", root, info.Name())) {
			err = os.Remove(path)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
