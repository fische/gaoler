package pkg

import (
	"errors"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"

	"github.com/fische/gaoler/pkg/filter"
)

type Set struct {
	packages    map[string]*Package
	Filter      filter.Factory
	Constructor Factory
}

func NewSet() *Set {
	return &Set{
		packages: make(map[string]*Package),
	}
}

func (s *Set) listDir(dir string, fset *token.FileSet, filter func(info os.FileInfo) bool) ([]string, error) {
	var nextDirectories []string
	pkgs, err := parser.ParseDir(fset, dir, filter, parser.ImportsOnly)
	if err != nil {
		return nil, err
	}
	for _, p := range pkgs {
		for _, file := range p.Files {
			for _, imp := range file.Imports {
				packagePath := GetPackagePathFromImport(imp)
				if !s.Has(packagePath) {
					add, dirs, err := s.Constructor.New(packagePath)
					if err != nil {
						return nil, err
					} else if add != nil {
						s.Insert(add, true)
					}
					if len(dirs) > 0 {
						nextDirectories = append(nextDirectories, dirs...)
					}
				}
			}
		}
	}
	return nextDirectories, nil
}

func (s *Set) list(directories []string, fset *token.FileSet) error {
	var nextDirectories []string
	for _, dir := range directories {
		if filepath.Base(dir) != "testdata" {
			var filter func(info os.FileInfo) bool
			if s.Filter != nil {
				filter = s.Filter.New(dir)
			}
			next, err := s.listDir(dir, fset, filter)
			if err != nil {
				return err
			}
			nextDirectories = append(nextDirectories, next...)
		}
	}
	if len(nextDirectories) > 0 {
		return s.list(nextDirectories, fset)
	}
	return nil
}

func (s *Set) ListFrom(srcPath string) error {
	if s.Constructor == nil {
		return errors.New("Could not complete set without a package constructor")
	}
	return s.list([]string{srcPath}, token.NewFileSet())
}

func (s *Set) Insert(p *Package, force bool) (added bool) {
	if !force {
		if s.Has(p.path) {
			return
		}
	}
	added = true
	s.packages[p.path] = p
	return
}

func (s *Set) Remove(key string) {
	delete(s.packages, key)
}

func (s Set) Has(packagePath string) bool {
	_, ok := s.packages[packagePath]
	return ok
}

func (s Set) Packages() map[string]*Package {
	return s.packages
}

func (s Set) IsVendored() bool {
	for _, p := range s.Packages() {
		if !p.IsVendored() {
			return false
		}
	}
	return true
}

func (s Set) IsSaved() bool {
	for _, p := range s.Packages() {
		if !p.IsSaved() {
			return false
		}
	}
	return true
}
