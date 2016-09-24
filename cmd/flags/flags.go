package flags

import "fmt"

type Flag struct {
	cb       func(v string) error
	dflt     string
	boolFlag bool
}

func New(cb func(v string) error, dflt interface{}, boolFlag bool) *Flag {
	return &Flag{
		cb:       cb,
		dflt:     fmt.Sprintf("%v", dflt),
		boolFlag: boolFlag,
	}
}

func (f *Flag) Set(v string) error {
	return f.cb(v)
}

func (f Flag) String() string {
	return f.dflt
}

func (f Flag) IsBoolFlag() bool {
	return f.boolFlag
}
