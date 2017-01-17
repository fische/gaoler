package cmd

import (
	"fmt"
	"os"
	"regexp"

	"github.com/fatih/color"
	"github.com/fische/gaoler/internal/cmd/middleware"
	"github.com/fische/gaoler/internal/config"
	"github.com/fische/gaoler/project"
	"github.com/fische/gaoler/project/dependency"
	cli "github.com/jawher/mow.cli"
)

func init() {
	Gaoler.Command("update", "Update vendored project's dependencies from config", func(cmd *cli.Cmd) {
		save := cmd.BoolOpt("p persistent", false, "Save dependencies state to config file")
		dependencies := cmd.StringsArg("DEPENDENCIES", []string{}, "Regular expressions for filtering dependencies to update")

		cmd.Spec = "[-p] [DEPENDENCIES...]"

		cmd.Before = middleware.Compute(
			ctx,
			initConfig(configPath, func() config.Flags {
				flags := config.Load
				if *save {
					flags |= config.Save
				}
				return flags
			}),
			initRegexps(depRegexpsKey, dependencies),
		)

		cmd.Action = func() {
			regexps := ctx.Value(depRegexpsKey).([]*regexp.Regexp)
			p := ctx.Value("project").(*project.Project)
			p.OnDecoded = func(dep *dependency.Dependency) error {
				for _, decoded := range dep.Packages() {
					decoded.Import(p.RootPath(), false)
				}
				return nil
			}
			var cfg *config.Config
			if cfg = ctx.Value("config").(*config.Config); cfg != nil {
				if err := cfg.Load(p); err != nil {
					fmt.Fprintf(os.Stderr, "Could not load project : %v\n", err)
					cli.Exit(ExitFailure)
				}
			}

			var (
				updated uint
				errors  uint

				// Colors
				dependencyColor = color.New(color.Bold, color.FgWhite)
				actionColor     = color.New(color.FgHiBlack)
				errorColor      = color.New(color.Bold, color.FgRed).SprintfFunc()
				resultColor     = color.New(color.Bold, color.FgWhite)
			)
			for rootPackage, dep := range p.Dependencies() {
				if dep.IsVendored() && (len(regexps) == 0 || checkRegexps(rootPackage, regexps)) {
					var (
						changed bool
						err     error
					)
					actionColor.Print("Updating:\t")
					dependencyColor.Println(rootPackage)
					if changed, err = dep.Update(p.VendorPath()); err != nil {
						fmt.Fprintf(os.Stderr, errorColor("Could not update dependency : %v\n", err))
						continue
					}
					actionColor.Print("Cleaning:\t")
					dependencyColor.Println(rootPackage)
					if err = dep.CleanVendor(p.VendorPath()); err != nil {
						fmt.Fprint(os.Stderr, errorColor("Could not clean vendored dependency : %v\n", err))
						errors++
						continue
					}
					if changed {
						updated++
					}
				}
			}
			resultColor.Printf("Updated: %d    Errors: %d\n", updated, errors)
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
