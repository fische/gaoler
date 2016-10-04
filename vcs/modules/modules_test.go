package modules_test

import (
	"errors"

	"github.com/fische/gaoler/vcs"
	. "github.com/fische/gaoler/vcs/modules"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Modules", func() {
	Describe("Register / GetVCS / VCS", func() {
		var (
			v          *VCSTest
			ret        vcs.VCS
			registered []string
		)

		JustBeforeEach(func() {
			Register(v)
			ret = GetVCS(v.name)
			registered = VCS()
		})

		Context("With a new VCS", func() {
			BeforeEach(func() {
				v = &VCSTest{
					name: "test",
				}
			})

			It("should have correctly registered VCS", func() {
				Expect(ret).To(Equal(v))
				passed := false
				for _, reg := range registered {
					if reg == v.name {
						passed = true
					}
				}
				Expect(passed).To(BeTrue())
			})
		})

		Context("With an existing VCS", func() {
			BeforeEach(func() {
				v = &VCSTest{
					name: "git",
				}
			})

			It("should have correctly overridden VCS", func() {
				Expect(ret).To(Equal(v))
				passed := false
				for _, reg := range registered {
					if reg == v.name {
						passed = true
					}
				}
				Expect(passed).To(BeTrue())
			})
		})
	})

	Describe("OpenRepository", func() {
		var (
			open func(path string) (vcs.Repository, error)

			r   vcs.Repository
			err error

			path string
		)

		JustBeforeEach(func() {
			for _, v := range VCS() {
				Register(VCSTest{
					name: v,
					open: open,
				})
			}
			r, err = OpenRepository(path)
		})

		Context("With an invalid path", func() {
			BeforeEach(func() {
				open = func(path string) (vcs.Repository, error) {
					return nil, errors.New("")
				}
				path = "/invalid"
			})

			It("should return an error", func() {
				Expect(err).ToNot(BeNil())
			})
		})

		Context("With an invalid path", func() {
			var (
				repo *Repository
			)

			BeforeEach(func() {
				repo = &Repository{
					name: "test",
				}
				open = func(path string) (vcs.Repository, error) {
					return repo, nil
				}
				path = "/valid"
			})

			It("should not return an error", func() {
				Expect(err).To(BeNil())
			})

			It("should have returned a valid repository", func() {
				Expect(r).To(Equal(repo))
			})
		})
	})
})
