package pkg

func (p Package) MarshalJSON() ([]byte, error) {
	return []byte("\"" + p.Path + "\""), nil
}

func (p Package) MarshalYAML() (interface{}, error) {
	return p.Path, nil
}
