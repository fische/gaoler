package pkg

import (
	"os"
	"strings"
)

type Flags int

const (
	NoTests Flags = 1 << iota
	NoPseudoPackage
	NoStandardPackage
	NoLocalPackage
	IgnoreVendor
)

func (f Flags) Filter(info os.FileInfo) bool {
	if f.Has(NoTests) && strings.HasSuffix(info.Name(), "_test.go") {
		return false
	}
	return true
}

func (f Flags) Has(flags Flags) bool {
	return f&flags == flags
}
