package pkg

import (
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
}

func New(packagePath string) *Package {
	return &Package{
		path: packagePath,
	}
}

func Import(packagePath, srcPath string, ignoreVendor bool) (*Package, error) {
	p := New(packagePath)
	if err := p.Import(srcPath, ignoreVendor); err != nil {
		return nil, err
	}
	return p, nil
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

func (p Package) Dir() string {
	return p.dir
}

func (p Package) Path() string {
	return p.path
}

func (p Package) IsStandardPackage() bool {
	return p.root
}

func (p Package) IsLocal() bool {
	return p.local
}

func (p Package) IsSaved() bool {
	return p.saved
}

func (p Package) IsVendored() bool {
	return p.vendored
}

func (p Package) IsPseudoPackage() bool {
	return IsPseudoPackage(p.path)
}
