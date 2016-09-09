package cmd

import (
	"io"

	log "github.com/Sirupsen/logrus"
	"github.com/fische/gaoler/config"
	"github.com/fische/gaoler/project"
	"github.com/jawher/mow.cli"
)

func init() {
	Gaoler.Command("list", "List dependencies of your project", func(cmd *cli.Cmd) {
		cmd.Spec = "[-t]"

		keepTests := cmd.BoolOpt("t test", false, "Keep test files")

		cmd.Before = func() {
			err := config.Setup(*configPath, false)
			if err != nil {
				log.Errorf("Could not setup config : %v", err)
				cli.Exit(ExitFailure)
			}
		}

		cmd.Action = func() {
			p, err := project.New(*root)
			if err != nil {
				log.Errorf("Could not create a new project : %v", err)
				cli.Exit(ExitFailure)
			} else if err = config.Load(p); err != nil && err != io.EOF {
				log.Errorf("Could not load config : %v", err)
				cli.Exit(ExitFailure)
			} else if err != nil {
				p, err = project.NewWithDependencies(*root, *keepTests, false)
				if err != nil {
					log.Errorf("Could not create a new project : %v", err)
					cli.Exit(ExitFailure)
				}
			}
			for _, dep := range p.Dependencies {
				for _, pkg := range dep.Packages {
					log.Printf("%s : %v", pkg.Path, pkg.IsVendored())
				}
			}
		}
	})
}
