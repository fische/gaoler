package project

import (
	"os"
	"path/filepath"
)

func isValidFile(path string) bool {
	return filepath.Ext(path) == ".go"
}

func isValidDir(path string) bool {
	name := filepath.Base(path)
	return name != ".git" && name != "vendor"
}

func skipDir(_ string) bool {
	return false
}

func sendError(ch chan<- error) func(path string, err error) error {
	return func(path string, err error) error {
		ch <- NewErrorMessage(err).WithField("path", path).WithMessage("Could not walk through file/folder")
		return err
	}
}

func walkError(_ string, err error) error {
	return err
}

func walk(p string, fileCond func(filepath string) bool, dirCond func(dirpath string) bool, errCond func(path string, err error) error) <-chan string {
	filepaths := make(chan string)
	go func() {
		err := filepath.Walk(p, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				if errCond != nil {
					return errCond(path, err)
				}
				return err
			}
			if !info.IsDir() && (fileCond == nil || fileCond(path)) {
				filepaths <- path
			} else if info.IsDir() && dirCond != nil && !dirCond(path) {
				return filepath.SkipDir
			}
			return nil
		})
		if err != nil {
			return
		}
		close(filepaths)
	}()
	return filepaths
}
