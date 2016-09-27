package config

import "os"

type Flags uint8

const (
	Load Flags = 1 << iota
	Save
)

func (f Flags) OpenFlags() int {
	if f.Has(Load) {
		if f.Has(Save) {
			return os.O_RDWR | os.O_CREATE
		}
		return os.O_RDONLY
	} else if f.Has(Save) {
		return os.O_WRONLY | os.O_CREATE
	}
	return 0
}

func (f Flags) Has(other Flags) bool {
	return f&other == other
}
