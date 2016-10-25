package dependency

import (
	"encoding/json"

	"github.com/fische/gaoler/pkg"
)

type dependencyFormat struct {
	Packages *pkg.Set
	VCS      string
	Remote   string
	Revision string
	Branch   string `json:",omitempty" yaml:",omitpempty"`
}

type setFormat struct {
	Dependencies map[string]*Dependency
}

func (d Dependency) encode() interface{} {
	return dependencyFormat{
		Packages: d.Set,
		VCS:      d.vcs,
		Remote:   d.remote,
		Revision: d.revision,
		Branch:   d.branch,
	}
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
