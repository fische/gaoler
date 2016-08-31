package cmd

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/fische/gaoler/project"
	"github.com/fische/gaoler/project/dependency"
	"github.com/jawher/mow.cli"
)

func init() {
	Gaoler.Command("vendor", "Vendor dependencies of your project", func(cmd *cli.Cmd) {
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
			err = os.RemoveAll(p.Vendor)
			if err != nil {
				log.Errorf("Could not clean vendor directory : %v", err)
				cli.Exit(ExitFailure)
			}
			err = os.MkdirAll(p.Vendor, 0775)
			if err != nil {
				log.Errorf("Could not create vendor directory : %v", err)
				cli.Exit(ExitFailure)
			}
			for _, dep := range deps {
				if !p.HasLocalDependency(dep) {
					log.Printf("Cloning of %s...", dep.RootPackage)
					err = dep.Vendor(p.Vendor, dependency.KeepOnlyGoFiles, dependency.RemoveTestFiles)
					if err != nil {
						log.Errorf("Could not clone repository of package %s : %v", dep.RootPackage, err)
						cli.Exit(ExitFailure)
					}
					log.Printf("Successful clone of %s", dep.RootPackage)
				}
			}
		}
	})
}
