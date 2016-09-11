package cmd

import (
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/fische/gaoler/project"
	"github.com/fische/gaoler/project/dependency/pkg"
	"github.com/jawher/mow.cli"
)

var (
	Gaoler = cli.App("goaler", "A Go package manager")

	ExitSuccess = 0
	ExitFailure = 1

	root       *string
	configPath *string
)

func init() {
	Gaoler.Spec = "[-v] [--config=<config-file>] [ROOT]"

	wd, err := os.Getwd()
	if err != nil {
		log.Errorf("Could not get working directory : %v", err)
		cli.Exit(ExitFailure)
	}
	dir, _ := project.GetProjectRootFromDir(wd)
	root = Gaoler.StringArg("ROOT", dir, "Root directory from a project")
	configPath = Gaoler.StringOpt("c config", filepath.Clean(dir+"/gaoler.json"), "Path to the configuration file")
	verbose := Gaoler.BoolOpt("v verbose", false, "Enable verbose mode")

	Gaoler.Before = func() {
		if *verbose {
			log.SetLevel(log.DebugLevel)
		} else {
			log.SetLevel(log.InfoLevel)
		}
		d, err := filepath.Abs(*root)
		if err != nil {
			log.Errorf("Could get absolute path for root directory : %v", err)
			cli.Exit(ExitFailure)
		}
		root = &d
		pkg.SetSourcePath(*root)
	}

	Gaoler.Action = func() {
		Gaoler.PrintLongHelp()
	}
}
