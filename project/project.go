package project

import (
	"path/filepath"

	"github.com/fische/gaoler/pkg"
	"github.com/fische/gaoler/pkg/set"
)

type Project struct {
	root   string
	vendor string

	Name string

	Dependencies Dependencies
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
		Dependencies: make(Dependencies),
	}, nil
}

func NewWithDependencies(root string, s *set.Packages, openRepository bool) (*Project, error) {
	p, err := New(root)
	if err != nil {
		return nil, err
	}
	p.Dependencies.Set(s, openRepository)
	return p, nil
}

func (p Project) Root() string {
	return p.root
}

func (p Project) Vendor() string {
	return p.vendor
}
