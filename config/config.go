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
	File    *os.File
	Format  formatter.Factory
	Project *project.Project
}

const (
	openPerm = 0664
)

func New(p *project.Project, configPath string, flags Flags) (*Config, error) {
	var (
		cfg = &Config{
			Project: p,
		}
		err error
	)
	if cfg.File, err = os.OpenFile(configPath, flags.OpenFlags(), openPerm); err != nil {
		return nil, err
	}
	ext := filepath.Ext(cfg.File.Name())
	if len(ext) > 0 && ext[0] == '.' {
		ext = ext[1:]
	}
	if cfg.Format = modules.GetFormatter(ext); cfg.Format == nil {
		return nil, errors.New("Could not find formatter")
	}
	return cfg, nil
}

func (cfg Config) Save() error {
	if offset, err := cfg.File.Seek(0, 0); err != nil {
		return err
	} else if offset != 0 {
		return errors.New("Could not seek beginning of the config")
	} else if err = cfg.File.Truncate(0); err != nil {
		return err
	}
	e := cfg.Format.NewEncoder(cfg.File)
	if i, ok := e.(formatter.PrettyEncoder); ok {
		return i.PrettyEncode(cfg.Project)
	}
	return e.Encode(cfg.Project)
}

func (cfg *Config) Load() error {
	return cfg.Format.NewDecoder(cfg.File).Decode(cfg.Project)
}
