package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/fische/gaoler/cmd/flags"
	"github.com/fische/gaoler/pkg/set"
	"github.com/jawher/mow.cli"
)

func init() {
	Gaoler.Command("list", "List dependencies of your project", func(cmd *cli.Cmd) {
		cmd.Spec = "[-t]"

		cmd.VarOpt("t test", flags.New(setNoTestPkg, false, true), "Include tests dependencies")

		cmd.Action = func() {
			s, err := set.ListPackagesFrom(*main, pkgFlags)
			if err != nil {
				log.Errorf("Could not list package : %v", err)
				cli.Exit(ExitFailure)
			}
			pkgs := s.Packages()
			for key := range pkgs {
				log.Printf("%v", key)
			}
		}
	})
}
