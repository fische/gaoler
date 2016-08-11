package project

import (
	"go/ast"
	"go/build"

	"github.com/fische/gaoler/vcs"
	"github.com/fische/gaoler/vcs/modules"
)

//Dependency represents a Go package import
type Dependency struct {
	Package *build.Package
	Name    string
}

//NewDependency creates a new `Dependency` and retrieves `build.Package` from `i`
func NewDependency(i *ast.ImportSpec) (*Dependency, error) {
	name := GetName(i)
	p, err := build.Import(name, "", build.FindOnly)
	if err != nil {
		return nil, err
	}
	return &Dependency{
		Package: p,
		Name:    name,
	}, nil
}

//GetRepository returns appropriate `vcs.Repository` of the `Dependency`
func (d Dependency) GetRepository() (vcsName string, repo vcs.Repository, err error) {
	return modules.GetRepository(d.Package.Dir)
}
