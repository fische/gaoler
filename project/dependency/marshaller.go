package dependency

import (
	"encoding/json"

	"github.com/fische/gaoler/pkg"
)

type dependencyFormat struct {
	Packages *pkg.Set
	VCS      string `json:",omitempty" yaml:",omitpempty"`
	Remote   string `json:",omitempty" yaml:",omitpempty"`
	Revision string `json:",omitempty" yaml:",omitpempty"`
	Branch   string `json:",omitempty" yaml:",omitpempty"`
}

type setFormat struct {
	Dependencies map[string]*Dependency
}

func (d Dependency) encode() interface{} {
	ret := dependencyFormat{
		Packages: d.Set,
	}
	if d.State != nil {
		ret.VCS = d.vcs
		ret.Remote = d.remote
		ret.Revision = d.revision
		ret.Branch = d.branch
	}
	return ret
}

func (d Dependency) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.encode())
}

func (d Dependency) MarshalYAML() (interface{}, error) {
	return d.encode(), nil
}

func (s Set) encode() interface{} {
	return setFormat{
		Dependencies: s.dependencies,
	}
}

func (s Set) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.encode())
}

func (s Set) MarshalYAML() (interface{}, error) {
	return s.encode(), nil
}
