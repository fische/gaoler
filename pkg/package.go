package pkg

import (
	"go/ast"
	"go/build"
	"path/filepath"
	"strings"
)

type Package struct {
	dir  string
	path string

	root     bool
	local    bool
	saved    bool
	vendored bool

	// TODO: importedBy []string
}

var (
	srcDirs = build.Default.SrcDirs()
)

func New(packagePath string) *Package {
	return &Package{
		path: packagePath,
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
	imp, err := build.Import(p.path, srcPath, f)
	if err != nil {
		return err
	}
	p.dir = imp.Dir
	p.root = imp.Goroot
	if !ignoreVendor {
		p.vendored = strings.HasPrefix(p.dir, filepath.Clean(srcPath+"/vendor/"))
	}
	p.local = !p.vendored && strings.HasPrefix(p.dir, srcPath)
	return nil
}

func (p Package) Path() string {
	return p.path
}

func (p Package) Dir() string {
	return p.dir
}

func (p Package) IsPseudoPackage() bool {
	return IsPseudoPackage(p.path)
}

func (p Package) IsStandardPackage() bool {
	return p.root
}

func (p Package) IsSaved() bool {
	return p.saved
}

func (p Package) IsVendored() bool {
	return p.vendored
}

func (p Package) IsLocal() bool {
	return p.local
}
