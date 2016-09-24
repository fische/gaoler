package cmd

import (
	"os"

	"github.com/fische/gaoler/config"
	"github.com/fische/gaoler/pkg"
)

const (
	trueValue = "true"
)

func resetDir(dir string) error {
	err := os.RemoveAll(dir)
	if err != nil {
		return err
	}
	return os.MkdirAll(dir, 0775)
}

func setNoTestPkg(v string) error {
	if v == trueValue {
		pkgFlags ^= pkg.NoTests
	}
	return nil
}

func setConfigFlag(f *config.Flags, n config.Flags) func(v string) error {
	return func(v string) error {
		if v == trueValue {
			*f |= n
		}
		return nil
	}
}

func execMiddlewares(middlewares ...func()) func() {
	return func() {
		for _, f := range middlewares {
			f()
		}
	}
}
