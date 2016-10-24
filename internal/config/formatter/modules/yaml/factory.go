package yaml

import (
	"io"

	"github.com/fische/gaoler/internal/config/formatter"
)

type Factory struct{}

const (
	typeName = "yaml"
)

func (f Factory) NewEncoder(w io.Writer) formatter.Encoder {
	return NewEncoder(w)
}

func (f Factory) NewDecoder(r io.Reader) formatter.Decoder {
	return NewDecoder(r)
}

func (f Factory) Types() []string {
	return []string{typeName, "yml"}
}
