package dependency

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fische/gaoler/pkg"
)

type CleanOption uint8

const (
	Pass CleanOption = 1 << iota
	SkipDir
	Remove
)

func RemoveGoTestFiles(info os.FileInfo) CleanOption {
	if pkg.IsNotGoTestFile(info) {
		return Pass
	}
	return Remove
}

func KeepGoTestFiles(info os.FileInfo) CleanOption {
	if !pkg.IsNotGoTestFile(info) && info.IsDir() {
		return SkipDir
	}
	return Pass
}

func (d Dependency) CleanVendor(vendorRoot string, checkers ...func(info os.FileInfo) CleanOption) error {
	root := filepath.Clean(fmt.Sprintf("%s/%s/", vendorRoot, d.RootPackage))
	removeRootFiles := !d.HasPackage(d.RootPackage)
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err == nil {
			for _, checker := range checkers {
				if opt := checker(info); opt == Pass {
					continue
				} else if opt == SkipDir && info.IsDir() {
					return filepath.SkipDir
				} else if opt == Remove {
					if err = os.RemoveAll(path); err != nil {
						return err
					} else if info.IsDir() {
						return filepath.SkipDir
					}
					return nil
				}
			}
			if info.IsDir() {
				var rel string
				if rel, err = filepath.Rel(vendorRoot, path); err == nil && rel != d.RootPackage && !d.HasPackage(rel) {
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
