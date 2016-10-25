package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fische/gaoler/internal/cmd/middleware"
	"github.com/fische/gaoler/project"
	"github.com/jawher/mow.cli"
)

var (
	Gaoler = cli.App("goaler", "A Go package manager")
	ctx    = middleware.NewContext()

	ExitSuccess = 0
	ExitFailure = 1

	rootPath   *string
	configPath *string
)

func init() {
	var dir string
	if wd, err := os.Getwd(); err != nil {
		fmt.Fprintf(os.Stderr, "Could not get working directory : %v\n", err)
		cli.Exit(ExitFailure)
	} else if dir, err = project.GetProjectRootFromDir(wd); err != nil {
		dir = wd
	}
	rootPath = Gaoler.StringOpt("r root", dir, "Path to the root package")
	configPath = Gaoler.StringOpt("c config", filepath.Clean(dir+"/gaoler.json"), "Path to the configuration file")

	Gaoler.Spec = "[--config=<config-file>] [--root=<root-package>]"

	Gaoler.Before = middleware.Compute(
		ctx,
		initProject(rootPath),
	)

	Gaoler.Action = func() {
		Gaoler.PrintLongHelp()
	}
}
