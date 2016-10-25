package pkg

import "encoding/json"

func (p Package) MarshalJSON() ([]byte, error) {
	return []byte("\"" + p.path + "\""), nil
}

func (p Package) MarshalYAML() (interface{}, error) {
	return p.path, nil
}

func (s Set) encode() []*Package {
	var packages []*Package
	for _, p := range s.packages {
		packages = append(packages, p)
	}
	return packages
}

func (s Set) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.encode())
}

func (s Set) MarshalYAML() (interface{}, error) {
	return s.encode(), nil
}
