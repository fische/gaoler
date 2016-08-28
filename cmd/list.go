package cmd

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/fische/gaoler/project"
	"github.com/jawher/mow.cli"
)

func init() {
	Gaoler.Command("list", "List imports from your project", func(cmd *cli.Cmd) {
		cmd.Spec = "[ROOT]"

		wd, err := os.Getwd()
		if err != nil {
			log.Errorf("Cannot get working directory : %v", err)
			cli.Exit(ExitFailure)
		}
		root := cmd.StringArg("ROOT", project.GetProjectRootFromDir(wd), "Root directory from a project")

		cmd.Action = func() {
			p := project.New(*root)
			deps, err := p.ListDependencies()
			if err != nil {
				log.Errorf("Could not get dependencies : %v", err)
				cli.Exit(ExitFailure)
			}
			for _, dep := range deps {
				log.Printf("%+v", dep.Name())
			}
		}
	})
}
