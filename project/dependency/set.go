package dependency

import (
	"strings"

	"github.com/fische/gaoler/pkg"
)

type Set struct {
	dependencies   map[string]*Dependency
	Filter         func(p *pkg.Package) bool
	OnPackageAdded func(p *pkg.Package, dep *Dependency) error
	OnDecoded      func(dep *Dependency) error
}

func NewSet() *Set {
	return &Set{
		dependencies: make(map[string]*Dependency),
	}
}

func (s *Set) fromSet(o *pkg.Set) error {
	for len(o.Packages()) > 0 {
		var dep *Dependency
		for pkgPath, p := range o.Packages() {
			if s.Filter == nil || s.Filter(p) {
				added := true
				if dep == nil {
					dep = New(p)
				} else if strings.HasPrefix(p.Path(), dep.rootPackage) ||
					strings.HasPrefix(dep.rootPackage, p.Path()) {
					added = dep.Add(p)
				} else {
					continue
				}
				if added && s.OnPackageAdded != nil {
					if err := s.OnPackageAdded(p, dep); err != nil {
						return err
					}
				}
			}
			o.Remove(pkgPath)
		}
		if dep != nil {
			s.dependencies[dep.rootPackage] = dep
		}
	}
	return nil
}

func (s *Set) MergePackageSet(o *pkg.Set) error {
	for rootPackage, dep := range s.dependencies {
		for pkgPath, p := range o.Packages() {
			var added bool
			if strings.HasPrefix(pkgPath, rootPackage) ||
				strings.HasPrefix(rootPackage, pkgPath) {
				added = dep.Add(p)
				if strings.HasPrefix(rootPackage, pkgPath) && rootPackage != pkgPath {
					delete(s.dependencies, rootPackage)
					s.dependencies[dep.rootPackage] = dep
				}
			} else {
				continue
			}
			if added && s.OnPackageAdded != nil {
				if err := s.OnPackageAdded(p, dep); err != nil {
					return err
				}
			}
			o.Remove(pkgPath)
		}
	}
	return s.fromSet(o)
}
