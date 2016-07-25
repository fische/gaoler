package cmd

import (
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/fische/gaoler/pkg"
	"github.com/fische/gaoler/project"
	"github.com/jawher/mow.cli"
)

func isValidFile(path string) bool {
	return filepath.Ext(path) == ".go"
}

func isValidDir(path string) bool {
	name := filepath.Base(path)
	return name != ".git" && name != "vendor"
}

func walkError(path string, err error) error {
	log.WithError(err).WithField("file", path).Errorf("Could not walk through file")
	return err
}

func init() {
	Gaoler.Command("list", "List imports from your project", func(cmd *cli.Cmd) {
		cmd.Action = func() {
			//TODO Take something else than PWD
			//TODO Check if vendor or testdata directory
			p := project.New(os.Getenv("PWD"))
			filepaths := p.Walk(isValidFile, isValidDir, walkError)
			m := make(map[string]*pkg.Import)
			for file := range filepaths {
				imports, err := pkg.GetImports(file)
				if err != nil {
					log.WithError(err).WithField("file", file).Error("Could not get imports from file")
				} else {
					for _, i := range imports {
						s, err := pkg.NewImport(i)
						if err != nil {
							log.WithError(err).WithField("import", i.Path.Value).Error("Could create new import")
						} else if _, ok := m[i.Path.Value]; !s.Goroot && !ok {
							m[i.Path.Value] = s
						}
					}
				}
			}
			log.Infof("%+v", m)
		}
	})
}
