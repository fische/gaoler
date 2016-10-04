package pkg

import (
	"os"
	"strings"
)

const (
	testExtension = "_test.go"
	testdataDir   = "testdata"
)

func IsNotGoTestFile(info os.FileInfo) bool {
	return (!info.IsDir() && !strings.HasSuffix(info.Name(), testExtension)) ||
		(info.IsDir() && info.Name() != testdataDir)
}
