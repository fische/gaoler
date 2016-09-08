package dependency

import (
	"strings"

	"github.com/fische/gaoler/project/dependency/pkg"
)

type Set map[string]*Dependency

func NewSet() Set {
	return make(map[string]*Dependency)
}

func (s Set) GetDependencyOf(p *pkg.Package) *Dependency {
	for root, dep := range s {
		if strings.HasPrefix(p.Path, root) {
			return dep
		}
	}
	return nil
}

func (s Set) ContainsDependencyOf(p *pkg.Package) bool {
	for root := range s {
		if strings.HasPrefix(p.Path, root) {
			return true
		}
	}
	return false
}

func (s Set) Add(p *pkg.Package, ignoreVendor bool) (added bool, err error) {
	if p.Root {
		return
	}
	if dep := s.GetDependencyOf(p); dep != nil {
		added = dep.Add(p)
	} else {
		var dep *Dependency
		dep, err = New(p)
		if err != nil {
			return
		}
		s[dep.RootPackage] = dep
		added = true
	}
	return
}

func (s Set) GetDependencies() []*Dependency {
	deps := make([]*Dependency, len(s))
	idx := 0
	for _, dep := range s {
		deps[idx] = dep
		idx++
	}
	return deps
}
