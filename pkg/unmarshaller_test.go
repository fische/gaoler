package pkg_test

import (
	"errors"
	"reflect"

	. "github.com/fische/gaoler/pkg"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Unmarshaller", func() {
	var (
		err error

		p *Package
	)

	BeforeEach(func() {
		p = &Package{}
	})

	Describe("UnmarshalJSON", func() {
		var data []byte

		JustBeforeEach(func() {
			err = p.UnmarshalJSON(data)
		})

		Context("With valid JSON", func() {
			BeforeEach(func() {
				data = []byte(`"test"`)
			})

			It("should not return an error", func() {
				Expect(err).To(BeNil())
			})

			It("should set Path and Saved correctly", func() {
				Expect(p.Path).To(Equal("test"))
				Expect(p.Saved).To(BeTrue())
			})
		})

		Context("With invalid JSON", func() {
			BeforeEach(func() {
				data = []byte(`"test`)
			})

			It("should return an error", func() {
				Expect(err).ToNot(BeNil())
			})
		})
	})

	Describe("UnmarshalYAML", func() {
		var unmarshaller func(interface{}) error

		JustBeforeEach(func() {
			err = p.UnmarshalYAML(unmarshaller)
		})

		Context("With valid YAML", func() {
			data := "test"

			BeforeEach(func() {
				unmarshaller = func(i interface{}) error {
					v := reflect.Indirect(reflect.ValueOf(i))
					v.Set(reflect.ValueOf(data))
					return nil
				}
			})

			It("should not return an error", func() {
				Expect(err).To(BeNil())
			})

			It("should set Path and Saved correctly", func() {
				Expect(p.Path).To(Equal(data))
				Expect(p.Saved).To(BeTrue())
			})
		})

		Context("With invalid YAML", func() {
			BeforeEach(func() {
				unmarshaller = func(i interface{}) error {
					return errors.New("")
				}
			})

			It("should return an error", func() {
				Expect(err).ToNot(BeNil())
			})
		})
	})
})
