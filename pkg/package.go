package pkg

import (
	"errors"
	"go/ast"
	"go/build"
	"path/filepath"
	"strings"
)

type Package struct {
	dir  string
	path string

	root     bool
	saved    bool
	vendored bool
	// TODO
	// importedBy []string
}

var (
	srcDirs = build.Default.SrcDirs()
)

func New(packagePath, srcPath string, flags Flags) (*Package, error) {
	p := &Package{
		path: packagePath,
	}
	if err := p.Import(srcPath, flags); err != nil {
		return nil, err
	}
	return p, nil
}

func NewFromImport(imp *ast.ImportSpec, srcPath string, flags Flags) (*Package, error) {
	return New(GetPackagePathFromImport(imp), srcPath, flags)
}

func GetPackagePath(dir string) (string, error) {
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

func IsInSrcDirs(dir string) bool {
	for _, src := range srcDirs {
		if strings.HasPrefix(dir, src) {
			return true
		}
	}
	return false
}

func (p *Package) Import(srcPath string, flags Flags) error {
	f := build.FindOnly | build.AllowBinary
	ignoreVendor := flags.Has(IgnoreVendor)
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

func (p Package) IsLocal(srcPath string) bool {
	return !p.vendored && strings.HasPrefix(p.dir, srcPath)
}
