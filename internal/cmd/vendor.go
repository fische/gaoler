package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/fische/gaoler/internal/cmd/middleware"
	"github.com/fische/gaoler/internal/config"
	"github.com/fische/gaoler/pkg"
	"github.com/fische/gaoler/pkg/filter"
	"github.com/fische/gaoler/project"
	"github.com/fische/gaoler/project/dependency"
	cli "github.com/jawher/mow.cli"
)

func init() {
	Gaoler.Command("vendor", "Vendor project's dependencies from config (unless it does not exist)", func(cmd *cli.Cmd) {
		save := cmd.BoolOpt("p persistent", false, "Save dependencies state to config file")
		sync := cmd.BoolOpt("s sync", false, "Synchronize with the actual dependency list")

		cmd.Spec = "[-p] [-s]"

		cmd.Before = middleware.Compute(
			ctx,
			initConfig(configPath, func() config.Flags {
				flags := config.Load
				if *save {
					flags |= config.Save
				}
				return flags
			}),
		)

		cmd.Action = func() {
			p := ctx.Value("project").(*project.Project)
			p.OnPackageAdded = func(p *pkg.Package, dep *dependency.Dependency) error {
				if !p.IsVendored() && p.Dir() != "" && !dep.HasOpenedRepository() {
					if err := dep.OpenRepository(p.Dir()); err != nil {
						return err
					} else if err = dep.SetRootPackage(); err != nil {
						return err
					} else if err = dep.LockCurrentState(); err != nil {
						return err
					}
				}
				return nil
			}
			p.OnDecoded = func(dep *dependency.Dependency) error {
				for _, decoded := range dep.Packages() {
					decoded.Import(p.RootPath(), false)
				}
				return nil
			}
			var cfg *config.Config
			if cfg = ctx.Value("config").(*config.Config); cfg != nil {
				if err := cfg.Load(p); err != nil {
					if err == io.EOF {
						*sync = true
					} else {
						fmt.Fprintf(os.Stderr, "Could not load project : %v\n", err)
						cli.Exit(ExitFailure)
					}
				}
			}

			var diff dependency.DiffSet
			if cfg == nil || *sync {
				s := pkg.NewSet()
				s.Constructor = &pkg.BasicFactory{
					SrcPath:              p.RootPath(),
					IgnoreVendor:         false,
					ImportCanFail:        true,
					IncludePseudoPackage: false,
				}
				s.Filter = filter.New(false, p.RootPath())
				if err := s.ListFrom(p.RootPath()); err != nil {
					fmt.Fprintf(os.Stderr, "Could not list packages : %v\n", err)
					cli.Exit(ExitFailure)
				}
				diff = p.Diff(s)
				if err := p.Apply(diff); err != nil {
					fmt.Fprintf(os.Stderr, "Could not apply diff to dependencies : %v\n", err)
					cli.Exit(ExitFailure)
				} else if err = p.MergePackageSet(s); err != nil {
					fmt.Fprintf(os.Stderr, "Could not merge package set : %v\n", err)
					cli.Exit(ExitFailure)
				}
			}

			var (
				vendored uint
				errors   uint
				removed  uint

				// Colors
				dependencyColor = color.New(color.Bold, color.FgWhite)
				actionColor     = color.New(color.FgHiBlack)
				errorColor      = color.New(color.Bold, color.FgRed).SprintfFunc()
				resultColor     = color.New(color.Bold, color.FgWhite)
			)
			for rootPackage, dep := range p.Dependencies() {
				if !dep.IsVendored() {
					actionColor.Print("Vendoring:\t")
					dependencyColor.Println(rootPackage)
					if err := dep.Vendor(p.VendorPath()); err != nil {
						fmt.Fprintf(os.Stderr, errorColor("Could not vendor dependency : %v\n", err))
						continue
					}
					actionColor.Print("Cleaning:\t")
					dependencyColor.Println(rootPackage)
					if err := dep.CleanVendor(p.VendorPath()); err != nil {
						fmt.Fprint(os.Stderr, errorColor("Could not clean vendored dependency : %v\n", err))
						errors++
						continue
					}
					vendored++
				}
			}
			for rootPackage, dep := range diff {
				if dep.IsRemoved() && dep.Removed().HasVendoredPackage() {
					actionColor.Print("Removing:\t")
					dependencyColor.Println(rootPackage)
					if err := removeEmptyParents(filepath.Clean(p.VendorPath()+rootPackage), p.VendorPath()); err != nil {
						fmt.Fprint(os.Stderr, errorColor("Could not remove dependency : %v\n", err))
						errors++
						continue
					}
					removed++
				}
			}
			resultColor.Printf("Vendored: %d    Removed: %d    Errors: %d\n", vendored, removed, errors)
			if *save && cfg != nil {
				if err := cfg.Save(p); err != nil {
					fmt.Fprint(os.Stderr, errorColor("Could not save config : %v\n", err))
					cli.Exit(ExitFailure)
				} else {
					resultColor.Println("New dependency set saved!")
				}
			}
		}

		cmd.After = middleware.Compute(
			ctx,
			closeConfig,
		)
	})
}
