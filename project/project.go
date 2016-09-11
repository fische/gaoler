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
	Root   string `json:"-" yaml:"-"`
	Vendor string `json:"-" yaml:"-"`

	Name string

	Dependencies dependency.Set
}

func New(root string) (*Project, error) {
	name, err := pkg.GetPackagePath(root)
	if err != nil {
		return nil, err
	}
	return &Project{
		Root:         root,
		Vendor:       filepath.Clean(root + "/vendor/"),
		Name:         name,
		Dependencies: make(dependency.Set),
	}, nil
}

func NewWithDependencies(root string, keepTests, ignoreVendor bool) (*Project, error) {
	p, err := New(root)
	if err != nil {
		return nil, err
	}
	p.Dependencies = make(dependency.Set)
	var filters []func(file os.FileInfo) bool
	if !keepTests {
		filters = append(filters, noTest)
	}
	if err = p.ListPackages(ignoreVendor, filters...); err != nil {
		return nil, err
	}
	return p, nil
}

func OpenCurrent(keepTests, ignoreVendor bool) (*Project, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	dir, err := GetProjectRootFromDir(wd)
	if err != nil {
		return nil, err
	}
	return NewWithDependencies(dir, keepTests, ignoreVendor)
}

func noTest(file os.FileInfo) bool {
	return !strings.HasSuffix(file.Name(), "_test.go")
}

func GatherFilters(filters ...func(file os.FileInfo) bool) func(file os.FileInfo) bool {
	if len(filters) == 0 {
		return nil
	}
	return func(file os.FileInfo) bool {
		for _, filter := range filters {
			if !filter(file) {
				return false
			}
		}
		return true
	}
}

func (project *Project) listPackages(directories []string, fset *token.FileSet, ignoreVendor bool, filter func(file os.FileInfo) bool) error {
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
					} else if p.Root {
						continue
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
		return project.listPackages(nextDirectories, fset, ignoreVendor, filter)
	}
	return nil
}

func (project *Project) ListPackages(ignoreVendor bool, filters ...func(file os.FileInfo) bool) error {
	return project.listPackages([]string{project.Root}, nil, ignoreVendor, GatherFilters(filters...))
}

func (project Project) IsDependency(d *dependency.Dependency) bool {
	return d.RootPackage == project.Name ||
		strings.HasPrefix(d.RootPackage, project.Name)
}
