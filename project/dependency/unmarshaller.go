package dependency

import "encoding/json"

func (s Set) UnmarshalJSON(data []byte) error {
	m := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	for root, value := range m {
		dep := &Dependency{
			RootPackage: root,
		}
		if err := json.Unmarshal(value, dep); err != nil {
			return err
		}
		s[root] = dep
	}
	return nil
}

func (s Set) UnmarshalYAML(unmarshal func(interface{}) error) error {
	m := make(map[string]*Dependency)
	err := unmarshal(&m)
	if err != nil {
		return err
	}
	for root, dep := range m {
		dep.RootPackage = root
		s[root] = dep
	}
	return nil
}
