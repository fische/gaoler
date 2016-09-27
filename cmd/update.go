package cmd

import (
	"os"

	"github.com/fische/gaoler/config"
	"github.com/fische/gaoler/project"
	"github.com/fische/gaoler/project/dependency"
	cli "github.com/jawher/mow.cli"
	"github.com/lunny/log"
)

func init() {
	Gaoler.Command("update", "Update dependencies of your project", func(cmd *cli.Cmd) {
		var (
			cfg *config.Config
			p   *project.Project

			save = cmd.BoolOpt("s save", false, "Save updated dependencies to CONFIG file")
			test = cmd.BoolOpt("t test", false, "Include tests")
		)

		cmd.Spec = "[-s] [-t]"

		cmd.Before = func() {
			flags := config.Load
			if *save {
				flags |= config.Save
			}

			var err error
			if p, err = project.New(*mainPath); err != nil {
				log.Errorf("Could not get project : %v", err)
				cli.Exit(ExitFailure)
			} else if cfg, err = config.New(p, *configPath, flags); err != nil {
				log.Errorf("Could not get config : %v", err)
				cli.Exit(ExitFailure)
			}
		}

		cmd.Action = func() {
			p.Dependencies.OnDecoded = importDependency(*mainPath, false, true)

			if err := cfg.Load(); err != nil {
				log.Errorf("Could not load config : %v", err)
				cli.Exit(ExitFailure)
			}

			var (
				opts    []func(info os.FileInfo) dependency.CleanOption
				updated bool
			)
			if !*test {
				opts = append(opts, dependency.RemoveTestFiles)
			} else {
				opts = append(opts, dependency.KeepTestFiles)
			}
			for _, dep := range p.Dependencies.Deps() {
				if dep.IsUpdatable() {
					log.Printf("Updating %s...", dep.RootPackage())
					if u, err := dep.Update(p.Vendor()); err != nil {
						log.Errorf("Could not load config : %v", err)
						cli.Exit(ExitFailure)
					} else if err = dep.CleanVendor(p.Vendor(), opts...); err != nil {
						log.Errorf("Could not clean repository of package %s : %v", dep.RootPackage(), err)
						cli.Exit(ExitFailure)
					} else {
						updated = updated || u
					}
					log.Printf("Successful update of %s", dep.RootPackage())
				}
			}
			if updated && *save {
				if err := cfg.Save(); err != nil {
					log.Errorf("Could not save config : %v", err)
					cli.Exit(ExitFailure)
				}
			}
		}
	})
}
