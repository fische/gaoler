package yaml_test

import (
	"bytes"

	"github.com/fische/gaoler/config/formatter"
	. "github.com/fische/gaoler/config/formatter/modules/yaml"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Factory", func() {
	var (
		f *Factory
	)

	BeforeEach(func() {
		f = &Factory{}
	})

	Describe("NewEncoder", func() {
		var (
			e formatter.Encoder

			buf *bytes.Buffer
		)

		BeforeEach(func() {
			buf = bytes.NewBufferString("")
		})

		JustBeforeEach(func() {
			e = f.NewEncoder(buf)
		})

		It("should a valid Encoder", func() {
			obj, ok := e.(*Encoder)
			Expect(ok).To(BeTrue())
			Expect(obj).ToNot(BeNil())
		})
	})

	Describe("NewDecoder", func() {
		var (
			e formatter.Decoder

			buf *bytes.Buffer
		)

		BeforeEach(func() {
			buf = bytes.NewBufferString("")
		})

		JustBeforeEach(func() {
			e = f.NewDecoder(buf)
		})

		It("should a valid Decoder", func() {
			obj, ok := e.(*Decoder)
			Expect(ok).To(BeTrue())
			Expect(obj).ToNot(BeNil())
		})
	})

	Describe("Types", func() {
		var (
			arr []string
		)

		JustBeforeEach(func() {
			arr = f.Types()
		})

		It("should return YAML extensions", func() {
			Expect(arr).To(Equal([]string{"yaml", "yml"}))
		})
	})
})
