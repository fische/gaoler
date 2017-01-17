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
	Gaoler.Command("diff", "Establish a diff of the current dependency tree to the config file", func(cmd *cli.Cmd) {
		list := cmd.BoolOpt("l list", false, "List packages diff for each dependency")
		dependencies := cmd.StringsArg("DEPENDENCIES", []string{}, "Regular expressions for filtering dependencies")

		cmd.Spec = "[-l] [DEPENDENCIES...]"

		cmd.Before = middleware.Compute(
			ctx,
			initConfig(configPath, func() config.Flags { return config.Load }),
			initRegexps(depRegexpsKey, dependencies),
		)

		cmd.Action = func() {
			regexps := ctx.Value(depRegexpsKey).([]*regexp.Regexp)
			p := ctx.Value("project").(*project.Project)
			if cfg := ctx.Value("config").(*config.Config); cfg == nil {
				fmt.Fprintln(os.Stderr, "Cannot diff depedencies without a config file")
				cli.Exit(ExitFailure)
			} else if err := cfg.Load(p); err != nil {
				fmt.Fprintf(os.Stderr, "Could not load project : %v\n", err)
				cli.Exit(ExitFailure)
			}

			s := pkg.NewSet()
			s.Constructor = &pkg.BasicFactory{
				SrcPath:              p.RootPath(),
				IgnoreVendor:         false,
				ImportCanFail:        false,
				IncludePseudoPackage: false,
			}
			s.Filter = filter.New(false, p.RootPath())
			if err := s.ListFrom(p.RootPath()); err != nil {
				fmt.Fprintf(os.Stderr, "Could not list packages : %v\n", err)
				cli.Exit(ExitFailure)
			}
			diff := p.Diff(s)
			remaining := dependency.NewSet()
			remaining.OnPackageAdded = func(p *pkg.Package, dep *dependency.Dependency) error {
				if !p.IsVendored() && p.Dir() != "" && !dep.HasOpenedRepository() {
					if err := dep.OpenRepository(p.Dir()); err == nil {
						return dep.SetRootPackage()
					}
				}
				return nil
			}
			if err := remaining.MergePackageSet(s); err != nil {
				fmt.Fprintf(os.Stderr, "Could not merge package set : %v\n", err)
				cli.Exit(ExitFailure)
			}

			var (
				// Colors
				untouchedDependency = color.New(color.Bold, color.FgWhite).SprintfFunc()
				removedDependency   = color.New(color.Bold, color.FgRed).SprintfFunc()
				addedDependency     = color.New(color.Bold, color.FgGreen).SprintfFunc()
				dependencyColor     = color.New(color.Bold, color.FgWhite).SprintfFunc()
				dependencyInfo      = color.New(color.Bold, color.FgBlue).SprintfFunc()
				packageInfo         = color.New(color.FgHiBlack).SprintfFunc()
				packageColor        = color.New(color.FgWhite).SprintfFunc()

				// Newline flag
				newline = false

				// Counters
				added     int
				removed   int
				untouched int

				// Output Buffer
				out string
			)
			for rootPackage, dep := range diff {
				if len(regexps) == 0 || checkRegexps(rootPackage, regexps) {
					if *list && newline {
						out += "\n"
					} else {
						newline = true
					}
					if dep.IsRemoved() {
						out += removedDependency("➖ ")
						removed++
					} else {
						out += untouchedDependency("• ")
						untouched++
					}
					nbAdded := len(dep.Added().Packages())
					nbRemoved := len(dep.Removed().Packages())
					nbUntouched := len(dep.Untouched().Packages())
					out += dependencyColor("%s ", rootPackage) + dependencyInfo("(%d Packages, %d Added , %d Removed)\n", nbAdded+nbRemoved+nbUntouched, nbAdded, nbRemoved)
					if *list {
						for _, p := range dep.Untouched().Packages() {
							out += packageColor("\t%s\n", p.Path())
						}
						for _, p := range dep.Removed().Packages() {
							out += packageColor("\t%s ", p.Path()) + packageInfo("Removed\n")
						}
						for _, p := range dep.Added().Packages() {
							out += packageColor("\t%s ", p.Path()) + packageInfo("Added\n")
						}
					}
				}
			}
			for rootPackage, dep := range remaining.Dependencies() {
				if len(regexps) == 0 || checkRegexps(rootPackage, regexps) {
					if *list && newline {
						out += "\n"
					} else {
						newline = true
					}
					nbAdded := len(dep.Packages())
					out += addedDependency("➕ ") + dependencyColor("%s ", rootPackage) + dependencyInfo("(%d Packages, %d Added , 0 Removed)\n", nbAdded, nbAdded)
					if *list {
						for _, p := range dep.Packages() {
							out += packageColor("\t%s ", p.Path()) + packageInfo("Added\n")
						}
					}
					added++
				}
			}
			if newline {
				out += "\n"
			}
			fmt.Println(out+"Dependencies: ", addedDependency("%d Added", added), removedDependency("%d Removed", removed), untouchedDependency("%d Untouched", untouched))
		}

		cmd.After = middleware.Compute(
			ctx,
			closeConfig,
		)
	})
}
