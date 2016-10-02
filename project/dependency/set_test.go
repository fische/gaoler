package dependency_test

import (
	"errors"
	"fmt"

	"github.com/fische/gaoler/pkg"
	. "github.com/fische/gaoler/project/dependency"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Set", func() {
	Describe("NewSet", func() {
		var (
			s *Set
		)

		JustBeforeEach(func() {
			s = NewSet()
		})

		It("should initialize dependencies map", func() {
			Expect(s.Deps).ToNot(BeNil())
		})
	})

	Describe("Methods", func() {
		var (
			s *Set
		)

		BeforeEach(func() {
			s = &Set{
				Deps: make(map[string]*Dependency),
			}
		})

		Describe("MergePackageSet", func() {
			var (
				err error

				o *pkg.Set
			)

			JustBeforeEach(func() {
				err = s.MergePackageSet(o)
			})

			Context("With new packages", func() {
				BeforeEach(func() {
					o = &pkg.Set{
						Packages: map[string]*pkg.Package{
							"test/subpackage": &pkg.Package{
								Path: "test/subpackage",
							},
							"test": &pkg.Package{
								Path: "test",
							},
							"dep": &pkg.Package{
								Path: "dep",
							},
							"dep/subdir": &pkg.Package{
								Path: "dep/subdir",
							},
						},
					}
				})

				Context("With valid OnPackageAdded callback", func() {
					var passed bool

					BeforeEach(func() {
						passed = false
						s.OnPackageAdded = func(p *pkg.Package, dep *Dependency) error {
							passed = true
							return nil
						}
					})

					It("should not return an error", func() {
						Expect(err).To(BeNil())
					})

					It("should have set dependency set", func() {
						Expect(len(s.Deps)).To(Equal(2))
						obj, ok := s.Deps["test"]
						Expect(ok).To(BeTrue())
						Expect(obj.RootPackage).To(Equal("test"))
						for _, p := range obj.Packages {
							if p.Path == "test" {
								Expect(p).To(Equal(&pkg.Package{
									Path: "test",
								}))
							} else if p.Path == "test/subpackage" {
								Expect(p).To(Equal(&pkg.Package{
									Path: "test/subpackage",
								}))
							} else {
								Fail(fmt.Sprintf("It should only have packages test and test/subpackage : %v", obj))
							}
						}
						obj, ok = s.Deps["dep"]
						Expect(ok).To(BeTrue())
						Expect(obj.RootPackage).To(Equal("dep"))
						for _, p := range obj.Packages {
							if p.Path == "dep" {
								Expect(p).To(Equal(&pkg.Package{
									Path: "dep",
								}))
							} else if p.Path == "dep/subdir" {
								Expect(p).To(Equal(&pkg.Package{
									Path: "dep/subdir",
								}))
							} else {
								Fail(fmt.Sprintf("It should only have packages dep and dep/subdir : %v", obj))
							}
						}
					})

					It("should have called OnPackageAdded callback", func() {
						Expect(passed).To(BeTrue())
					})
				})

				Context("With failing OnPackageAdded callback", func() {
					BeforeEach(func() {
						s.OnPackageAdded = func(p *pkg.Package, dep *Dependency) error {
							return errors.New("")
						}
					})

					It("should not return an error", func() {
						Expect(err).ToNot(BeNil())
					})
				})
			})

			Context("With subpackages from packages from dependencies set", func() {
				BeforeEach(func() {
					s.Deps = map[string]*Dependency{
						"test": &Dependency{
							RootPackage: "test",
							Packages: []*pkg.Package{
								&pkg.Package{
									Path: "test",
								},
							},
						},
						"dep/subdir": &Dependency{
							RootPackage: "dep/subdir",
							Packages: []*pkg.Package{
								&pkg.Package{
									Path: "dep/subdir",
								},
							},
						},
					}
					o = &pkg.Set{
						Packages: map[string]*pkg.Package{
							"test/subpackage": &pkg.Package{
								Path: "test/subpackage",
							},
							"dep": &pkg.Package{
								Path: "dep",
							},
						},
					}
				})

				Context("With valid OnPackageAdded callback", func() {
					var passed bool

					BeforeEach(func() {
						passed = false
						s.OnPackageAdded = func(p *pkg.Package, dep *Dependency) error {
							passed = true
							return nil
						}
					})

					It("should not return an error", func() {
						Expect(err).To(BeNil())
					})

					It("should have set dependency set", func() {
						Expect(s.Deps).To(Equal(map[string]*Dependency{
							"test": &Dependency{
								RootPackage: "test",
								Packages: []*pkg.Package{
									&pkg.Package{
										Path: "test",
									},
									&pkg.Package{
										Path: "test/subpackage",
									},
								},
							},
							"dep": &Dependency{
								RootPackage: "dep",
								Packages: []*pkg.Package{
									&pkg.Package{
										Path: "dep/subdir",
									},
									&pkg.Package{
										Path: "dep",
									},
								},
							},
						}))
					})

					It("should have called OnPackageAdded callback", func() {
						Expect(passed).To(BeTrue())
					})
				})

				Context("With failing OnPackageAdded callback", func() {
					BeforeEach(func() {
						s.OnPackageAdded = func(p *pkg.Package, dep *Dependency) error {
							return errors.New("")
						}
					})

					It("should not return an error", func() {
						Expect(err).ToNot(BeNil())
					})
				})
			})
		})

		Describe("ToPackageSet", func() {
			var (
				o *pkg.Set
			)

			BeforeEach(func() {
				s.Deps = map[string]*Dependency{
					"test": &Dependency{
						RootPackage: "test",
						Packages: []*pkg.Package{
							&pkg.Package{
								Path: "test",
							},
							&pkg.Package{
								Path: "test/subpackage",
							},
						},
					},
					"dep": &Dependency{
						RootPackage: "dep",
						Packages: []*pkg.Package{
							&pkg.Package{
								Path: "dep/subdir",
							},
							&pkg.Package{
								Path: "dep",
							},
						},
					},
				}
			})

			JustBeforeEach(func() {
				o = s.ToPackageSet()
			})

			It("should fill package map correctly", func() {
				Expect(o.Packages).To(Equal(map[string]*pkg.Package{
					"test": &pkg.Package{
						Path: "test",
					},
					"test/subpackage": &pkg.Package{
						Path: "test/subpackage",
					},
					"dep/subdir": &pkg.Package{
						Path: "dep/subdir",
					},
					"dep": &pkg.Package{
						Path: "dep",
					},
				}))
			})
		})
	})
})
