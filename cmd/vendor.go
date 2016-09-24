package cmd

import (
	"io"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/fische/gaoler/cmd/flags"
	"github.com/fische/gaoler/config"
	"github.com/fische/gaoler/pkg"
	"github.com/fische/gaoler/pkg/set"
	"github.com/fische/gaoler/project"
	"github.com/fische/gaoler/project/dependency"
	"github.com/jawher/mow.cli"
)

func forceReset(f *config.Flags) {
	if err := resetDir(filepath.Clean(*main + "/vendor/")); err != nil {
		log.Errorf("Could not clean vendor directory : %v", err)
		cli.Exit(ExitFailure)
	}
}

func init() {
	Gaoler.Command("vendor", "Vendor dependencies of your project", func(cmd *cli.Cmd) {
		cmd.Spec = "[-t] [-s] [-f]"

		configFlags := config.LoadFlag
		cmd.VarOpt("t test", flags.New(setNoTestPkg, false, true), "Include tests")
		cmd.VarOpt("s save", flags.New(setConfigFlag(&configFlags, config.SaveFlag), false, true), "Save vendored dependencies to CONFIG file")
		force := cmd.BoolOpt("f force", false, "Force the regeneration of the vendor directory")

		cmd.Before = func() {
			if *force {
				forceReset(&configFlags)
			}
			if err := config.Setup(*configPath, configFlags); err != nil {
				log.Errorf("Could not setup config : %v", err)
				cli.Exit(ExitFailure)
			}
		}

		cmd.Action = func() {
			p, err := project.New(*main)
			if err != nil {
				log.Errorf("Could not create a new project : %v", err)
				cli.Exit(ExitFailure)
			} else if err = config.Load(p); err != nil && err != io.EOF {
				log.Errorf("Could not load config : %v", err)
				cli.Exit(ExitFailure)
			} else if err == io.EOF {
				var s *set.Packages
				s, err = set.ListPackagesFrom(*main, pkgFlags)
				if err != nil {
					log.Errorf("Could not list package : %v", err)
					cli.Exit(ExitFailure)
				}
				p.Dependencies.Set(s, true)
			} else {
				if err = p.Dependencies.Import(*main, pkgFlags); err != nil {
					log.Errorf("Could not import saved dependencies : %v", err)
					cli.Exit(ExitFailure)
				}
				s := p.Dependencies.GetPackagesSet()
				if err = s.CompleteFrom(*main, pkgFlags); err != nil {
					log.Errorf("Could not complete packages list : %v", err)
					cli.Exit(ExitFailure)
				}
				p.Dependencies.Merge(s, true)
			}
			var opts []dependency.CleanCheck
			if pkgFlags.Has(pkg.NoTests) {
				opts = append(opts, dependency.RemoveTestFiles)
			} else {
				opts = append(opts, dependency.KeepTestFiles)
			}
			for _, dep := range p.Dependencies {
				if *force || !dep.IsVendored() {
					log.Printf("Cloning of %s...", dep.RootPackage())
					err = dep.Vendor(p.Vendor())
					if err != nil {
						log.Errorf("Could not clone repository of package %s : %v", dep.RootPackage(), err)
						cli.Exit(ExitFailure)
					}
					err = dep.CleanVendor(p.Vendor(), opts...)
					if err != nil {
						log.Errorf("Could not clean repository of package %s : %v", dep.RootPackage(), err)
						cli.Exit(ExitFailure)
					}
					log.Printf("Successful clone of %s", dep.RootPackage())
				}
			}
			if configFlags.Has(config.SaveFlag) {
				err = config.Save(p)
				if err != nil {
					log.Errorf("Could not save config : %v", err)
					cli.Exit(ExitFailure)
				}
			}
		}
	})
}
