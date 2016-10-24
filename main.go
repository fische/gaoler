package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/fische/gaoler/internal/cmd"
)

func main() {
	if err := cmd.Gaoler.Run(os.Args); err != nil {
		log.WithError(err).Fatal("Error while running CLI.")
	}
}
