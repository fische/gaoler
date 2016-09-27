package dependency

import "encoding/json"

func (s *Set) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &s.deps); err != nil {
		return err
	} else if s.OnDecoded != nil {
		for rootPackage, dep := range s.deps {
			dep.rootPackage = rootPackage
			if err = s.OnDecoded(dep); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Set) UnmarshalYAML(unmarshal func(interface{}) error) error {
	if err := unmarshal(&s.deps); err != nil {
		return err
	} else if s.OnDecoded != nil {
		for rootPackage, dep := range s.deps {
			dep.rootPackage = rootPackage
			if err = s.OnDecoded(dep); err != nil {
				return err
			}
		}
	}
	return nil
}
