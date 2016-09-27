package formatter

import "io"

type PrettyEncoder interface {
	PrettyEncode(i interface{}) error
}

type Encoder interface {
	Encode(i interface{}) error
}

type Decoder interface {
	Decode(i interface{}) error
}

type Factory interface {
	Types() []string
	NewEncoder(w io.Writer) Encoder
	NewDecoder(r io.Reader) Decoder
}
