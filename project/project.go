package project

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/fische/gaoler/project/dependency"
	"github.com/fische/gaoler/project/dependency/pkg"
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

func (p Project) listPackages(directories []string, dependencies dependency.Set, fset *token.FileSet, ignoreVendor bool) (dependency.Set, error) {
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
					if pkg.IsPseudoPackage(imp) {
						continue
					} else if p, err := pkg.GetFromImport(imp, ignoreVendor); err != nil {
						return nil, err
					} else if added, err := dependencies.Add(p, ignoreVendor); err != nil {
						return nil, err
					} else if added && !p.Root {
						nextDirectories = append(nextDirectories, p.Dir)
					}
				}
			}
		}
	}
	if len(nextDirectories) > 0 {
		return p.listPackages(nextDirectories, dependencies, fset, ignoreVendor)
	}
	return dependencies, nil
}

func (p Project) ListDependencies(ignoreVendor bool) (dependency.Set, error) {
	return p.listPackages([]string{p.Root}, nil, nil, ignoreVendor)
}

func (p Project) IsDependency(d *dependency.Dependency) bool {
	path, err := d.Repository.GetPath()
	if err != nil {
		return false
	}
	return path == p.Root
}
