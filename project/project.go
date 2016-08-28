package project

import (
	"go/parser"
	"go/token"
	"os"
	"strings"

	"github.com/fische/gaoler/pkg"
)

type Project struct {
	Root string
}

func GetProjectRootFromDir(dir string) string {
	return dir
}

func New(root string) *Project {
	return &Project{
		Root: root,
	}
}

func OpenCurrent() (*Project, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return &Project{
		Root: GetProjectRootFromDir(wd),
	}, nil
}

func noVendor(file os.FileInfo) bool {
	return !strings.HasSuffix(file.Name(), "_test.go")
}

func listPackages(directories []string, packages *pkg.Set, fset *token.FileSet) ([]*pkg.Package, error) {
	if packages == nil {
		packages = pkg.NewSet()
	}
	if fset == nil {
		fset = token.NewFileSet()
	}
	var nextDirectories []string
	for _, dir := range directories {
		pkgs, err := parser.ParseDir(fset, dir, noVendor, parser.ImportsOnly)
		if err != nil {
			return nil, err
		}
		for _, p := range pkgs {
			for _, file := range p.Files {
				for _, imp := range file.Imports {
					if pkg.IsPseudoPackage(imp) {
						continue
					} else if item, added, err := packages.Add(imp); err != nil {
						return nil, err
					} else if added && !item.IsRoot() {
						nextDirectories = append(nextDirectories, item.Path())
					}
				}
			}
		}
	}
	if len(nextDirectories) > 0 {
		return listPackages(nextDirectories, packages, fset)
	}
	return packages.GetPackages(), nil
}

func (p Project) ListDependencies() ([]*pkg.Package, error) {
	return listPackages([]string{p.Root}, nil, nil)
}
