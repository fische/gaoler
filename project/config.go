package project

import (
	"encoding/json"
	"os"
)

const (
	openPerm = 0664
)

func (project Project) Save(configPath string) error {
	file, err := os.OpenFile(configPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, openPerm)
	if err != nil {
		return err
	}
	e := json.NewEncoder(file)
	e.SetIndent("", "\t")
	return e.Encode(project)
}
