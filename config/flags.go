package config

import "os"

type Flags int

const (
	LoadFlag Flags = 1 << iota
	SaveFlag
)

func (f Flags) OpenFlags() int {
	var ret int
	if f.Has(LoadFlag) {
		if f.Has(SaveFlag) {
			ret |= os.O_RDWR | os.O_CREATE
		} else {
			ret |= os.O_RDONLY
		}
	} else if f.Has(SaveFlag) {
		ret |= os.O_WRONLY | os.O_CREATE
	}
	return ret
}

func (f Flags) Has(s Flags) bool {
	return f&s == s
}
