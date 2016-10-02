package project

import (
	"path/filepath"

	"github.com/fische/gaoler/pkg"
	"github.com/fische/gaoler/project/dependency"
)

type Project struct {
	Root   string `json:"-" yaml:"-"`
	Vendor string `json:"-" yaml:"-"`

	Name string

	Dependencies *dependency.Set
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
		Dependencies: dependency.NewSet(),
	}, nil
}
