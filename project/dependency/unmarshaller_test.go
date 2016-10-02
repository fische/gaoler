package dependency_test

import (
	"errors"
	"reflect"

	"github.com/fische/gaoler/pkg"
	. "github.com/fische/gaoler/project/dependency"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Unmarshaller", func() {
	var (
		err error

		s *Set
	)

	BeforeEach(func() {
		s = &Set{}
	})

	Describe("UnmarshalJSON", func() {
		var (
			data []byte
		)

		JustBeforeEach(func() {
			err = s.UnmarshalJSON(data)
		})

		Context("With a valid json", func() {
			BeforeEach(func() {
				data = []byte(`{"test":{"Packages":["test"]}}`)
			})

			Context("With a valid OnDecoded callback", func() {
				var passed bool

				BeforeEach(func() {
					passed = false
					s.OnDecoded = func(dep *Dependency) error {
						passed = true
						return nil
					}
				})

				It("should not return an error", func() {
					Expect(err).To(BeNil())
				})

				It("should have called OnDecoded callback", func() {
					Expect(passed).To(BeTrue())
				})

				It("should have filled dependencies", func() {
					Expect(s.Deps).To(Equal(map[string]*Dependency{
						"test": &Dependency{
							RootPackage: "test",
							Packages: []*pkg.Package{
								&pkg.Package{
									Path:  "test",
									Saved: true,
								},
							},
						},
					}))
				})
			})

			Context("With a failing OnDecoded callback", func() {
				BeforeEach(func() {
					s.OnDecoded = func(dep *Dependency) error {
						return errors.New("")
					}
				})

				It("should return an error", func() {
					Expect(err).ToNot(BeNil())
				})

			})
		})

		Context("With an invalid json", func() {
			BeforeEach(func() {
				data = []byte(`["test"]`)
			})

			It("should return an error", func() {
				Expect(err).ToNot(BeNil())
			})
		})
	})

	Describe("UnmarshalYAML", func() {
		var (
			unmarshaller func(interface{}) error
		)

		JustBeforeEach(func() {
			err = s.UnmarshalYAML(unmarshaller)
		})

		Context("With a valid unmarshaller", func() {
			BeforeEach(func() {
				unmarshaller = func(i interface{}) error {
					v := reflect.Indirect(reflect.ValueOf(i))
					v.Set(reflect.ValueOf(map[string]*Dependency{
						"test": &Dependency{
							RootPackage: "test",
							Packages: []*pkg.Package{
								&pkg.Package{
									Path:  "test",
									Saved: true,
								},
							},
						},
					}))
					return nil
				}
			})

			Context("With a valid OnDecoded callback", func() {
				var passed bool

				BeforeEach(func() {
					passed = false
					s.OnDecoded = func(dep *Dependency) error {
						passed = true
						return nil
					}
				})

				It("should not return an error", func() {
					Expect(err).To(BeNil())
				})

				It("should have called OnDecoded callback", func() {
					Expect(passed).To(BeTrue())
				})

				It("should have filled dependencies", func() {
					Expect(s.Deps).To(Equal(map[string]*Dependency{
						"test": &Dependency{
							RootPackage: "test",
							Packages: []*pkg.Package{
								&pkg.Package{
									Path:  "test",
									Saved: true,
								},
							},
						},
					}))
				})
			})

			Context("With a failing OnDecoded callback", func() {
				BeforeEach(func() {
					s.OnDecoded = func(dep *Dependency) error {
						return errors.New("")
					}
				})

				It("should return an error", func() {
					Expect(err).ToNot(BeNil())
				})

			})
		})

		Context("With an invalid json", func() {
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
