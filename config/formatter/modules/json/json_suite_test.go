package json_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

type testStruct struct {
	Test string
}

func TestJson(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "JSON Suite")
}
