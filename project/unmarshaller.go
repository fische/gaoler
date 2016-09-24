package project

import (
	"encoding/json"

	"github.com/fische/gaoler/project/dependency"
)

func (deps Dependencies) UnmarshalJSON(data []byte) error {
	decoded := make(map[string]json.RawMessage)
	for k, v := range decoded {
		dep := &dependency.Dependency{}
		if err := json.Unmarshal(v, dep); err != nil {
			return err
		}
		deps[k] = dep
		dep.SetRootPackage(k)
	}
	return nil
}

// func (deps Dependencies) UnmarshalYAML(unmarshal func(interface{}) error) error {
// 	if err := unmarshal(deps); err != nil {
// 		return err
// 	}
// 	for k, v := range deps {
// 		v.SetRootPackage(k)
// 	}
// 	return nil
// }
