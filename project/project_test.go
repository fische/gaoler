package project_test

import (
	"go/build"
	"path/filepath"

	. "github.com/fische/gaoler/project"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Project", func() {
	var (
		p *Project
	)

	Describe("New", func() {
		var (
			err error

			root string
		)

		JustBeforeEach(func() {
			p, err = New(root)
		})

		Context("When root directory is outside GO environment", func() {
			BeforeEach(func() {
				root = "/nowhere"
			})

			It("should return an error", func() {
				Expect(err).ToNot(BeNil())
			})
		})

		Context("When root directory is valid", func() {
			BeforeEach(func() {
				root = filepath.Clean(build.Default.SrcDirs()[0] + "/test")
			})

			It("should not return an error", func() {
				Expect(err).To(BeNil())
			})

			It("should fill Project fields correctly", func() {
				Expect(p.Root).To(Equal(root))
				Expect(p.Vendor).To(Equal(filepath.Clean(root + "/vendor/")))
				Expect(p.Name).To(Equal("test"))
				Expect(p.Dependencies).ToNot(BeNil())
			})
		})
	})
})
