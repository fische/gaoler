package set

import (
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/fische/gaoler/pkg"
	"github.com/fische/gaoler/project/dependency"
)

type Packages struct {
	packages map[string]*pkg.Package
}

func NewPackages() *Packages {
	return &Packages{
		packages: make(map[string]*pkg.Package),
	}
}

func ListPackagesFrom(srcPath string, flags pkg.Flags) (*Packages, error) {
	s := NewPackages()
	if err := s.list([]string{srcPath}, srcPath, flags, token.NewFileSet()); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Packages) list(directories []string, srcPath string, flags pkg.Flags, fset *token.FileSet) error {
	var nextDirectories []string
	for _, dir := range directories {
		pkgs, err := parser.ParseDir(fset, dir, flags.Filter, parser.ImportsOnly)
		if err != nil {
			return err
		}
		for _, p := range pkgs {
			for _, file := range p.Files {
				for _, imp := range file.Imports {
					if p, added, err := s.InsertImport(imp, srcPath, flags); err != nil {
						return err
					} else if p != nil && (added || (flags.Has(pkg.NoLocalPackage) && p.IsLocal(srcPath))) {
						nextDirectories = append(nextDirectories, p.Dir())
					}
				}
			}
		}
	}
	if len(nextDirectories) > 0 {
		return s.list(nextDirectories, srcPath, flags, fset)
	}
	return nil
}

func (s *Packages) CompleteFrom(srcPath string, flags pkg.Flags) error {
	return s.list([]string{srcPath}, srcPath, flags, token.NewFileSet())
}

func (s *Packages) InsertImport(imp *ast.ImportSpec, srcPath string, flags pkg.Flags) (p *pkg.Package, added bool, err error) {
	if flags.Has(pkg.NoPseudoPackage) && pkg.IsPseudoPackage(pkg.GetPackagePathFromImport(imp)) {
		return
	}
	packagePath := pkg.GetPackagePathFromImport(imp)
	var ok bool
	if p, ok = s.packages[packagePath]; ok {
		return
	}
	p, err = pkg.New(packagePath, srcPath, flags)
	if err != nil {
		return
	} else if flags.Has(pkg.NoStandardPackage) && p.IsStandardPackage() {
		return
	} else if flags.Has(pkg.NoLocalPackage) && p.IsLocal(srcPath) {
		return
	}
	added = true
	s.packages[packagePath] = p
	return
}

func (s *Packages) ForceInsertImport(imp *ast.ImportSpec, srcPath string, flags pkg.Flags) (p *pkg.Package, err error) {
	packagePath := pkg.GetPackagePathFromImport(imp)
	p, err = pkg.New(packagePath, srcPath, flags)
	if err != nil {
		return
	}
	s.packages[packagePath] = p
	return
}

func (s *Packages) InsertDependency(dep *dependency.Dependency, srcPath string, flags pkg.Flags) (pkgs []*pkg.Package) {
	if flags.Has(pkg.NoPseudoPackage) && pkg.IsPseudoPackage(dep.RootPackage()) {
		return
	}
	for _, p := range dep.Packages {
		pkgPath := p.Path()
		if stored, ok := s.packages[pkgPath]; ok {
			pkgs = append(pkgs, stored)
			continue
		} else if flags.Has(pkg.NoStandardPackage) && p.IsStandardPackage() {
			continue
		} else if flags.Has(pkg.NoLocalPackage) && p.IsLocal(srcPath) {
			continue
		}
		s.packages[pkgPath] = p
		pkgs = append(pkgs, p)
	}
	return
}

func (s *Packages) ForceInsertDependency(dep *dependency.Dependency) (pkgs []*pkg.Package) {
	for _, p := range dep.Packages {
		pkgPath := p.Path()
		s.packages[pkgPath] = p
		pkgs = append(pkgs, p)
	}
	return
}

func (s *Packages) ForEach(cb func(key string, value *pkg.Package) error) error {
	for k, v := range s.packages {
		if err := cb(k, v); err != nil {
			return err
		}
	}
	return nil
}

func (s *Packages) Delete(key string) {
	delete(s.packages, key)
}

func (s Packages) Packages() map[string]*pkg.Package {
	return s.packages
}
