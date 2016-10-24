package dependency

import "encoding/json"

func (s *State) decode(f *stateFormat) {
	s.vcs = f.VCS
	s.revision = f.Revision
	s.remote = f.Remote
	s.branch = f.Branch
}

func (s *State) UnmarshalJSON(data []byte) error {
	f := &stateFormat{}
	if err := json.Unmarshal(data, f); err != nil {
		return err
	}
	s.decode(f)
	return nil
}

func (s *State) UnmarshalYAML(unmarshal func(interface{}) error) error {
	f := &stateFormat{}
	if err := unmarshal(f); err != nil {
		return err
	}
	s.decode(f)
	return nil
}

func (s *Set) decode(f *setFormat) error {
	s.dependencies = f.Dependencies
	if s.OnDecoded != nil {
		for rootPackage, dep := range s.dependencies {
			dep.rootPackage = rootPackage
			if err := s.OnDecoded(dep); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Set) UnmarshalJSON(data []byte) error {
	f := &setFormat{}
	if err := json.Unmarshal(data, f); err != nil {
		return err
	}
	return s.decode(f)
}

func (s *Set) UnmarshalYAML(unmarshal func(interface{}) error) error {
	f := &setFormat{}
	if err := unmarshal(f); err != nil {
		return err
	}
	return s.decode(f)
}
