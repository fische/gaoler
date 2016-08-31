package dependency

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fische/gaoler/vcs"
	"github.com/fische/gaoler/vcs/modules"
)

type CleanCheck func(info os.FileInfo) CleanOption
type CleanOption uint8

const (
	Pass CleanOption = 1 << iota
	SkipDir
	Remove
)

func RemoveTestFiles(info os.FileInfo) CleanOption {
	if strings.HasSuffix(info.Name(), "_test.go") {
		return Remove
	}
	return Pass
}

func (d Dependency) Vendor(vendorRoot string, checkers ...CleanCheck) error {
	path := filepath.Clean(fmt.Sprintf("%s/%s/", vendorRoot, d.RootPackage))
	v, _ := modules.GetVCS(d.Repository.GetVCSName())
	_, err := vcs.CloneRepository(v, path, d.Repository)
	if err != nil {
		return err
	}
	return d.CleanVendor(vendorRoot, checkers...)
}

func (d Dependency) CleanVendor(vendorRoot string, checkers ...CleanCheck) error {
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
					return os.RemoveAll(path)
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
