package cmd

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/fische/gaoler/project"
	"github.com/jawher/mow.cli"
)

func init() {
	Gaoler.Command("list", "List imports from your project", func(cmd *cli.Cmd) {
		cmd.Action = func() {
			//TODO Take something else than PWD
			//TODO Check if vendor or testdata directory
			p := project.New(os.Getenv("PWD"))
			deps, errch := p.GetDependencies()
			for {
				select {
				case dep := <-deps:
					if dep == nil {
						return
					}
					log.Infof("%+v", dep.Path.Value)
				case err := <-errch:
					if err == nil {
						return
					} else if e, ok := err.(*project.ErrorMessage); ok {
						log.WithError(e.Err).WithFields(log.Fields(e.Fields)).Error(e.Message)
					} else {
						log.WithError(err).Error("Could not get imports")
					}
					return
				}
			}
		}
	})
}
