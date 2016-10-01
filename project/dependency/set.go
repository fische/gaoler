package dependency

import (
	"strings"

	"github.com/fische/gaoler/pkg"
)

type Set struct {
	deps           map[string]*Dependency
	Filter         func(p *pkg.Package) bool
	OnPackageAdded func(p *pkg.Package, dep *Dependency) error
	OnDecoded      func(dep *Dependency) error
}

func NewSet() *Set {
	return &Set{
		deps: make(map[string]*Dependency),
	}
}

func (deps *Set) fromSet(s *pkg.Set) error {
	for len(s.Packages) > 0 {
		var dep *Dependency
		for pkgPath, p := range s.Packages {
			if deps.Filter != nil && deps.Filter(p) {
				added := true
				if dep == nil {
					dep = New(p)
				} else if strings.HasPrefix(p.Path, dep.RootPackage) {
					added = dep.Add(p)
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
			deps.deps[dep.RootPackage] = dep
		}
	}
	return nil
}

func (deps *Set) MergePackageSet(s *pkg.Set) error {
	for rootPackage, dep := range deps.deps {
		for pkgPath, p := range s.Packages {
			if strings.HasPrefix(pkgPath, rootPackage) {
				if dep.Add(p) && deps.OnPackageAdded != nil {
					if err := deps.OnPackageAdded(p, dep); err != nil {
						return err
					}
				}
				s.Remove(pkgPath)
			}
		}
	}
	return deps.fromSet(s)
}

func (deps *Set) ToPackageSet() *pkg.Set {
	s := pkg.NewSet()
	for _, dep := range deps.deps {
		for _, p := range dep.Packages {
			s.Insert(p, true)
		}
	}
	return s
}

func (deps Set) Deps() map[string]*Dependency {
	return deps.deps
}
