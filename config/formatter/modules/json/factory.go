package json

import (
	"encoding/json"
	"io"

	"github.com/fische/gaoler/config/formatter"
)

type Factory struct{}

const (
	typeName = "json"
)

func (f Factory) NewEncoder(w io.Writer) formatter.Encoder {
	return json.NewEncoder(w)
}

func (f Factory) NewDecoder(r io.Reader) formatter.Decoder {
	return json.NewDecoder(r)
}

func (f Factory) Types() []string {
	return []string{typeName}
}
