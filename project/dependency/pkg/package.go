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
	srcDirs = build.Default.SrcDirs()
)

func GetFromPath(packagePath string, ignoreVendor bool) (*Package, error) {
	flags := build.FindOnly | build.AllowBinary
	if ignoreVendor {
		flags |= build.IgnoreVendor
	}
	p, err := build.Import(packagePath, "", flags)
	if err != nil {
		return nil, err
	}
	return &Package{
		Path: p.ImportPath,
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
