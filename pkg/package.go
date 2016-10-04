package pkg

import (
	"go/ast"
	"go/build"
	"path/filepath"
	"strings"
)

type Package struct {
	Dir  string
	Path string

	Root     bool
	Local    bool
	Saved    bool
	Vendored bool

	// TODO: importedBy []string
}

func New(packagePath string) *Package {
	return &Package{
		Path: packagePath,
	}
}

func NewFromImport(imp *ast.ImportSpec) *Package {
	return New(GetPackagePathFromImport(imp))
}

func (p *Package) Import(srcPath string, ignoreVendor bool) error {
	f := build.FindOnly | build.AllowBinary
	if ignoreVendor {
		f |= build.IgnoreVendor
	}
	imp, err := build.Import(p.Path, srcPath, f)
	if err != nil {
		return err
	}
	p.Dir = imp.Dir
	p.Root = imp.Goroot
	if !ignoreVendor {
		p.Vendored = strings.HasPrefix(p.Dir, filepath.Clean(srcPath+"/vendor/"))
	}
	p.Local = !p.Vendored && strings.HasPrefix(p.Dir, srcPath)
	return nil
}

func (p Package) IsPseudoPackage() bool {
	return IsPseudoPackage(p.Path)
}
