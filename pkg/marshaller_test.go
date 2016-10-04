package pkg_test

import (
	"fmt"

	. "github.com/fische/gaoler/pkg"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Marshaller", func() {
	var (
		err error

		p Package
	)

	BeforeEach(func() {
		p = Package{
			Path: "testpath",
		}
	})

	Describe("MarshalJSON", func() {
		var data []byte

		JustBeforeEach(func() {
			data, err = p.MarshalJSON()
		})

		It("should not return an error", func() {
			Expect(err).To(BeNil())
		})

		It("should return path surronded by double quotes", func() {
			Expect(string(data)).To(Equal(fmt.Sprintf(`"%s"`, p.Path)))
		})
	})

	Describe("MarshalYAML", func() {
		var data interface{}

		JustBeforeEach(func() {
			data, err = p.MarshalYAML()
		})

		It("should not return an error", func() {
			Expect(err).To(BeNil())
		})

		It("should return path ended by a carriage return", func() {
			Expect(data).To(Equal(p.Path))
		})
	})
})
