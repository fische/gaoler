package dependency_test

import (
	"github.com/fische/gaoler/pkg"
	. "github.com/fische/gaoler/project/dependency"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Marshaller", func() {
	var (
		err error

		s *Set
	)

	BeforeEach(func() {
		s = &Set{
			Deps: map[string]*Dependency{
				"test": &Dependency{
					RootPackage: "test",
					Packages: []*pkg.Package{
						&pkg.Package{
							Path: "test",
						},
					},
				},
			},
		}
	})

	Describe("MarshalJSON", func() {
		var (
			data []byte
		)

		JustBeforeEach(func() {
			data, err = s.MarshalJSON()
		})

		It("should not return an error", func() {
			Expect(err).To(BeNil())
		})

		It("should only marshal deps map", func() {
			Expect(data).To(Equal([]byte(`{"test":{"Packages":["test"]}}`)))
		})
	})

	Describe("MarshalYAML", func() {
		var (
			data interface{}
		)

		JustBeforeEach(func() {
			data, err = s.MarshalYAML()
		})

		It("should not return an error", func() {
			Expect(err).To(BeNil())
		})

		It("should only marshal deps map", func() {
			Expect(data).To(Equal(s.Deps))
		})
	})
})
