package modules_test

import (
	"io"

	"github.com/fische/gaoler/config/formatter"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

type factory struct {
	types []string
}

func (f factory) Types() []string {
	return f.types
}

func (f factory) NewEncoder(w io.Writer) formatter.Encoder {
	return nil
}

func (f factory) NewDecoder(r io.Reader) formatter.Decoder {
	return nil
}

func TestModules(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Modules Suite")
}
