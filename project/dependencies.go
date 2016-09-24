package project

import (
	"strings"

	"github.com/fische/gaoler/pkg"
	"github.com/fische/gaoler/pkg/set"
	"github.com/fische/gaoler/project/dependency"
)

type Dependencies map[string]*dependency.Dependency

func (deps Dependencies) Set(s *set.Packages, openRepository bool) {
	for len(s.Packages()) > 0 {
		var (
			dep *dependency.Dependency
		)
		for k, p := range s.Packages() {
			if dep == nil {
				dep = dependency.New(p)
			} else if strings.HasPrefix(p.Path(), dep.RootPackage()) {
				dep.Add(p)
			} else {
				continue
			}
			if openRepository && !dep.IsOpened() && !p.IsVendored() {
				dep.OpenRepository(p) //TODO Handle when error is returned
			}
			s.Delete(k)
		}
		deps[dep.RootPackage()] = dep
	}
}

func (deps Dependencies) Merge(s *set.Packages, openRepository bool) {
	for rootPackage, dep := range deps {
		for pkgPath, p := range s.Packages() {
			if strings.HasPrefix(pkgPath, rootPackage) {
				dep.Add(p)
				s.Delete(pkgPath)
			}
		}
	}
	deps.Set(s, openRepository)
}

func (deps Dependencies) GetPackagesSet() (pkgs *set.Packages) {
	pkgs = set.NewPackages()
	for _, dep := range deps {
		pkgs.ForceInsertDependency(dep)
	}
	return
}

func (deps Dependencies) Import(srcPath string, flags pkg.Flags) error {
	for _, dep := range deps {
		if err := dep.Import(srcPath, flags); err != nil {
			return err
		}
	}
	return nil
}
