package config

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/fische/gaoler/internal/config/formatter"
	"github.com/fische/gaoler/internal/config/formatter/modules"
	"github.com/fische/gaoler/project"
)

type Config struct {
	file   *os.File
	format formatter.Factory
}

const (
	openPerm = 0664
)

func New(configPath string, flags Flags) (*Config, error) {
	var (
		cfg = &Config{}
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

func (cfg Config) Save(p *project.Project) error {
	if offset, err := cfg.File.Seek(0, 0); err != nil {
		return err
	} else if offset != 0 {
		return errors.New("Could not seek beginning of the config")
	} else if err = cfg.File.Truncate(0); err != nil {
		return err
	}
	e := cfg.Format.NewEncoder(cfg.File)
	if i, ok := e.(formatter.PrettyEncoder); ok {
		return i.PrettyEncode(p)
	}
	return e.Encode(p)
}

func (cfg *Config) Load(p *project.Project) error {
	return cfg.Format.NewDecoder(cfg.File).Decode(p)
}
