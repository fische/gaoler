package pkg

import (
	"errors"
	"path/filepath"
	"strings"
)

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
