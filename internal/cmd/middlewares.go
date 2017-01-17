package cmd

import (
	"fmt"
	"os"
	"regexp"

	"github.com/fische/gaoler/internal/cmd/middleware"
	"github.com/fische/gaoler/internal/config"
	"github.com/fische/gaoler/project"
	cli "github.com/jawher/mow.cli"
)

func initProject(rootPath *string) middleware.Middleware {
	return func(ctx *middleware.Context) {
		ctx.Set("project", project.New(*rootPath))
	}
}

func initConfig(configPath *string, computeFlags func() config.Flags) middleware.Middleware {
	return func(ctx *middleware.Context) {
		if cfg, err := config.New(*configPath, computeFlags()); err != nil && !os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Could not init config : %v\n", err)
			cli.Exit(ExitFailure)
		} else {
			ctx.Set("config", cfg)
		}
	}
}

func initRegexps(key string, regexps *[]string) middleware.Middleware {
	return func(ctx *middleware.Context) {
		var (
			arr = make([]*regexp.Regexp, len(*regexps))
			err error
		)
		for idx, r := range *regexps {
			if arr[idx], err = regexp.Compile(r); err != nil {
				fmt.Fprintf(os.Stderr, "Could not compile regexp %s : %v\n", r, err)
				cli.Exit(ExitFailure)
			}
		}
		ctx.Set(key, arr)
	}
}

func closeConfig(ctx *middleware.Context) {
	if cfg := ctx.Value("config").(*config.Config); cfg != nil {
		if err := cfg.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Could not close config : %v\n", err)
			cli.Exit(ExitFailure)
		}
	}
}
