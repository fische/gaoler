package project

import (
	"path/filepath"

	"github.com/fische/gaoler/project/dependency"
)

type Project struct {
	rootPath   string
	vendorPath string

	*dependency.Set
}

func New(rootPath string) *Project {
	return &Project{
		rootPath:   rootPath,
		vendorPath: filepath.Clean(rootPath + "/vendor/"),
		Set:        dependency.NewSet(),
	}
}

func (p Project) RootPath() string {
	return p.rootPath
}

func (p Project) VendorPath() string {
	return p.vendorPath
}
