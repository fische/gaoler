package pkg

import (
	"go/ast"
	"go/build"
)

type Package struct {
	*build.Package
}

func GetPackageFromImport(imp *ast.ImportSpec) (*Package, error) {
	var (
		err error

		p = &Package{}
	)
	p.Package, err = build.Import(GetNameFromImport(imp), "", build.FindOnly|build.IgnoreVendor|build.AllowBinary)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func New(packagePath string) (*Package, error) {
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

func (p *Package) Name() string {
	return p.ImportPath
}

func (p *Package) Path() string {
	return p.Dir
}

func (p *Package) IsRoot() bool {
	return p.Goroot
}
