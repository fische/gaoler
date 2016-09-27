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
	file    *os.File
	format  formatter.Factory
	project *project.Project
}

const (
	openPerm = 0664
)

func New(p *project.Project, configPath string, flags Flags) (*Config, error) {
	var (
		cfg = &Config{
			project: p,
		}
		err error
	)
	if cfg.file, err = os.OpenFile(configPath, flags.OpenFlags(), openPerm); err != nil {
		return nil, err
	}
	ext := filepath.Ext(cfg.file.Name())
	if len(ext) > 0 && ext[0] == '.' {
		ext = ext[1:]
	}
	if cfg.format = modules.GetFormatter(ext); cfg.format == nil {
		return nil, errors.New("Could not find formatter")
	}
	return cfg, nil
}

func (cfg Config) Save() error {
	if offset, err := cfg.file.Seek(0, 0); err != nil {
		return err
	} else if offset != 0 {
		return errors.New("Could not seek beginning of the config")
	} else if err = cfg.file.Truncate(0); err != nil {
		return err
	}
	e := cfg.format.NewEncoder(cfg.file)
	if i, ok := e.(formatter.PrettyEncoder); ok {
		return i.PrettyEncode(cfg.project)
	}
	return e.Encode(cfg.project)
}

func (cfg *Config) Load() error {
	return cfg.format.NewDecoder(cfg.file).Decode(cfg.project)
}
