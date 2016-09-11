package yaml

import (
	"io"

	yaml "gopkg.in/yaml.v2"
)

type Encoder struct {
	w io.Writer
}

func (e *Encoder) Encode(obj interface{}) error {
	data, err := yaml.Marshal(obj)
	if err != nil {
		return err
	}
	_, err = e.w.Write(data)
	return err
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w}
}
