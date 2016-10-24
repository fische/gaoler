package filter

import "os"

type Factory interface {
	New(dir string) func(info os.FileInfo) bool
}
