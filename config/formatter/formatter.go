package formatter

import "io"

type IndentableEncoder interface {
	SetIndent(prefix, indent string)
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
