package project

import (
	"go/parser"
	"go/token"
	"os"
	"strings"

	"github.com/fische/gaoler/project/dependency"
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

//TODO Do not add dependencies to packages of the project itself
func listPackages(directories []string, dependencies *dependency.Set, fset *token.FileSet) ([]*dependency.Dependency, error) {
	if dependencies == nil {
		dependencies = dependency.NewSet()
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
					if dependency.IsPseudoPackage(imp) {
						continue
					} else if item, added, err := dependencies.Add(imp); err != nil {
						return nil, err
					} else if added && !item.IsRoot() {
						nextDirectories = append(nextDirectories, item.Path())
					}
				}
			}
		}
	}
	if len(nextDirectories) > 0 {
		return listPackages(nextDirectories, dependencies, fset)
	}
	return dependencies.GetDependencies(), nil
}

func (p Project) ListDependencies() ([]*dependency.Dependency, error) {
	return listPackages([]string{p.Root}, nil, nil)
}
