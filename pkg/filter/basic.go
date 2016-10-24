package filter

import (
	"os"
	"path/filepath"
	"strings"
)

type Filter struct {
	test       bool
	vendorPath string
}

func New(test bool, srcPath string) Factory {
	return &Filter{
		test:       test,
		vendorPath: filepath.Clean(srcPath + "/vendor/"),
	}
}

func (f Filter) New(dir string) func(info os.FileInfo) bool {
	return func(info os.FileInfo) bool {
		test := strings.HasSuffix(info.Name(), "_test.go")
		return ((f.test && !(test && strings.HasPrefix(dir, f.vendorPath))) || !test)
	}
}
