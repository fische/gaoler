package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/fische/gaoler/cmd"
)

//IDEA Add Dev dependencies

func main() {
	if err := cmd.Gaoler.Run(os.Args); err != nil {
		log.WithError(err).Fatal("Error while running CLI.")
	}
}
