package pkg_test

import (
	"errors"

	. "github.com/fische/gaoler/pkg"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Set", func() {
	var (
		s *Set
	)

	Describe("NewSet", func() {
		JustBeforeEach(func() {
			s = NewSet()
		})

		It("should initialize the map", func() {
			Expect(s.Packages).ToNot(BeNil())
		})
	})

	Describe("Methods", func() {
		BeforeEach(func() {
			s = &Set{
				Packages: map[string]*Package{
					"test": &Package{
						Path: "test",
					},
					"dir": &Package{
						Path: "dir",
					},
					"path": &Package{
						Path: "path",
					},
				},
			}
		})

		Describe("ListFrom", func() {
			var (
				err error

				srcPath string
			)

			JustBeforeEach(func() {
				err = s.ListFrom(srcPath)
			})

			Context("With a valid source path", func() {
				BeforeEach(func() {
					srcPath = "testdata"
				})

				Context("Without OnAdded callback defined", func() {
					It("should not return an error", func() {
						Expect(err).To(BeNil())
					})

					It("should have added package", func() {
						obj, ok := s.Packages["package"]
						Expect(ok).To(BeTrue())
						Expect(obj.Path).To(Equal("package"))
					})
				})

				Context("With failing OnAdded callback defined", func() {
					BeforeEach(func() {
						s.OnAdded = func(p *Package) (nextDirectory string, err error) {
							return "", errors.New("")
						}
					})

					It("should return an error", func() {
						Expect(err).ToNot(BeNil())
					})
				})

				Context("With valid OnAdded callback defined", func() {
					var passed bool

					BeforeEach(func() {
						passed = false
						s.OnAdded = func(p *Package) (nextDirectory string, err error) {
							Expect(p.Path).To(Equal("package"))
							passed = true
							return "testdata/vendor/package", nil
						}
					})

					It("should not return an error", func() {
						Expect(err).To(BeNil())
					})

					It("should have called OnAdded callback", func() {
						Expect(passed).To(BeTrue())
					})
				})
			})

			Context("With unknown source path", func() {
				BeforeEach(func() {
					srcPath = "package_nowhere"
				})

				It("should return an error", func() {
					Expect(err).ToNot(BeNil())
				})
			})
		})

		Describe("Insert", func() {
			var (
				added bool

				p     *Package
				force bool
			)

			JustBeforeEach(func() {
				added = s.Insert(p, force)
			})

			Context("With existing Package", func() {
				BeforeEach(func() {
					p = &Package{
						Path: "test",
					}
				})

				Context("With force enabled", func() {
					BeforeEach(func() {
						force = true
					})

					It("should return true", func() {
						Expect(added).To(BeTrue())
					})

					It("should add package to set", func() {
						obj, ok := s.Packages[p.Path]
						Expect(ok).To(BeTrue())
						Expect(obj).To(Equal(p))
					})
				})

				Context("With force disabled", func() {
					BeforeEach(func() {
						force = false
					})

					It("should return false", func() {
						Expect(added).To(BeFalse())
					})
				})
			})

			Context("With non existing Package", func() {
				BeforeEach(func() {
					force = false
					p = &Package{
						Path: "nonexisting",
					}
				})

				It("should return true", func() {
					Expect(added).To(BeTrue())
				})

				It("should add package to set", func() {
					obj, ok := s.Packages[p.Path]
					Expect(ok).To(BeTrue())
					Expect(obj).To(Equal(p))
				})
			})
		})

		Describe("ForEach", func() {
			var (
				err error

				fc func(key string, value *Package) error
			)

			JustBeforeEach(func() {
				err = s.ForEach(fc)
			})

			Context("When callback does not fail", func() {
				var (
					count int
				)

				BeforeEach(func() {
					count = 0
					fc = func(key string, value *Package) error {
						_, ok := s.Packages[key]
						Expect(value).ToNot(BeNil())
						Expect(ok).To(BeTrue())
						delete(s.Packages, key)
						count++
						return nil
					}
				})

				It("should not return an error", func() {
					Expect(err).To(BeNil())
				})

				It("should walk through all packages", func() {
					Expect(count).To(Equal(3))
				})
			})

			Context("When callback fails", func() {
				BeforeEach(func() {
					fc = func(key string, value *Package) error {
						return errors.New("")
					}
				})

				It("should return error", func() {
					Expect(err).ToNot(BeNil())
				})
			})
		})

		Describe("Remove", func() {
			var key string

			JustBeforeEach(func() {
				s.Remove(key)
			})

			BeforeEach(func() {
				key = "test"
			})

			It("should remove key from Packages map", func() {
				obj, ok := s.Packages[key]
				Expect(ok).To(BeFalse())
				Expect(obj).To(BeNil())
			})
		})
	})
})
