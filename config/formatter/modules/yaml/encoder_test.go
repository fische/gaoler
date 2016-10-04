package yaml_test

import (
	"bytes"
	"io/ioutil"
	"log"

	. "github.com/fische/gaoler/config/formatter/modules/yaml"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Encoder", func() {
	var (
		d *Encoder

		buf *bytes.Buffer
	)

	Describe("NewEncoder", func() {
		BeforeEach(func() {
			buf = bytes.NewBufferString("")
		})

		JustBeforeEach(func() {
			d = NewEncoder(buf)
		})

		It("should a valid Encoder", func() {
			Expect(d.Writer).To(Equal(buf))
		})
	})

	Describe("Methods", func() {
		JustBeforeEach(func() {
			buf = bytes.NewBufferString("")
			d = &Encoder{buf}
		})

		Describe("Encode", func() {
			var (
				err error

				obj *testStruct
			)

			JustBeforeEach(func() {
				err = d.Encode(obj)
			})

			Context("With valid YAML", func() {
				var (
					data []byte
				)

				BeforeEach(func() {
					obj = &testStruct{
						Test: "test",
					}
					data, err = ioutil.ReadFile("testdata/test.yml")
					if err != nil {
						log.Fatalf("Could not read testdata/test.yml : %v", err)
					}
				})

				It("should not return an error", func() {
					Expect(err).To(BeNil())
				})

				It("should fill obj correctly", func() {
					Expect(buf.Bytes()).To(Equal(data))
				})
			})
		})
	})
})
