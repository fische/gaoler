package dependency

import "encoding/json"

func (s *Set) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &s.Deps); err != nil {
		return err
	} else if s.OnDecoded != nil {
		for rootPackage, dep := range s.Deps {
			dep.RootPackage = rootPackage
			if err = s.OnDecoded(dep); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Set) UnmarshalYAML(unmarshal func(interface{}) error) error {
	if err := unmarshal(&s.Deps); err != nil {
		return err
	} else if s.OnDecoded != nil {
		for rootPackage, dep := range s.Deps {
			dep.RootPackage = rootPackage
			if err = s.OnDecoded(dep); err != nil {
				return err
			}
		}
	}
	return nil
}
