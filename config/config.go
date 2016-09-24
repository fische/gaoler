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
	err := file.Truncate(0)
	if err != nil {
		return err
	}
	e := format.NewEncoder(file)
	if i, ok := e.(formatter.PrettyEncoder); ok {
		return i.PrettyEncode(project)
	}
	return e.Encode(project)
}

func Load(p *project.Project) error {
	return format.NewDecoder(file).Decode(p)
}

func Setup(configPath string, flags Flags) (err error) {
	file, err = os.OpenFile(configPath, flags.OpenFlags(), openPerm)
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
