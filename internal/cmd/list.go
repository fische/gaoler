package cmd

import (
	"fmt"
	"os"

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
		update := cmd.BoolOpt("u update", false, "Update project's dependencies compared to the current dependency tree")

		cmd.Spec = "[-l] [-u]"

		cmd.Before = middleware.Compute(
			ctx,
			initConfig(configPath, config.Load),
		)

		cmd.Action = func() {
			p := ctx.Value("project").(*project.Project)
			p.OnPackageAdded = func(p *pkg.Package, dep *dependency.Dependency) error {
				if p.Dir() != "" && !dep.HasOpenedRepository() {
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
					fmt.Fprintf(os.Stderr, "Could not load project : %v", err)
					cli.Exit(ExitFailure)
				}
			} else {
				*update = true
			}

			if *update {
				s := pkg.NewSet()
				s.Constructor = &pkg.BasicFactory{
					SrcPath:              p.RootPath(),
					IgnoreVendor:         false,
					ImportCanFail:        true,
					IncludePseudoPackage: false,
				}
				s.Filter = filter.New(false, p.RootPath())
				if err := s.ListFrom(p.RootPath()); err != nil {
					fmt.Fprintf(os.Stderr, "Could not list packages : %v", err)
					cli.Exit(ExitFailure)
				}
				diff := p.Diff(s)
				if err := p.Apply(diff); err != nil {
					fmt.Fprintf(os.Stderr, "Could not apply diff to dependencies : %v", err)
					cli.Exit(ExitFailure)
				} else if err := p.MergePackageSet(s); err != nil {
					fmt.Fprintf(os.Stderr, "Could not merge package set : %v", err)
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
				if *update {
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
			fmt.Print(out)
			if *update {
				var nl string
				if newline {
					nl = "\n"
				}
				fmt.Println(nl+"Dependencies: ", savedDependency("%d Saved", saved), notsavedDependency("%d Not Saved", notsaved))
			}
		}

		cmd.After = middleware.Compute(
			ctx,
			closeConfig,
		)
	})
}
