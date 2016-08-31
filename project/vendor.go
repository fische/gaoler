package project

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fische/gaoler/project/dependency"
)

func CleanVendor(vendorPath string, dependencies []*dependency.Dependency) error {
	for _, dep := range dependencies {
		root := filepath.Clean(fmt.Sprintf("%s/%s/", vendorPath, dep.RootPackage))
		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				rel, err := filepath.Rel(vendorPath, path)
				if err != nil {
					return err
				}
				if !(dep.HasPackage(rel)) && rel != dep.RootPackage {
					err := os.RemoveAll(path)
					if err != nil {
						return err
					}
					return filepath.SkipDir
				}
				//TODO Remove files from root
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}
