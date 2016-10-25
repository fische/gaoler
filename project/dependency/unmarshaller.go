package dependency

import (
	"encoding/json"

	"github.com/fische/gaoler/pkg"
)

func (d *Dependency) decode(f *dependencyFormat) {
	d.Set = f.Packages
	d.State = &State{
		vcs:      f.VCS,
		remote:   f.Remote,
		revision: f.Revision,
		branch:   f.Branch,
	}
}

func (d *Dependency) UnmarshalJSON(data []byte) error {
	f := &dependencyFormat{
		Packages: pkg.NewSet(),
	}
	if err := json.Unmarshal(data, f); err != nil {
		return err
	}
	d.decode(f)
	return nil
}

func (d *Dependency) UnmarshalYAML(unmarshal func(interface{}) error) error {
	f := &dependencyFormat{
		Packages: pkg.NewSet(),
	}
	if err := unmarshal(f); err != nil {
		return err
	}
	d.decode(f)
	return nil
}

func (s *Set) decode(f *setFormat) error {
	s.dependencies = f.Dependencies
	for rootPackage, dep := range s.dependencies {
		dep.rootPackage = rootPackage
	}
	if s.OnDecoded != nil {
		for _, dep := range s.dependencies {
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
