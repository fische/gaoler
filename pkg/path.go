package pkg

import (
	"errors"
	"go/ast"
	"go/build"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	srcDirs        = build.Default.SrcDirs()
	pseudoPackages = []string{
		"C",
	}
	vendorRule = regexp.MustCompile("vendor/")
)

func GetPackagePathFromImport(imp *ast.ImportSpec) string {
	return imp.Path.Value[1 : len(imp.Path.Value)-1]
}

func GetPackagePath(dir string) (string, error) {
	for _, src := range srcDirs {
		if strings.HasPrefix(dir, src) {
			rel, err := filepath.Rel(src, dir)
			if err != nil {
				return "", err
			}
			matchs := vendorRule.Split(rel, -1)
			return matchs[len(matchs)-1], nil
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

func IsPseudoPackage(pkgPath string) bool {
	for _, p := range pseudoPackages {
		if pkgPath == p {
			return true
		}
	}
	return false
}
