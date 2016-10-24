package cmd

import (
	"github.com/fische/gaoler/internal/cmd/middleware"
	"github.com/fische/gaoler/project"
)

func setProject(root *string) middleware.Middleware {
	return func(ctx *middleware.Context) {
		ctx.Set("project", project.New(*root))
	}
}
