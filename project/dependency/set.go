package dependency

import (
	"go/ast"
	"strings"
)

type Set struct {
	Dependencies map[string]*Dependency
}

func NewSet() *Set {
	return &Set{
		Dependencies: make(map[string]*Dependency),
	}
}

func (s Set) GetDependencyOf(pkg *Package) *Dependency {
	p := pkg.Name()
	for root, dep := range s.Dependencies {
		if strings.HasPrefix(p, root) {
			return dep
		}
	}
	return nil
}

func (s Set) ContainsDependencyOf(pkg *Package) bool {
	p := pkg.Path()
	for root := range s.Dependencies {
		if strings.HasPrefix(p, root) {
			return true
		}
	}
	return false
}

func (s *Set) Add(imp *ast.ImportSpec) (pkg *Package, added bool, err error) {
	pkg, err = GetPackageFromImport(imp)
	if err != nil || pkg.IsRoot() {
		return
	}
	if dep := s.GetDependencyOf(pkg); dep != nil {
		added = dep.Add(pkg)
	} else {
		var dep *Dependency
		dep, err = New(pkg)
		if err != nil {
			return
		}
		s.Dependencies[dep.RootPackage] = dep
		added = true
	}
	return
}

func (s Set) GetDependencies() []*Dependency {
	deps := make([]*Dependency, len(s.Dependencies))
	idx := 0
	for _, dep := range s.Dependencies {
		deps[idx] = dep
		idx++
	}
	return deps
}
