package config_test

import (
	"io"

	"github.com/fische/gaoler/config/formatter"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

type factory struct {
	types   []string
	encoder formatter.Encoder
	decoder formatter.Decoder
}

type encoder struct {
	encode func(i interface{}) error
}

type prettyEncoder struct {
	encoder
	prettyEncode func(i interface{}) error
}

type decoder struct {
	decode func(i interface{}) error
}

func (f factory) Types() []string {
	return f.types
}

func (f factory) NewEncoder(w io.Writer) formatter.Encoder {
	return f.encoder
}

func (f factory) NewDecoder(r io.Reader) formatter.Decoder {
	return f.decoder
}

func (e encoder) Encode(i interface{}) error {
	return e.encode(i)
}

func (e prettyEncoder) PrettyEncode(i interface{}) error {
	return e.prettyEncode(i)
}

func (d decoder) Decode(i interface{}) error {
	return d.decode(i)
}

func TestConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Config Suite")
}
