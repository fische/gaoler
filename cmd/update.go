package cmd

// import (
// 	log "github.com/Sirupsen/logrus"
// 	"github.com/fische/gaoler/config"
// 	"github.com/fische/gaoler/project"
// 	cli "github.com/jawher/mow.cli"
// )
//
// func init() {
// 	Gaoler.Command("update", "Update dependencies of your project", func(cmd *cli.Cmd) {
// 		cmd.Spec = "[-s] [-t]"
//
// 		save := cmd.BoolOpt("s save", false, "Save vendored dependencies to CONFIG file")
// 		// test := cmd.BoolOpt("t test", false, "Include tests")
//
// 		cmd.Before = func() {
// 			if err := config.Setup(*configPath, false); err != nil {
// 				if err != nil {
// 					log.Errorf("Could not setup config : %v", err)
// 					cli.Exit(ExitFailure)
// 				}
// 			}
// 		}
//
// 		cmd.Action = func() {
// 			p, err := project.New(*root)
// 			if err != nil {
// 				log.Errorf("Could not create a new project : %v", err)
// 				cli.Exit(ExitFailure)
// 			} else if err = config.Load(p); err != nil {
// 				log.Errorf("Could not load config : %v", err)
// 				cli.Exit(ExitFailure)
// 			}
// 			// var opts []dependency.CleanCheck
// 			// if *keepTests {
// 			// 	opts = append(opts, dependency.KeepTestFiles)
// 			// } else {
// 			// 	opts = append(opts, dependency.RemoveTestFiles)
// 			// }
// 			// for _, dep := range p.Dependencies {
// 			// 	if !p.IsDependency(dep) && (*force || !dep.IsVendored()) {
// 			// 		log.Printf("Cloning of %s...", dep.RootPackage)
// 			// 		err = dep.Vendor(p.Vendor, opts...)
// 			// 		if err != nil {
// 			// 			log.Errorf("Could not clone repository of package %s : %v", dep.RootPackage, err)
// 			// 			cli.Exit(ExitFailure)
// 			// 		}
// 			// 		log.Printf("Successful clone of %s", dep.RootPackage)
// 			// 	}
// 			// }
// 			if *save {
// 				err = config.Save(p)
// 				if err != nil {
// 					log.Errorf("Could not save config : %v", err)
// 					cli.Exit(ExitFailure)
// 				}
// 			}
// 		}
// 	})
// }
