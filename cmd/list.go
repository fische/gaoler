package cmd

import (
	"fmt"
	"os"

	"github.com/fische/gaoler/config"
	"github.com/fische/gaoler/pkg"
	"github.com/fische/gaoler/project"
	"github.com/jawher/mow.cli"
	"github.com/lunny/log"
	"github.com/ttacon/chalk"
)

func init() {
	Gaoler.Command("list", "List dependencies of your project", func(cmd *cli.Cmd) {
		var (
			test = cmd.BoolOpt("t test", false, "Include tests dependencies")

			cfg *config.Config
			p   *project.Project
		)

		cmd.Spec = "[-t]"

		cmd.Before = func() {
			var err error
			if p, err = project.New(*mainPath); err != nil {
				log.Errorf("Could not get project : %v", err)
				cli.Exit(ExitFailure)
			} else if cfg, err = config.New(p, *configPath); err != nil {
				log.Errorf("Could not get config : %v", err)
				cli.Exit(ExitFailure)
			}
		}

		cmd.Action = func() {
			p.Dependencies.Filter = filterUsefulDependencies
			p.Dependencies.OnDecoded = importDependency(*mainPath, false, true)

			var s *pkg.Set
			if err := cfg.Load(); err != nil && !os.IsNotExist(err) {
				log.Errorf("Could not load config : %v", err)
				cli.Exit(ExitFailure)
			} else if os.IsNotExist(err) {
				s = pkg.NewSet()
			} else {
				s = p.Dependencies.ToPackageSet()
			}

			s.OnAdded = importPackage(*mainPath, false)
			if !*test {
				s.Filter = pkg.NoTestFiles
			}

			if err := s.ListFrom(*mainPath); err != nil {
				log.Errorf("Could not list packages : %v", err)
				cli.Exit(ExitFailure)
			}

			for key, p := range s.Packages {
				if !(p.IsLocal() || p.IsPseudoPackage() || p.IsStandardPackage()) {
					if p.IsSaved() {
						if p.IsVendored() {
							fmt.Println(chalk.Green, key)
						} else {
							fmt.Println(chalk.Yellow, key)
						}
					} else {
						fmt.Println(chalk.Red, key)
					}
				}
			}
			fmt.Print(chalk.Reset)
		}
	})
}
