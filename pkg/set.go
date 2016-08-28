package pkg

import (
	"go/ast"
	"sync"
)

type Set struct {
	locker   *sync.RWMutex
	packages map[string]*Package
}

func (s *Set) Add(imp *ast.ImportSpec) (item *Package, added bool, err error) {
	s.locker.Lock()
	defer s.locker.Unlock()
	if p, ok := s.packages[GetNameFromImport(imp)]; ok {
		return p, false, nil
	}
	p, err := GetPackageFromImport(imp)
	if err != nil {
		return nil, false, err
	}
	s.packages[p.Name()] = p
	return p, true, nil
}

func (s Set) Get(packagePath string) (*Package, bool) {
	p, ok := s.packages[packagePath]
	return p, ok
}

func (s Set) GetPackages() []*Package {
	var (
		idx int

		arr = make([]*Package, len(s.packages))
	)
	for _, p := range s.packages {
		arr[idx] = p
		idx++
	}
	return arr
}

func NewSet() *Set {
	return &Set{
		locker:   new(sync.RWMutex),
		packages: make(map[string]*Package),
	}
}
