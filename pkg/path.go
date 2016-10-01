package pkg

import (
	"errors"
	"go/build"
	"path/filepath"
	"strings"
)

var (
	srcDirs = build.Default.SrcDirs()
)

func GetPackagePath(dir string) (string, error) {
	for _, src := range srcDirs {
		if strings.HasPrefix(dir, src) {
			rel, _ := filepath.Rel(src, dir)
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
