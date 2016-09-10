package dependency

import (
	"encoding/json"
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

func (s Set) UnmarshalJSON(data []byte) error {
	m := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	for root, value := range m {
		dep := &Dependency{
			RootPackage: root,
		}
		if err := json.Unmarshal(value, dep); err != nil {
			return err
		}
		s[root] = dep
	}
	return nil
}
