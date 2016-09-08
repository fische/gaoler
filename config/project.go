package config

import (
	"encoding/json"
	"os"

	"github.com/fische/gaoler/project"
	"github.com/fische/gaoler/project/dependency"
	"github.com/fische/gaoler/project/dependency/pkg"
)

type Project struct {
	Name         string
	Dependencies map[string]*Dependency
}

const (
	openPerm = 0664
)

func NewProject(p *project.Project, deps dependency.Set) (*Project, error) {
	var (
		err error

		cfgDeps = make(map[string]*Dependency, len(deps))
	)
	for root, dep := range deps {
		if !p.IsDependency(dep) {
			cfgDeps[root], err = NewDependency(dep)
			if err != nil {
				return nil, err
			}
		}
	}
	name, err := pkg.GetPackagePath(p.Root)
	if err != nil {
		return nil, err
	}
	return &Project{
		Name:         name,
		Dependencies: cfgDeps,
	}, nil
}

func (p Project) Save(configPath string) error {
	file, err := os.OpenFile(configPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, openPerm)
	if err != nil {
		return err
	}
	e := json.NewEncoder(file)
	e.SetIndent("", "\t")
	return e.Encode(p)
}
