package cmd

import (
	"io"
	"os"

	"github.com/fische/gaoler/config"
	"github.com/fische/gaoler/pkg"
	"github.com/fische/gaoler/project"
	"github.com/fische/gaoler/project/dependency"
	cli "github.com/jawher/mow.cli"
	"github.com/lunny/log"
)

func init() {
	Gaoler.Command("vendor", "Vendor dependencies of your project", func(cmd *cli.Cmd) {
		var (
			cfg *config.Config
			p   *project.Project

			test  = cmd.BoolOpt("t test", false, "Include tests")
			save  = cmd.BoolOpt("s save", false, "Save vendored dependencies to CONFIG file")
			force = cmd.BoolOpt("f force", false, "Force the regeneration of the vendor directory")
		)

		cmd.Spec = "[-t] [-s] [-f]"

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

			if *force {
				if err := resetDir(p.Vendor()); err != nil {
					log.Errorf("Could not reset vendor directory : %v", err)
					cli.Exit(ExitFailure)
				}
			}
		}

		cmd.Action = func() {
			changed := false
			p.Dependencies.Filter = filterUsefulDependencies
			p.Dependencies.OnPackageAdded = func(p *pkg.Package, dep *dependency.Dependency) error {
				if !p.Vendored && dep.Repository == nil {
					if err := dep.OpenRepository(p.Dir); err == nil {
						if err = dep.LockCurrentState(); err != nil {
							return err
						}
					} else {
						log.Warnf("Could not open repository of %s : %v", p.Path, err)
					}
				}
				changed = true
				return nil
			}
			p.Dependencies.OnDecoded = importDependency(*mainPath, false, true)

			var s *pkg.Set
			if *force {
				s = pkg.NewSet()
			} else {
				if err := cfg.Load(); err != nil && err != io.EOF {
					log.Errorf("Could not load config : %v", err)
					cli.Exit(ExitFailure)
				} else if err == io.EOF {
					s = pkg.NewSet()
				} else {
					s = p.Dependencies.ToPackageSet()
				}
			}

			s.OnAdded = importPackage(*mainPath, false)
			if !*test {
				s.Filter = pkg.IsNotGoTestFile
			}

			if err := s.ListFrom(*mainPath); err != nil {
				log.Errorf("Could not list packages : %v", err)
				cli.Exit(ExitFailure)
			} else if err = p.Dependencies.MergePackageSet(s); err != nil {
				log.Errorf("Could not merge dependencies : %v", err)
				cli.Exit(ExitFailure)
			}

			var opts []func(info os.FileInfo) dependency.CleanOption
			if !*test {
				opts = append(opts, dependency.RemoveGoTestFiles)
			} else {
				opts = append(opts, dependency.KeepGoTestFiles)
			}
			for _, dep := range p.Dependencies.Deps {
				if dep.IsVendorable() && (*force || !dep.IsVendored()) {
					log.Printf("Cloning of %s...", dep.RootPackage)
					if err := dep.Vendor(p.Vendor()); err != nil {
						log.Errorf("Could not clone repository of package %s : %v", dep.RootPackage, err)
						cli.Exit(ExitFailure)
					} else if err = dep.CleanVendor(p.Vendor(), opts...); err != nil {
						log.Errorf("Could not clean repository of package %s : %v", dep.RootPackage, err)
						cli.Exit(ExitFailure)
					}
					log.Printf("Successful clone of %s", dep.RootPackage)
				}
			}
			if changed && *save {
				if err := cfg.Save(); err != nil {
					log.Errorf("Could not save config : %v", err)
					cli.Exit(ExitFailure)
				}
			}
		}
	})
}
