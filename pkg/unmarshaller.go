package pkg

import (
	"encoding/json"
	"errors"
)

func (p *Package) UnmarshalJSON(data []byte) error {
	if data[0] == '"' && data[len(data)-1] == '"' {
		p.path = string(data[1 : len(data)-1])
		p.saved = true
		return nil
	}
	return errors.New("Could not unmarshal package")
}

func (p *Package) UnmarshalYAML(unmarshal func(interface{}) error) error {
	if err := unmarshal(&p.path); err != nil {
		return err
	}
	p.saved = true
	return nil
}

func (s *Set) decode(packages []*Package) {
	for _, p := range packages {
		s.packages[p.Path()] = p
	}
}

func (s *Set) UnmarshalJSON(data []byte) error {
	var decode []*Package
	if err := json.Unmarshal(data, &decode); err != nil {
		return err
	}
	s.decode(decode)
	return nil
}

func (s *Set) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var decode []*Package
	if err := unmarshal(&decode); err != nil {
		return err
	}
	s.decode(decode)
	return nil
}
