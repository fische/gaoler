package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/fische/gaoler/cmd"
)

func main() {
	if err := cmd.App.Run(os.Args); err != nil {
		log.WithError(err).Fatal("Error while running CLI.")
	}
}
