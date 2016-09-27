package pkg

import (
	"os"
	"strings"
)

const (
	testExtension = "_test.go"
)

func NoTestFiles(info os.FileInfo) bool {
	return !strings.HasSuffix(info.Name(), testExtension)
}
