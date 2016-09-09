package pkg

import (
	"errors"
	"go/ast"
	"go/build"
	"path/filepath"
	"strings"
)

type Package struct {
	Dir  string
	Path string
	Root bool
}

var (
	srcDirs    = build.Default.SrcDirs()
	srcPath    = ""
	vendorPath = ""
)

func SetSourcePath(p string) {
	srcPath = p
	vendorPath = filepath.Clean(p + "/vendor/")
}

func GetFromPath(packagePath string, ignoreVendor bool) (*Package, error) {
	flags := build.FindOnly | build.AllowBinary
	if ignoreVendor {
		flags |= build.IgnoreVendor
	}
	p, err := build.Import(packagePath, srcPath, flags)
	if err != nil {
		return nil, err
	}
	return &Package{
		Path: packagePath,
		Dir:  p.Dir,
		Root: p.Goroot,
	}, nil
}

func GetFromImport(imp *ast.ImportSpec, ignoreVendor bool) (*Package, error) {
	return GetFromPath(GetNameFromImport(imp), ignoreVendor)
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

func (p Package) MarshalJSON() ([]byte, error) {
	return []byte("\"" + p.Path + "\""), nil
}

func (p *Package) UnmarshalJSON(data []byte) error {
	if data[0] == '"' && data[len(data)-1] == '"' {
		p.Path = string(data[1 : len(data)-1])
		i, err := build.Import(p.Path, srcPath, build.FindOnly|build.AllowBinary)
		if err == nil {
			p.Dir = i.Dir
			p.Root = i.Goroot
		}
		return nil
	}
	return errors.New("Could not unmarshal package")
}

func (p Package) IsVendored() bool {
	if len(vendorPath) == 0 {
		return false
	}
	return strings.HasPrefix(p.Dir, vendorPath)
}
