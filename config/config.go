package config

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/fische/gaoler/config/formatter"
	"github.com/fische/gaoler/config/formatter/modules"
	"github.com/fische/gaoler/project"
)

const (
	openPerm = 0664
)

var (
	file   *os.File
	format formatter.Factory
)

func Save(project *project.Project) error {
	e := format.NewEncoder(file)
	if i, ok := format.(formatter.IndentableEncoder); ok {
		i.SetIndent("", "\t")
	}
	return e.Encode(project)
}

func Load(p *project.Project) error {
	return format.NewDecoder(file).Decode(p)
}

func Setup(configPath string, force bool) (err error) {
	var flag int
	if force {
		flag |= os.O_TRUNC
	}
	file, err = os.OpenFile(configPath, os.O_RDWR|os.O_CREATE|flag, openPerm)
	if err != nil {
		return
	}
	ext := filepath.Ext(file.Name())
	if len(ext) > 0 && ext[0] == '.' {
		ext = ext[1:]
	}
	if format = modules.GetFormatter(ext); format == nil {
		err = errors.New("Could not find formatter")
	}
	return
}

func Close() error {
	return file.Close()
}
