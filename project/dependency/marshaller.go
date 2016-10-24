package dependency

import "encoding/json"

type stateFormat struct {
	VCS      string
	Remote   string
	Revision string
	Branch   string `json:",omitempty" yaml:",omitpempty"`
}

type setFormat struct {
	Dependencies map[string]*Dependency
}

func (s State) encode() interface{} {
	return stateFormat{
		VCS:      s.vcs,
		Remote:   s.remote,
		Revision: s.revision,
		Branch:   s.branch,
	}
}

func (s State) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.encode())
}

func (s State) MarshalYAML() (interface{}, error) {
	return s.encode(), nil
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
