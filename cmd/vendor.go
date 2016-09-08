package cmd

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/fische/gaoler/project"
	"github.com/fische/gaoler/project/dependency"
	"github.com/jawher/mow.cli"
)

func cleanVendor(dir string) error {
	err := os.RemoveAll(dir)
	if err != nil {
		return err
	}
	return os.MkdirAll(dir, 0775)
}

func init() {
	Gaoler.Command("vendor", "Vendor dependencies of your project", func(cmd *cli.Cmd) {
		cmd.Spec = "[-t] [-s [-c]] [ROOT]"

		wd, err := os.Getwd()
		if err != nil {
			log.Errorf("Cannot get working directory : %v", err)
			cli.Exit(ExitFailure)
		}
		root := cmd.StringArg("ROOT", project.GetProjectRootFromDir(wd), "Root directory from a project")
		keepTests := cmd.BoolOpt("t test", false, "Keep test files")
		configPath := cmd.StringOpt("c config", "gaoler.json", "Path to the configuration file")
		save := cmd.BoolOpt("s save", false, "Save vendored dependencies to CONFIG file")

		cmd.Action = func() {
			p, err := project.New(*root, true)
			if err != nil {
				log.Errorf("Could not get dependencies : %v", err)
				cli.Exit(ExitFailure)
			} else if err = cleanVendor(p.Vendor); err != nil {
				log.Errorf("Could not clean vendor directory : %v", err)
				cli.Exit(ExitFailure)
			}
			var opts []dependency.CleanCheck
			if *keepTests {
				opts = append(opts, dependency.KeepTestFiles)
			} else {
				opts = append(opts, dependency.RemoveTestFiles)
			}
			for _, dep := range p.Dependencies {
				if !p.IsDependency(dep) {
					log.Printf("Cloning of %s...", dep.RootPackage)
					err = dep.Vendor(p.Vendor, opts...)
					if err != nil {
						log.Errorf("Could not clone repository of package %s : %v", dep.RootPackage, err)
						cli.Exit(ExitFailure)
					}
					log.Printf("Successful clone of %s", dep.RootPackage)
				}
			}
			if *save {
				err = p.Save(*configPath)
				if err != nil {
					log.Errorf("Could not save config : %v", err)
					cli.Exit(ExitFailure)
				}
			}
		}
	})
}
