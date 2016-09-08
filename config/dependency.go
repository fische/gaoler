package config

import "github.com/fische/gaoler/project/dependency"

type Dependency struct {
	Packages []string

	//VCS
	Revision string
	Remote   string
	VCS      string
	Branch   string `json:",omitempty"`
}

func NewDependency(dep *dependency.Dependency) (*Dependency, error) {
	ret := &Dependency{
		Packages: make([]string, len(dep.Packages)),
		VCS:      dep.Repository.GetVCSName(),
	}
	for idx, pkg := range dep.Packages {
		ret.Packages[idx] = pkg.Path
	}
	var err error
	if ret.Revision, err = dep.Repository.GetRevision(); err != nil {
		return nil, err
	} else if ret.Remote, err = dep.Repository.GetRemote(); err != nil {
		return nil, err
	}
	return ret, nil
}
