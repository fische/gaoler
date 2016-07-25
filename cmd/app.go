package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/jawher/mow.cli"
)

//Gaoler is the CLI of Gaoler
var Gaoler = cli.App("goaler", "A Go package manager")

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
