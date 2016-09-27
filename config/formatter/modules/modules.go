package modules

import (
	"github.com/fische/gaoler/config/formatter"
	"github.com/fische/gaoler/config/formatter/modules/json"
	"github.com/fische/gaoler/config/formatter/modules/yaml"
)

var impl []formatter.Factory

func Register(f formatter.Factory) {
	impl = append(impl, f)
}

func GetFormatter(ext string) formatter.Factory {
	for _, f := range impl {
		for _, t := range f.Types() {
			if t == ext {
				return f
			}
		}
	}
	return nil
}

func init() {
	// Register modules here
	impl = []formatter.Factory{
		&json.Factory{},
		&yaml.Factory{},
	}
}
