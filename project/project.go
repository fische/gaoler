package project

import (
	"path/filepath"

	"github.com/fische/gaoler/pkg"
	"github.com/fische/gaoler/project/dependency"
)

type Project struct {
	root   string
	vendor string

	Name string

	Dependencies *dependency.Set
}

func New(root string) (*Project, error) {
	name, err := pkg.GetPackagePath(root)
	if err != nil {
		return nil, err
	}
	return &Project{
		root:         root,
		vendor:       filepath.Clean(root + "/vendor/"),
		Name:         name,
		Dependencies: dependency.NewSet(),
	}, nil
}

func (p Project) Root() string {
	return p.root
}

func (p Project) Vendor() string {
	return p.vendor
}
