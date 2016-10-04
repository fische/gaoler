package json_test

import (
	"bytes"
	"encoding/json"

	"github.com/fische/gaoler/config/formatter"
	. "github.com/fische/gaoler/config/formatter/modules/json"

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
			Expect(obj.Encoder).ToNot(BeNil())
		})

		It("should implement formatter.PrettyEncoder", func() {
			_, ok := e.(formatter.PrettyEncoder)
			Expect(ok).To(BeTrue())
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
			obj, ok := e.(*json.Decoder)
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

		It("should return JSON extensions", func() {
			Expect(arr).To(Equal([]string{"json"}))
		})
	})
})
