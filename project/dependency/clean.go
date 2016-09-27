package dependency

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type CleanOption uint8

const (
	Pass CleanOption = 1 << iota
	SkipDir
	Remove
)

func RemoveTestFiles(info os.FileInfo) CleanOption {
	if !info.IsDir() {
		if strings.HasSuffix(info.Name(), "_test.go") {
			return Remove
		}
	} else if info.Name() == "testdata" {
		return Remove
	}
	return Pass
}

func KeepTestFiles(info os.FileInfo) CleanOption {
	if info.IsDir() && info.Name() == "testdata" {
		return SkipDir
	}
	return Pass
}

func (d Dependency) CleanVendor(vendorRoot string, checkers ...func(info os.FileInfo) CleanOption) error {
	root := filepath.Clean(fmt.Sprintf("%s/%s/", vendorRoot, d.rootPackage))
	removeRootFiles := !d.HasPackage(d.rootPackage)
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err == nil {
			for _, checker := range checkers {
				if opt := checker(info); opt == Pass {
					continue
				} else if opt == SkipDir && info.IsDir() {
					return filepath.SkipDir
				} else if opt == Remove {
					return os.RemoveAll(path)
				}
			}
			if info.IsDir() {
				var rel string
				if rel, err = filepath.Rel(vendorRoot, path); err == nil && rel != d.rootPackage && !d.HasPackage(rel) {
					if err = os.RemoveAll(path); err != nil {
						return err
					}
					return filepath.SkipDir
				}
			} else if removeRootFiles && path == filepath.Clean(fmt.Sprintf("%s/%s", root, info.Name())) {
				err = os.Remove(path)
			}
		}
		return err
	})
}
