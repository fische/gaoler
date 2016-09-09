package cmd

import (
	"os"

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
	Gaoler.Spec = "[-v] [-c] [ROOT]"

	wd, err := os.Getwd()
	if err != nil {
		log.Errorf("Cannot get working directory : %v", err)
		cli.Exit(ExitFailure)
	}
	root = Gaoler.StringArg("ROOT", project.GetProjectRootFromDir(wd), "Root directory from a project")
	configPath = Gaoler.StringOpt("c config", "gaoler.json", "Path to the configuration file")
	verbose := Gaoler.BoolOpt("v verbose", false, "Enable verbose mode")

	Gaoler.Before = func() {
		if *verbose {
			log.SetLevel(log.DebugLevel)
		} else {
			log.SetLevel(log.InfoLevel)
		}
		pkg.SetSourcePath(*root)
	}

	Gaoler.Action = func() {
		Gaoler.PrintLongHelp()
	}
}
