package dependency_test

import (
	"errors"
	"go/build"
	"path/filepath"

	. "github.com/fische/gaoler/project/dependency"
	"github.com/fische/gaoler/vcs"
	"github.com/fische/gaoler/vcs/modules"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type VCS struct {
	name string
	new  func(path string) (vcs.Repository, error)
	open func(path string) (vcs.Repository, error)
}

type Repository struct {
	name     string
	revision func() (string, error)
	remote   func() (string, error)
	path     func() (string, error)
	branch   func() (string, error)

	addRemote      func(remote string) error
	fetch          func() error
	checkout       func(revision string) error
	checkoutBranch func(branch string) error
}

func (v VCS) New(path string) (vcs.Repository, error) {
	return v.new(path)
}

func (v VCS) Open(path string) (vcs.Repository, error) {
	return v.open(path)
}

func (v VCS) GetName() string {
	return v.name
}

func (r Repository) GetRevision() (string, error) {
	return r.revision()
}

func (r Repository) GetRemote() (string, error) {
	return r.remote()
}

func (r Repository) GetVCSName() string {
	return r.name
}

func (r Repository) GetPath() (string, error) {
	return r.path()
}

func (r Repository) GetBranch() (string, error) {
	return r.branch()
}

func (r Repository) AddRemote(remote string) error {
	return r.addRemote(remote)
}

func (r Repository) Fetch() error {
	return r.fetch()
}

func (r Repository) Checkout(revision string) error {
	return r.checkout(revision)
}

func (r Repository) CheckoutBranch(branch string) error {
	return r.checkoutBranch(branch)
}

func getter(v string, err error) func() (string, error) {
	return func() (string, error) {
		return v, err
	}
}

var _ = Describe("Repository", func() {
	var (
		dep *Dependency
	)

	BeforeEach(func() {
		dep = &Dependency{}
	})

	Describe("OpenRepository", func() {
		var (
			err error

			dir string
		)

		JustBeforeEach(func() {
			err = dep.OpenRepository(dir)
		})

		Context("With an invalid directory", func() {
			BeforeEach(func() {
				keys := modules.VCS()
				for _, key := range keys {
					modules.Register(VCS{
						name: key,
						open: func(path string) (vcs.Repository, error) {
							return nil, errors.New("")
						},
					})
				}
				dir = "nowhere"
			})

			It("should return an error", func() {
				Expect(err).ToNot(BeNil())
			})
		})

		Context("With a valid repository", func() {
			BeforeEach(func() {
				keys := modules.VCS()
				for _, key := range keys {
					modules.Register(VCS{
						name: key,
						open: func(path string) (vcs.Repository, error) {
							return &Repository{
								name: key,
							}, nil
						},
					})
				}
				dir = "valid"
			})

			It("should not return an error", func() {
				Expect(err).To(BeNil())
			})

			It("should set Repository", func() {
				Expect(dep.Repository).ToNot(BeNil())
			})
		})
	})

	Describe("LockCurrentState", func() {
		var (
			err error

			path     func() (string, error)
			revision func() (string, error)
			remote   func() (string, error)
			branch   func() (string, error)
			name     string
		)

		BeforeEach(func() {
			path = getter(filepath.Clean(build.Default.SrcDirs()[0]+"/os"), nil)
			revision = getter("Revision", nil)
			remote = getter("Remote", nil)
			branch = getter("Branch", nil)
			name = "VCS"
		})

		JustBeforeEach(func() {
			dep.Repository = Repository{
				path:     path,
				revision: revision,
				remote:   remote,
				branch:   branch,
				name:     name,
			}
			err = dep.LockCurrentState()
		})

		Context("When getting path fails", func() {
			BeforeEach(func() {
				path = getter("", errors.New(""))
			})

			It("should return an error", func() {
				Expect(err).ToNot(BeNil())
			})
		})

		Context("With an invalid path", func() {
			BeforeEach(func() {
				path = getter("nowhere", nil)
			})

			It("should return an error", func() {
				Expect(err).ToNot(BeNil())
			})
		})

		Context("When getting revision fails", func() {
			BeforeEach(func() {
				revision = getter("", errors.New(""))
			})

			It("should return an error", func() {
				Expect(err).ToNot(BeNil())
			})
		})

		Context("When getting remote fails", func() {
			BeforeEach(func() {
				remote = getter("", errors.New(""))
			})

			It("should return an error", func() {
				Expect(err).ToNot(BeNil())
			})
		})

		Context("When getting branch fails", func() {
			BeforeEach(func() {
				branch = getter("", errors.New(""))
			})

			It("should return an error", func() {
				Expect(err).ToNot(BeNil())
			})
		})

		Context("When getting values doesn't fail", func() {
			It("should not return an error", func() {
				Expect(err).To(BeNil())
			})

			It("should set fields correctly", func() {
				Expect(dep.VCS).To(Equal("VCS"))
				Expect(dep.Branch).To(Equal("Branch"))
				Expect(dep.Remote).To(Equal("Remote"))
				Expect(dep.Revision).To(Equal("Revision"))
				Expect(dep.RootPackage).To(Equal("os"))
			})
		})
	})
})
