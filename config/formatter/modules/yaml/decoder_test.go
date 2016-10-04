package yaml_test

import (
	"bytes"
	"errors"
	"io"
	"log"
	"os"

	. "github.com/fische/gaoler/config/formatter/modules/yaml"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Decoder", func() {
	var (
		d *Decoder

		buf io.Reader
	)

	Describe("NewDecoder", func() {
		BeforeEach(func() {
			buf = bytes.NewBufferString("")
		})

		JustBeforeEach(func() {
			d = NewDecoder(buf)
		})

		It("should a valid Decoder", func() {
			Expect(d.Reader).To(Equal(buf))
		})
	})

	Describe("Methods", func() {
		JustBeforeEach(func() {
			d = &Decoder{buf}
		})

		Describe("Decode", func() {
			var (
				err error

				obj *testStruct
			)

			JustBeforeEach(func() {
				err = d.Decode(obj)
			})

			Context("With an empty YAML", func() {
				BeforeEach(func() {
					buf = bytes.NewBufferString("")
				})

				It("should return an error", func() {
					Expect(err).ToNot(BeNil())
				})
			})

			Context("When reading returns an error", func() {
				BeforeEach(func() {
					buf = testReader{
						read: func(p []byte) (n int, err error) {
							err = errors.New("")
							return
						},
					}
				})

				It("should return an error", func() {
					Expect(err).ToNot(BeNil())
				})
			})

			Context("With valid YAML", func() {
				BeforeEach(func() {
					buf, err = os.Open("testdata/test.yml")
					if err != nil {
						log.Fatalf("Could not open testdata/test.yml: %v", err)
					}
					obj = &testStruct{}
				})

				It("should not return an error", func() {
					Expect(err).To(BeNil())
				})

				It("should fill obj correctly", func() {
					Expect(obj.Test).To(Equal("test"))
				})
			})
		})
	})
})
