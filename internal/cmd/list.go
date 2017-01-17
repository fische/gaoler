package cmd

import (
	"fmt"
	"os"
	"regexp"

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
	Gaoler.Command("list", "List project's dependencies from config (unless it does not exist)", func(cmd *cli.Cmd) {
		list := cmd.BoolOpt("l list", false, "List packages' status for each dependency")
		sync := cmd.BoolOpt("s sync", false, "Synchronize with the actual dependency list")
		dependencies := cmd.StringsArg("DEPENDENCIES", []string{}, "Regular expressions for filtering dependencies")

		cmd.Spec = "[-l] [-s] [DEPENDENCIES...]"

		cmd.Before = middleware.Compute(
			ctx,
			initConfig(configPath, func() config.Flags { return config.Load }),
			initRegexps(depRegexpsKey, dependencies),
		)

		cmd.Action = func() {
			regexps := ctx.Value(depRegexpsKey).([]*regexp.Regexp)
			p := ctx.Value("project").(*project.Project)
			p.OnPackageAdded = func(p *pkg.Package, dep *dependency.Dependency) error {
				if !p.IsVendored() && p.Dir() != "" && !dep.HasOpenedRepository() {
					if err := dep.OpenRepository(p.Dir()); err == nil {
						return dep.SetRootPackage()
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
			if cfg := ctx.Value("config").(*config.Config); cfg != nil {
				if err := cfg.Load(p); err != nil {
					fmt.Fprintf(os.Stderr, "Could not load project : %v\n", err)
					cli.Exit(ExitFailure)
				}
			} else {
				*sync = true
			}

			if *sync {
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
				} else if err = p.Apply(p.Diff(s)); err != nil {
					fmt.Fprintf(os.Stderr, "Could not apply diff to dependencies : %v\n", err)
					cli.Exit(ExitFailure)
				} else if err = p.MergePackageSet(s); err != nil {
					fmt.Fprintf(os.Stderr, "Could not merge package set : %v\n", err)
					cli.Exit(ExitFailure)
				}
			}

			var (
				// Colors
				savedDependency    = color.New(color.Bold, color.FgGreen).SprintfFunc()
				notsavedDependency = color.New(color.Bold, color.FgRed).SprintfFunc()
				dependencyColor    = color.New(color.Bold, color.FgWhite).SprintfFunc()
				dependencyInfo     = color.New(color.Bold, color.FgBlue).SprintfFunc()
				validColor         = color.New(color.FgGreen).SprintfFunc()
				invalidColor       = color.New(color.FgRed).SprintfFunc()
				packageColor       = color.New(color.FgWhite).SprintfFunc()
				packageInfo        = color.New(color.FgHiBlack).SprintfFunc()

				// Newline flag
				newline = false

				// Counters
				saved    int
				notsaved int

				// Output Buffer
				out string
			)
			for rootPackage, dep := range p.Dependencies() {
				if len(regexps) == 0 || checkRegexps(rootPackage, regexps) {
					if *list && newline {
						out += "\n"
					} else {
						newline = true
					}
					var (
						pkgSaved    int
						pkgVendored int
						pkgOut      string
					)
					for packagePath, p := range dep.Packages() {
						pkgOut += packageColor("\t%s", packagePath) + packageInfo(" Saved(")
						if p.IsSaved() {
							pkgSaved++
							pkgOut += validColor("✓")
						} else {
							pkgOut += invalidColor("×")
						}
						pkgOut += packageInfo(") Vendored(")
						if p.IsVendored() {
							pkgVendored++
							pkgOut += validColor("✓")
						} else {
							pkgOut += invalidColor("×")
						}
						pkgOut += packageInfo(")\n")
					}
					if *sync {
						if dep.IsSaved() {
							out += savedDependency("✔ ")
							saved++
						} else {
							out += notsavedDependency("❌ ")
							notsaved++
						}
					}
					out += dependencyColor("%s ", rootPackage) + dependencyInfo("(%d Packages, %d Saved, %d Vendored)\n", len(dep.Packages()), pkgSaved, pkgVendored)
					if *list {
						out += pkgOut
					}
				}
			}
			fmt.Print(out)
			if *sync {
				var nl string
				if newline {
					nl = "\n"
				}
				fmt.Println(nl+"Dependencies: ", savedDependency("%d Saved", saved), notsavedDependency("%d Unsaved", notsaved))
			}
		}

		cmd.After = middleware.Compute(
			ctx,
			closeConfig,
		)
	})
}
