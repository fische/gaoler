package project

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/fische/gaoler/project/dependency"
)

type Project struct {
	Root   string
	Vendor string
}

func GetProjectRootFromDir(dir string) string {
	return dir
}

func New(root string) *Project {
	return &Project{
		Root:   root,
		Vendor: filepath.Clean(root + "/vendor/"),
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

func (p Project) listPackages(directories []string, dependencies *dependency.Set, fset *token.FileSet) ([]*dependency.Dependency, error) {
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
		for _, pkg := range pkgs {
			for _, file := range pkg.Files {
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
		return p.listPackages(nextDirectories, dependencies, fset)
	}
	return dependencies.GetDependencies(), nil
}

func (p Project) ListDependencies() ([]*dependency.Dependency, error) {
	return p.listPackages([]string{p.Root}, nil, nil)
}

func (p Project) HasLocalDependency(d *dependency.Dependency) bool {
	path, err := d.Repository.GetPath()
	if err != nil {
		return false
	}
	return path == p.Root
}
