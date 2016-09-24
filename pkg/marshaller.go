package pkg

func (p Package) MarshalJSON() ([]byte, error) {
	return []byte("\"" + p.path + "\""), nil
}

func (p Package) MarshalYAML() (interface{}, error) {
	return p.path, nil
}
