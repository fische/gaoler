package dependency

import "encoding/json"

func (s Set) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Deps)
}

func (s *Set) MarshalYAML() (interface{}, error) {
	return s.Deps, nil
}
