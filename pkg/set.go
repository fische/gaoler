package pkg

import (
	"go/parser"
	"go/token"
	"os"
)

type Set struct {
	Packages map[string]*Package
	Filter   func(info os.FileInfo) bool
	OnAdded  func(p *Package) (nextDirectory string, err error)
}

func NewSet() *Set {
	return &Set{
		Packages: make(map[string]*Package),
	}
}

func (s *Set) ListFrom(srcPath string) error {
	return s.list([]string{srcPath}, token.NewFileSet())
}

func (s *Set) list(directories []string, fset *token.FileSet) error {
	var nextDirectories []string
	for _, dir := range directories {
		pkgs, err := parser.ParseDir(fset, dir, s.Filter, parser.ImportsOnly)
		if err != nil {
			return err
		}
		for _, p := range pkgs {
			for _, file := range p.Files {
				for _, imp := range file.Imports {
					n := NewFromImport(imp)
					if added := s.Insert(n, false); added && s.OnAdded != nil {
						if next, err := s.OnAdded(n); err != nil {
							return err
						} else if next != "" {
							nextDirectories = append(nextDirectories, next)
						}
					}
				}
			}
		}
	}
	if len(nextDirectories) > 0 {
		return s.list(nextDirectories, fset)
	}
	return nil
}

func (s *Set) Insert(p *Package, force bool) (added bool) {
	if !force {
		if _, ok := s.Packages[p.Path()]; ok {
			return
		}
	}
	added = true
	s.Packages[p.Path()] = p
	return
}

func (s *Set) ForEach(cb func(key string, value *Package) error) error {
	for k, v := range s.Packages {
		if err := cb(k, v); err != nil {
			return err
		}
	}
	return nil
}

func (s *Set) Remove(key string) {
	delete(s.Packages, key)
}
