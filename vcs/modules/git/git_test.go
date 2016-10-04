package git_test

import (
	"log"
	"os/exec"
	"path/filepath"

	"github.com/fische/gaoler/vcs"
	. "github.com/fische/gaoler/vcs/modules/git"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Git", func() {
	var (
		g vcs.VCS
	)

	BeforeEach(func() {
		g = &Git{}
	})

	Describe("New", func() {
		var (
			r   vcs.Repository
			err error

			path string
		)

		JustBeforeEach(func() {
			r, err = g.New(path)
		})

		Context("When path is invalid", func() {
			BeforeEach(func() {
				path = "/invalid"
			})

			It("should return an error", func() {
				Expect(err).ToNot(BeNil())
			})
		})

		Context("When path is invalid", func() {
			BeforeEach(func() {
				path = "testdata/git"
				resetDirectory(path)
			})

			AfterEach(func() {
				removeDirectory(path)
			})

			It("should not return an error", func() {
				Expect(err).To(BeNil())
			})

			It("should return a valid Repository", func() {
				Expect(r).ToNot(BeNil())
				repo, ok := r.(*Repository)
				Expect(ok).To(BeTrue())
				dir, err := filepath.Abs(path)
				Expect(err).To(BeNil())
				Expect(repo.Path).To(Equal(dir))
			})
		})
	})

	Describe("Open", func() {
		var (
			r   vcs.Repository
			err error

			path string
		)

		JustBeforeEach(func() {
			r, err = g.Open(path)
		})

		Context("When path is invalid", func() {
			BeforeEach(func() {
				path = "/invalid"
			})

			It("should return an error", func() {
				Expect(err).ToNot(BeNil())
			})
		})

		Context("When path is invalid", func() {
			BeforeEach(func() {
				path = "testdata/git"
				resetDirectory(path)
				if o, err := exec.Command("git", "init", path).CombinedOutput(); err != nil {
					log.Fatalf("Could not init repository : %s", string(o))
				}
			})

			AfterEach(func() {
				removeDirectory(path)
			})

			It("should not return an error", func() {
				Expect(err).To(BeNil())
			})

			It("should return a valid Repository", func() {
				Expect(r).ToNot(BeNil())
				repo, ok := r.(*Repository)
				Expect(ok).To(BeTrue())
				dir, err := filepath.Abs(path)
				Expect(err).To(BeNil())
				Expect(repo.Path).To(Equal(dir))
			})
		})
	})

	Describe("GetName", func() {
		var (
			name string
		)

		JustBeforeEach(func() {
			name = g.GetName()
		})

		It("should return git", func() {
			Expect(name).To(Equal("git"))
		})
	})
})
