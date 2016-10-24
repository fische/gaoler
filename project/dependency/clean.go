package dependency

import (
	"fmt"
	"os"
	"path/filepath"
)

type CleanOption uint8

const (
	Pass CleanOption = 1 << iota
	Skip
	Remove
)

func (opt CleanOption) filter(path string, info os.FileInfo) error {
	if opt == Skip {
		if info.IsDir() {
			return filepath.SkipDir
		}
	} else if opt == Remove {
		if err := os.RemoveAll(path); err != nil {
			return err
		} else if info.IsDir() {
			return filepath.SkipDir
		}
	}
	return nil
}

func (d Dependency) CleanVendor(vendorRoot string, checkers ...func(info os.FileInfo) CleanOption) error {
	root := filepath.Clean(fmt.Sprintf("%s/%s/", vendorRoot, d.rootPackage))
	removeRootFiles := !d.Has(d.rootPackage)
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err == nil {
			for _, checker := range checkers {
				opt := checker(info)
				if opt == Pass {
					continue
				}
				return opt.filter(path, info)
			}
			if info.IsDir() {
				var rel string
				if rel, err = filepath.Rel(vendorRoot, path); err == nil && rel != d.rootPackage && !d.Has(rel) {
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
