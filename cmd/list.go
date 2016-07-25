package cmd

import (
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/fische/gaoler/pkg"
	"github.com/jawher/mow.cli"
)

func init() {
	Gaoler.Command("list", "List imports from your project", func(cmd *cli.Cmd) {
		cmd.Spec = "[FILE...]"

		files := cmd.StringsArg("FILE", nil, "Files from which imports will be listed")

		cmd.Action = func() {
			if *files != nil {
				for _, file := range *files {
					if filepath.Ext(file) == ".go" {
						imports, err := pkg.ListImports(file)
						if err != nil {
							log.WithError(err).Errorf("Could not get imports from : %v", file)
						} else {
							for _, i := range imports {
								log.Infof("%+v", i.Path)
							}
						}
					} else {
						log.WithField("file", file).Warn("Cannot get import from files without `.go` extension")
					}
				}
			}
		}
	})
}
