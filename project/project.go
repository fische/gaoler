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
	Root   string `json:"-"`
	Vendor string `json:"-"`

	Name string

	Dependencies dependency.Set
}

func GetProjectRootFromDir(dir string) string {
	return dir
}

func New(root string, ignoreVendor bool) (*Project, error) {
	name, err := pkg.GetPackagePath(root)
	if err != nil {
		return nil, err
	}
	p := &Project{
		Root:         root,
		Vendor:       filepath.Clean(root + "/vendor/"),
		Name:         name,
		Dependencies: make(dependency.Set),
	}
	if err = p.listPackages([]string{root}, nil, ignoreVendor); err != nil {
		return nil, err
	}
	return p, nil
}

func OpenCurrent(ignoreVendor bool) (*Project, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return New(GetProjectRootFromDir(wd), ignoreVendor)
}

func noTest(file os.FileInfo) bool {
	return !strings.HasSuffix(file.Name(), "_test.go")
}

func (project Project) listPackages(directories []string, fset *token.FileSet, ignoreVendor bool) error {
	if fset == nil {
		fset = token.NewFileSet()
	}
	var nextDirectories []string
	for _, dir := range directories {
		pkgs, err := parser.ParseDir(fset, dir, noTest, parser.ImportsOnly)
		if err != nil {
			return err
		}
		for _, p := range pkgs {
			for _, file := range p.Files {
				for _, imp := range file.Imports {
					if pkg.IsPseudoPackage(imp) {
						continue
					} else if p, err := pkg.GetFromImport(imp, ignoreVendor); err != nil {
						return err
					} else if added, err := project.Dependencies.Add(p, ignoreVendor); err != nil {
						return err
					} else if added && !p.Root {
						nextDirectories = append(nextDirectories, p.Dir)
					}
				}
			}
		}
	}
	if len(nextDirectories) > 0 {
		return project.listPackages(nextDirectories, fset, ignoreVendor)
	}
	return nil
}

func (project Project) IsDependency(d *dependency.Dependency) bool {
	return d.RootPackage == project.Name ||
		strings.HasPrefix(d.RootPackage, project.Name)
}
