package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/jawher/mow.cli"
)

var (
	Gaoler = cli.App("goaler", "A Go package manager")

	ExitSuccess = 0
	ExitFailure = 1
)

func init() {
	Gaoler.Spec = "[-v]"

	verbose := Gaoler.Bool(cli.BoolOpt{Name: "v verbose", Value: false, Desc: "Enable verbose mode"})

	Gaoler.Before = func() {
		if *verbose {
			log.SetLevel(log.DebugLevel)
		} else {
			log.SetLevel(log.InfoLevel)
		}
	}

	Gaoler.Action = func() {
		Gaoler.PrintLongHelp()
	}
}
