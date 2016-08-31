package dependency

import (
	"errors"
	"go/ast"
	"go/build"
	"path/filepath"
	"strings"
)

type Package struct {
	*build.Package
}

var (
	srcDirs = build.Default.SrcDirs()
)

func GetPackageFromPath(packagePath string) (*Package, error) {
	var (
		err error

		p = &Package{}
	)
	p.Package, err = build.Import(packagePath, "", build.FindOnly|build.IgnoreVendor|build.AllowBinary)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func GetPackageFromImport(imp *ast.ImportSpec) (*Package, error) {
	return GetPackageFromPath(GetNameFromImport(imp))
}

func GetPackagePathFromPath(dir string) (string, error) {
	for _, src := range srcDirs {
		if strings.HasPrefix(dir, src) {
			rel, err := filepath.Rel(src, dir)
			if err != nil {
				return "", err
			}
			return rel, nil
		}
	}
	return "", errors.New("Could not find package in src directories")
}

func (p Package) Name() string {
	return p.ImportPath
}

func (p Package) Path() string {
	return p.Dir
}

func (p Package) IsRoot() bool {
	return p.Goroot
}
