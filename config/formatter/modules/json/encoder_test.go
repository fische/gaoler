package json_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"

	. "github.com/fische/gaoler/config/formatter/modules/json"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Encoder", func() {
	var (
		e *Encoder

		buf *bytes.Buffer
	)

	BeforeEach(func() {
		buf = bytes.NewBufferString("")
		e = &Encoder{
			Encoder: json.NewEncoder(buf),
		}
	})

	Describe("PrettyEncode", func() {
		var (
			err error

			obj *testStruct
		)

		BeforeEach(func() {
			obj = &testStruct{
				Test: "test",
			}
		})

		JustBeforeEach(func() {
			err = e.PrettyEncode(obj)
		})

		Context("When encoding succeeds", func() {
			var (
				out []byte
			)

			BeforeEach(func() {
				out, err = ioutil.ReadFile("testdata/test.json")
				if err != nil {
					log.Fatalf("Could not read testdata/test.json : %v", err)
				}
			})

			It("should not return an error", func() {
				Expect(err).To(BeNil())
			})

			It("should return an indented JSON", func() {
				Expect(buf.String()).To(Equal(string(out)))
			})
		})
	})
})
