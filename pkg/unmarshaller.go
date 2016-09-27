package pkg

import "errors"

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
