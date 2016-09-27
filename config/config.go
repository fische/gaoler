package config

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/fische/gaoler/config/formatter"
	"github.com/fische/gaoler/config/formatter/modules"
	"github.com/fische/gaoler/project"
)

type Config struct {
	path    string
	format  formatter.Factory
	project *project.Project
}

const (
	openPerm = 0664
)

func New(p *project.Project, configPath string) (*Config, error) {
	var cfg = &Config{
		project: p,
		path:    configPath,
	}
	ext := filepath.Ext(configPath)
	if len(ext) > 0 && ext[0] == '.' {
		ext = ext[1:]
	}
	if cfg.format = modules.GetFormatter(ext); cfg.format == nil {
		return nil, errors.New("Could not find formatter")
	}
	return cfg, nil
}

func (cfg Config) Save() error {
	var (
		file *os.File
		err  error
	)
	if file, err = os.OpenFile(cfg.path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, openPerm); err != nil {
		return err
	}
	e := cfg.format.NewEncoder(file)
	if i, ok := e.(formatter.PrettyEncoder); ok {
		return i.PrettyEncode(cfg.project)
	} else if err = e.Encode(cfg.project); err != nil {
		return err
	}
	return file.Close()
}

func (cfg *Config) Load() error {
	var (
		file *os.File
		err  error
	)
	if file, err = os.OpenFile(cfg.path, os.O_RDONLY, openPerm); err != nil {
		return err
	} else if err = cfg.format.NewDecoder(file).Decode(cfg.project); err != nil {
		return err
	}
	return file.Close()
}
