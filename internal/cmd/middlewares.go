package cmd

import (
	"fmt"
	"os"

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

func initConfig(configPath *string, flags config.Flags) middleware.Middleware {
	return func(ctx *middleware.Context) {
		if cfg, err := config.New(*configPath, flags); err != nil && !os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Could not init config : %v\n", err)
			cli.Exit(ExitFailure)
		} else {
			ctx.Set("config", cfg)
		}
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
