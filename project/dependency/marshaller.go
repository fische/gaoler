package dependency

import "encoding/json"

func (s Set) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.deps)
}

func (s *Set) MarshalYAML() (interface{}, error) {
	return s.deps, nil
}
