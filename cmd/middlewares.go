package cmd

import (
	"os"

	"github.com/fische/gaoler/pkg"
	"github.com/fische/gaoler/project/dependency"
)

func resetDir(dir string) error {
	err := os.RemoveAll(dir)
	if err != nil {
		return err
	}
	return os.MkdirAll(dir, 0775)
}

func importPackage(srcPath string, ignoreVendor bool) func(p *pkg.Package) (string, error) {
	return func(p *pkg.Package) (nextDirectory string, err error) {
		if p.IsPseudoPackage() {
			return
		} else if err = p.Import(srcPath, ignoreVendor); err != nil {
			return
		} else if !p.IsStandardPackage() {
			nextDirectory = p.Dir()
		}
		return
	}
}

func importDependency(srcPath string, ignoreVendor, force bool) func(dep *dependency.Dependency) error {
	return func(dep *dependency.Dependency) error {
		for _, p := range dep.Packages {
			if err := p.Import(srcPath, ignoreVendor); !force && err != nil {
				return err // TODO: Print warning on error if force is false
			}
		}
		return nil
	}
}

func filterUsefulDependencies(p *pkg.Package) bool {
	return !(p.IsLocal() || p.IsPseudoPackage() || p.IsStandardPackage())
}
