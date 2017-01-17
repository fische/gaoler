package cmd

import (
	"io"
	"os"
	"path/filepath"
	"regexp"
)

func resetDirectory(path string) error {
	if err := os.RemoveAll(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return os.Mkdir(path, 0775)
}

func checkRegexps(s string, regexps []*regexp.Regexp) bool {
	for _, r := range regexps {
		if len(r.FindAllString(s, -1)) > 0 {
			return true
		}
	}
	return false
}

func removeEmptyParents(basePath, endPath string) error {
	if err := os.RemoveAll(basePath); err != nil {
		return err
	} else if basePath != endPath {
		next := filepath.Dir(basePath)
		if empty, err := isEmptyDirectory(next); err != nil {
			return err
		} else if empty {
			return removeEmptyParents(next, endPath)
		}
	}
	return nil
}

func isEmptyDirectory(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()
	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}
