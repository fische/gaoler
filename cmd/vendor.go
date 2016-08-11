package cmd

import (
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/fische/gaoler/project"
	"github.com/fische/gaoler/vcs"
	"github.com/fische/gaoler/vcs/modules"
	"github.com/jawher/mow.cli"
)

//TODO Add --save, --save-dev
//TODO Clean restore if something went wrong
//TODO Support bare repositories for git

func init() {
	Gaoler.Command("vendor", "Vendor all dependencies", func(cmd *cli.Cmd) {
		cmd.Action = func() {
			//TODO Take something else than PWD
			p := project.New(os.Getenv("PWD"))
			//TODO Operation on vendor dir through a unique absolute path
			if err := os.RemoveAll("vendor/"); err != nil {
				log.WithError(err).Fatal("Could not remove vendor directory")
			} else if err = os.Mkdir("vendor/", 0755); err != nil {
				log.WithError(err).Fatal("Could not create vendor directory")
			}
			deps, errch := p.GetDependencies()
			for {
				select {
				case dep := <-deps:
					if dep == nil {
						return
					}
					vcsName, repo, err := dep.GetRepository()
					if err != nil {
						log.WithError(err).WithField("dependency", dep.Name).Error("Could not get repository from dependency")
						return
					}
					rev, err := repo.GetRevision()
					if err != nil {
						log.WithError(err).WithField("dependency", dep.Name).Error("Could not get revision from repository")
						return
					} else if err = os.MkdirAll(filepath.Join(p.Root, "vendor/"+dep.Name), 0755); err != nil {
						log.WithError(err).WithField("dependency", dep.Name).Error("Could not create vendored package directory")
						return
					}
					remotes, err := repo.GetRemotes()
					if err != nil {
						log.WithError(err).WithField("dependency", dep.Name).Error("Could not get remotes from repository")
						return
					}

					v, ok := modules.GetVCS(vcsName)
					if !ok {
						log.WithError(err).WithField("dependency", dep.Name).WithField("module", vcsName).Error("Unknown VCS module")
						return
					}
					_, err = vcs.CloneAtRevision(v, filepath.Join(p.Root, "vendor/"+dep.Name), rev, remotes)
					if err != nil {
						log.WithError(err).WithField("dependency", dep.Name).Error("Could not clone")
						return
					}
					log.Infof("%s - %s: %s", dep.Name, vcsName, rev)
				case err := <-errch:
					if err == nil {
						return
					} else if e, ok := err.(*project.ErrorMessage); ok {
						log.WithError(e.Err).WithFields(log.Fields(e.Fields)).Error(e.Message)
					} else {
						log.WithError(err).Error("Could not get imports")
					}
					return
				}
			}
		}
	})
}
