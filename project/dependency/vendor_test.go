package dependency_test

import (
	"errors"
	"fmt"
	"path/filepath"

	. "github.com/fische/gaoler/project/dependency"
	"github.com/fische/gaoler/vcs"
	"github.com/fische/gaoler/vcs/modules"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Vendor", func() {
	var (
		err error

		vendorRoot string
		dep        *Dependency
	)

	BeforeEach(func() {
		dep = &Dependency{}
	})

	Describe("Vendor", func() {
		JustBeforeEach(func() {
			err = dep.Vendor(vendorRoot)
		})

		AfterEach(func() {
			removeDirectory(vendorRoot)
		})

		Context("With an unknown vendor path", func() {
			BeforeEach(func() {
				vendorRoot = "nowhere/unkown"
			})

			It("should return an error", func() {
				Expect(err).ToNot(BeNil())
			})
		})

		Context("With an unknown VCS", func() {
			BeforeEach(func() {
				vendorRoot = "testdata/vendor"
				resetDirectory(vendorRoot)

				dep.VCS = "___unknown___"
			})

			It("should return an error", func() {
				Expect(err).ToNot(BeNil())
			})
		})

		Context("When opening VCS repository fails", func() {
			BeforeEach(func() {
				vendorRoot = "testdata/vendor"
				resetDirectory(vendorRoot)

				v := "VCS"
				dep.VCS = v
				modules.Register(VCS{
					name: v,
					new: func(path string) (vcs.Repository, error) {
						return nil, errors.New("")
					},
				})
			})

			It("should return an error", func() {
				Expect(err).ToNot(BeNil())
			})
		})

		Context("When opening VCS reposity succeeds", func() {
			var (
				repo *Repository
			)

			BeforeEach(func() {
				vendorRoot = "testdata/vendor"
				resetDirectory(vendorRoot)
				dep.VCS = "VCS"
				dep.Remote = "Remote"
				dep.Revision = "Revision"
				dep.RootPackage = "RootPackage"

				repo = &Repository{
					fetch: func() error {
						return nil
					},
					addRemote: func(remote string) error {
						Expect(remote).To(Equal(dep.Remote))
						return nil
					},
					checkout: func(revision string) error {
						Expect(revision).To(Equal(dep.Revision))
						return nil
					},
				}
				modules.Register(VCS{
					name: dep.VCS,
					new: func(path string) (vcs.Repository, error) {
						Expect(path).To(Equal(filepath.Clean(fmt.Sprintf("%s/%s/", vendorRoot, dep.RootPackage))))
						return repo, nil
					},
				})
			})

			It("should not return an error", func() {
				Expect(err).To(BeNil())
			})

			It("should set dependency repository", func() {
				Expect(dep.Repository).To(Equal(repo))
			})
		})
	})

	Describe("Update", func() {
		var (
			updated bool
		)

		JustBeforeEach(func() {
			updated, err = dep.Update(vendorRoot)
		})

		Context("When vendoring fails", func() {
			BeforeEach(func() {
				vendorRoot = "nowhere/unkown"
			})

			It("should return an error", func() {
				Expect(err).ToNot(BeNil())
			})
		})

		Context("When checking out branch fails", func() {
			BeforeEach(func() {
				vendorRoot = "testdata/vendor"
				resetDirectory(vendorRoot)
				dep.VCS = "VCS"
				dep.Remote = "Remote"
				dep.Revision = "Revision"
				dep.RootPackage = "RootPackage"
				dep.Branch = "Branch"

				modules.Register(VCS{
					name: dep.VCS,
					new: func(path string) (vcs.Repository, error) {
						Expect(path).To(Equal(filepath.Clean(fmt.Sprintf("%s/%s/", vendorRoot, dep.RootPackage))))
						return &Repository{
							fetch: func() error {
								return nil
							},
							addRemote: func(remote string) error {
								return nil
							},
							checkout: func(revision string) error {
								return nil
							},
							checkoutBranch: func(branch string) error {
								Expect(branch).To(Equal(dep.Branch))
								return errors.New("")
							},
						}, nil
					},
				})
			})

			It("should return an error", func() {
				Expect(err).ToNot(BeNil())
			})
		})

		Context("When getting revision fails", func() {
			BeforeEach(func() {
				vendorRoot = "testdata/vendor"
				resetDirectory(vendorRoot)
				dep.VCS = "VCS"
				dep.Remote = "Remote"
				dep.Revision = "Revision"
				dep.RootPackage = "RootPackage"
				dep.Branch = "Branch"

				modules.Register(VCS{
					name: dep.VCS,
					new: func(path string) (vcs.Repository, error) {
						Expect(path).To(Equal(filepath.Clean(fmt.Sprintf("%s/%s/", vendorRoot, dep.RootPackage))))
						return &Repository{
							fetch: func() error {
								return nil
							},
							addRemote: func(remote string) error {
								return nil
							},
							checkout: func(revision string) error {
								return nil
							},
							checkoutBranch: func(branch string) error {
								Expect(branch).To(Equal(dep.Branch))
								return nil
							},
							revision: func() (string, error) {
								return "", errors.New("")
							},
						}, nil
					},
				})
			})

			It("should return an error", func() {
				Expect(err).ToNot(BeNil())
			})
		})

		Context("When revision is the same", func() {
			BeforeEach(func() {
				vendorRoot = "testdata/vendor"
				resetDirectory(vendorRoot)
				dep.VCS = "VCS"
				dep.Remote = "Remote"
				dep.Revision = "Revision"
				dep.RootPackage = "RootPackage"
				dep.Branch = "Branch"

				modules.Register(VCS{
					name: dep.VCS,
					new: func(path string) (vcs.Repository, error) {
						Expect(path).To(Equal(filepath.Clean(fmt.Sprintf("%s/%s/", vendorRoot, dep.RootPackage))))
						return &Repository{
							fetch: func() error {
								return nil
							},
							addRemote: func(remote string) error {
								return nil
							},
							checkout: func(revision string) error {
								return nil
							},
							checkoutBranch: func(branch string) error {
								Expect(branch).To(Equal(dep.Branch))
								return nil
							},
							revision: func() (string, error) {
								return dep.Revision, nil
							},
						}, nil
					},
				})
			})

			It("should not return an error", func() {
				Expect(err).To(BeNil())
			})

			It("should not have updated revision", func() {
				Expect(updated).To(BeFalse())
				Expect(dep.Revision).To(Equal("Revision"))
			})
		})

		Context("When revision is not the same", func() {
			BeforeEach(func() {
				vendorRoot = "testdata/vendor"
				resetDirectory(vendorRoot)
				dep.VCS = "VCS"
				dep.Remote = "Remote"
				dep.Revision = "Revision"
				dep.RootPackage = "RootPackage"
				dep.Branch = "Branch"

				modules.Register(VCS{
					name: dep.VCS,
					new: func(path string) (vcs.Repository, error) {
						Expect(path).To(Equal(filepath.Clean(fmt.Sprintf("%s/%s/", vendorRoot, dep.RootPackage))))
						return &Repository{
							fetch: func() error {
								return nil
							},
							addRemote: func(remote string) error {
								return nil
							},
							checkout: func(revision string) error {
								return nil
							},
							checkoutBranch: func(branch string) error {
								Expect(branch).To(Equal(dep.Branch))
								return nil
							},
							revision: func() (string, error) {
								return dep.Revision + dep.Revision, nil
							},
						}, nil
					},
				})
			})

			It("should not return an error", func() {
				Expect(err).To(BeNil())
			})

			It("should not have updated revision", func() {
				Expect(updated).To(BeTrue())
				Expect(dep.Revision).To(Equal("RevisionRevision"))
			})
		})
	})
})
