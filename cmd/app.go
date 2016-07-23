package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/jawher/mow.cli"
)

//App is the CLI of Goaler
var App = cli.App("goaler", "A Go package manager")

func init() {
	App.Spec = "[-v]"

	verbose := App.Bool(cli.BoolOpt{Name: "v verbose", Value: false, Desc: "Enable verbose mode"})

	App.Before = func() {
		if *verbose {
			log.SetLevel(log.DebugLevel)
		} else {
			log.SetLevel(log.InfoLevel)
		}
	}

	App.Action = func() {
		App.PrintLongHelp()
	}
}
