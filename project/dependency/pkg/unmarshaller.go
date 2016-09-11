package pkg

import (
	"errors"
	"go/build"
)

func (p *Package) importFromPath() {
	i, err := build.Import(p.Path, srcPath, build.FindOnly|build.AllowBinary)
	if err == nil {
		p.Dir = i.Dir
		p.Root = i.Goroot
	}
}

func (p *Package) UnmarshalJSON(data []byte) error {
	if data[0] == '"' && data[len(data)-1] == '"' {
		p.Path = string(data[1 : len(data)-1])
		p.importFromPath()
		return nil
	}
	return errors.New("Could not unmarshal package")
}

func (p *Package) UnmarshalYAML(unmarshal func(interface{}) error) error {
	err := unmarshal(&p.Path)
	if err != nil {
		return err
	}
	p.importFromPath()
	return nil
}
