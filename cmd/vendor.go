package cmd

import (
	"io"
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/fische/gaoler/config"
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
		cmd.Spec = "[-t] [-s] [-f]"

		keepTests := cmd.BoolOpt("t test", false, "Keep test files")
		save := cmd.BoolOpt("s save", false, "Save vendored dependencies to CONFIG file")
		force := cmd.BoolOpt("f force", false, "Force the regeneration of the vendor directory")

		cmd.Before = func() {
			if *force {
				if err := cleanVendor(filepath.Clean(*root + "/vendor/")); err != nil {
					log.Errorf("Could not clean vendor directory : %v", err)
					cli.Exit(ExitFailure)
				}
			}
			if err := config.Setup(*configPath, *force); err != nil {
				if err != nil {
					log.Errorf("Could not setup config : %v", err)
					cli.Exit(ExitFailure)
				}
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
			var opts []dependency.CleanCheck
			if *keepTests {
				opts = append(opts, dependency.KeepTestFiles)
			} else {
				opts = append(opts, dependency.RemoveTestFiles)
			}
			for _, dep := range p.Dependencies {
				if !p.IsDependency(dep) && (*force || !dep.IsVendored()) {
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
				err = config.Save(p)
				if err != nil {
					log.Errorf("Could not save config : %v", err)
					cli.Exit(ExitFailure)
				}
			}
		}
	})
}
