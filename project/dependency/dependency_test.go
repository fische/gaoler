package dependency_test

import (
	"github.com/fische/gaoler/pkg"
	. "github.com/fische/gaoler/project/dependency"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Dependency", func() {
	Describe("New", func() {
		var (
			dep *Dependency

			p *pkg.Package
		)

		BeforeEach(func() {
			p = &pkg.Package{
				Path: "test",
			}
		})

		JustBeforeEach(func() {
			dep = New(p)
		})

		It("should fill in the fields correctly", func() {
			Expect(dep.RootPackage).To(Equal(p.Path))
			Expect(dep.Packages[0]).To(Equal(p))
		})
	})

	Describe("Methods", func() {
		var (
			dep *Dependency
		)

		BeforeEach(func() {
			dep = &Dependency{}
		})

		Describe("Add", func() {
			var (
				added bool

				p *pkg.Package
			)

			JustBeforeEach(func() {
				added = dep.Add(p)
			})

			Context("With existing package", func() {
				BeforeEach(func() {
					p = &pkg.Package{
						Path: "test",
					}
					dep.Packages = append(dep.Packages, p)
				})

				It("should return false", func() {
					Expect(added).To(BeFalse())
				})

				It("should not add package", func() {
					Expect(dep.Packages).To(Equal([]*pkg.Package{p}))
				})
			})

			Context("With non existing package", func() {
				BeforeEach(func() {
					p = &pkg.Package{
						Path: "test",
					}
				})

				It("should return true", func() {
					Expect(added).To(BeTrue())
				})

				It("should add package", func() {
					Expect(dep.Packages).To(Equal([]*pkg.Package{p}))
				})
			})
		})

		Describe("IsVendored", func() {
			var (
				vendored bool
			)

			BeforeEach(func() {
				dep.Packages = []*pkg.Package{
					&pkg.Package{
						Path: "test",
					},
					&pkg.Package{
						Path: "third",
					},
					&pkg.Package{
						Path: "another",
					},
					&pkg.Package{
						Path: "unknown",
					},
				}
			})

			JustBeforeEach(func() {
				vendored = dep.IsVendored()
			})

			Context("With all packages vendored", func() {
				BeforeEach(func() {
					for _, p := range dep.Packages {
						p.Vendored = true
					}
				})

				It("should return true", func() {
					Expect(vendored).To(BeTrue())
				})
			})

			Context("With all packages vendored except one", func() {
				BeforeEach(func() {
					for _, p := range dep.Packages {
						p.Vendored = true
					}
					dep.Packages[0].Vendored = false
				})

				It("should return false", func() {
					Expect(vendored).To(BeFalse())
				})
			})
		})

		Describe("HasPackage", func() {
			var (
				ok bool

				packagePath string
			)

			BeforeEach(func() {
				dep.Packages = []*pkg.Package{
					&pkg.Package{
						Path: "test",
					},
					&pkg.Package{
						Path: "third",
					},
					&pkg.Package{
						Path: "another",
					},
					&pkg.Package{
						Path: "unknown",
					},
				}
			})

			JustBeforeEach(func() {
				ok = dep.HasPackage(packagePath)
			})

			Context("With existing package", func() {
				BeforeEach(func() {
					packagePath = "unknown"
				})

				It("should return true", func() {
					Expect(ok).To(BeTrue())
				})
			})

			Context("With non existing package", func() {
				BeforeEach(func() {
					packagePath = "nowhere"
				})

				It("should return false", func() {
					Expect(ok).To(BeFalse())
				})
			})
		})

		Describe("IsVendorable", func() {
			var (
				ok bool
			)

			JustBeforeEach(func() {
				ok = dep.IsVendorable()
			})

			Context("With VCS fields filled", func() {
				BeforeEach(func() {
					dep.VCS = "git"
					dep.Remote = "remote"
					dep.Revision = "revision"
				})

				It("should return true", func() {
					Expect(ok).To(BeTrue())
				})
			})

			Context("Without VCS fields filled", func() {
				It("should return false", func() {
					Expect(ok).To(BeFalse())
				})
			})
		})

		Describe("IsUpdatable", func() {
			var (
				ok bool
			)

			JustBeforeEach(func() {
				ok = dep.IsUpdatable()
			})

			Context("With all VCS fields filled", func() {
				BeforeEach(func() {
					dep.VCS = "git"
					dep.Remote = "remote"
					dep.Revision = "revision"
					dep.Branch = "branch"
				})

				It("should return true", func() {
					Expect(ok).To(BeTrue())
				})
			})

			Context("Without VCS fields filled", func() {
				It("should return false", func() {
					Expect(ok).To(BeFalse())
				})
			})
		})
	})
})
