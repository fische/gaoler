package dependency

import (
	"strings"

	"github.com/fische/gaoler/pkg"
)

type Set struct {
	Deps           map[string]*Dependency
	Filter         func(p *pkg.Package) bool
	OnPackageAdded func(p *pkg.Package, dep *Dependency) error
	OnDecoded      func(dep *Dependency) error
}

func NewSet() *Set {
	return &Set{
		Deps: make(map[string]*Dependency),
	}
}

func (deps *Set) fromSet(s *pkg.Set) error {
	for len(s.Packages) > 0 {
		var dep *Dependency
		for pkgPath, p := range s.Packages {
			if deps.Filter == nil || deps.Filter(p) {
				added := true
				if dep == nil {
					dep = New(p)
				} else if strings.HasPrefix(p.Path, dep.RootPackage) {
					added = dep.Add(p)
				} else if strings.HasPrefix(dep.RootPackage, p.Path) {
					added = dep.Add(p)
					dep.RootPackage = p.Path
				} else {
					continue
				}
				if added && deps.OnPackageAdded != nil {
					if err := deps.OnPackageAdded(p, dep); err != nil {
						return err
					}
				}
			}
			s.Remove(pkgPath)
		}
		if dep != nil {
			deps.Deps[dep.RootPackage] = dep
		}
	}
	return nil
}

func (deps *Set) MergePackageSet(s *pkg.Set) error {
	for rootPackage, dep := range deps.Deps {
		for pkgPath, p := range s.Packages {
			var added bool
			if strings.HasPrefix(pkgPath, rootPackage) {
				added = dep.Add(p)
			} else if strings.HasPrefix(rootPackage, pkgPath) {
				added = dep.Add(p)
				dep.RootPackage = pkgPath
				delete(deps.Deps, rootPackage)
				deps.Deps[dep.RootPackage] = dep
			} else {
				continue
			}
			if added && deps.OnPackageAdded != nil {
				if err := deps.OnPackageAdded(p, dep); err != nil {
					return err
				}
			}
			s.Remove(pkgPath)
		}
	}
	return deps.fromSet(s)
}

func (deps *Set) ToPackageSet() *pkg.Set {
	s := pkg.NewSet()
	for _, dep := range deps.Deps {
		for _, p := range dep.Packages {
			s.Insert(p, true)
		}
	}
	return s
}
