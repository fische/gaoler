package project

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fische/gaoler/project/dependency"
)

func (p Project) CleanVendor(vendorPath string, dependencies []*dependency.Dependency) error {
	for _, dep := range dependencies {
		if !p.HasLocalDependency(dep) {
			root := filepath.Clean(fmt.Sprintf("%s/%s/", vendorPath, dep.RootPackage))
			removeRootFiles := !dep.HasPackage(dep.RootPackage)
			err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if info.IsDir() {
					var rel string
					rel, err = filepath.Rel(vendorPath, path)
					if err != nil {
						return err
					}
					if rel != dep.RootPackage && !dep.HasPackage(rel) {
						err = os.RemoveAll(path)
						if err != nil {
							return err
						}
						return filepath.SkipDir
					}
				} else if removeRootFiles && path == filepath.Clean(fmt.Sprintf("%s/%s", root, info.Name())) {
					err = os.Remove(path)
					if err != nil {
						return err
					}
				}
				return nil
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}
