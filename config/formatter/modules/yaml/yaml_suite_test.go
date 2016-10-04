package yaml_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

type testStruct struct {
	Test string
}

type testReader struct {
	read func(p []byte) (n int, err error)
}

type testWriter struct {
	write func(p []byte) (n int, err error)
}

func (r testReader) Read(p []byte) (n int, err error) {
	return r.read(p)
}

func (r testWriter) Write(p []byte) (n int, err error) {
	return r.write(p)
}

func TestYaml(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "YAML Suite")
}
