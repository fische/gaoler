package project

import (
	"os"
	"path/filepath"
)

//Walk walks through the whole project sending path of the files
//that satisfy the given `fileCond`, through the returned channel
//If `fileCond` is nil, all project's files will be sent.
//If `dirCond` is defined, it should return whether Walk should walks
//through this directory.
//If `errCond` is defined, the returned error is returned to `filepath.Walk`
//stopping the process. If it's not defined, in case of an error,
//the process will stop.
func (p Project) Walk(fileCond func(filepath string) bool, dirCond func(dirpath string) bool, errCond func(path string, err error) error) <-chan string {
	filepaths := make(chan string)
	go func() {
		err := filepath.Walk(p.Root, func(path string, info os.FileInfo, err error) error {
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
