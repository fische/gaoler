package project

import (
	"errors"
	"go/parser"
	"go/token"
	"path/filepath"

	"github.com/fische/gaoler/project/dependency/pkg"
)

func isProjectRoot(dir string) bool {
	pkgs, err := parser.ParseDir(token.NewFileSet(), dir, nil, parser.ImportsOnly)
	if err != nil {
		return false
	}
	_, ok := pkgs["main"]
	return ok
}

func getProjectRootFromDir(dir string) (string, error) {
	if isProjectRoot(dir) {
		return dir, nil
	}
	next := filepath.Dir(dir)
	if pkg.IsInSrcDirs(next) {
		return getProjectRootFromDir(next)
	}
	return "", errors.New("Could not find root directory")
}

func GetProjectRootFromDir(dir string) (string, error) {
	if !pkg.IsInSrcDirs(dir) {
		return "", errors.New("This directory is not located in the GO environment")
	}
	return getProjectRootFromDir(dir)
}
