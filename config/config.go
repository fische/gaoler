package config

import (
	"encoding/json"
	"os"

	"github.com/fische/gaoler/project"
)

const (
	openPerm = 0664
)

var (
	file *os.File
)

func Save(project *project.Project) error {
	e := json.NewEncoder(file)
	e.SetIndent("", "\t")
	return e.Encode(project)
}

func Load(p *project.Project) error {
	return json.NewDecoder(file).Decode(p)
}

func Setup(configPath string, force bool) (err error) {
	var flag int
	if force {
		flag |= os.O_TRUNC
	}
	file, err = os.OpenFile(configPath, os.O_RDWR|os.O_CREATE|flag, openPerm)
	return
}

func Close() error {
	return file.Close()
}
