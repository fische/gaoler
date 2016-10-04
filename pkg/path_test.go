package pkg_test

import (
	"go/build"
	"path/filepath"

	. "github.com/fische/gaoler/pkg"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Path", func() {
	var (
		dir string
	)

	Describe("GetPackagePath", func() {
		var (
			path string
			err  error
		)

		JustBeforeEach(func() {
			path, err = GetPackagePath(dir)
		})

		Context("With valid dir in Go SRC directories", func() {
			var res string

			BeforeEach(func() {
				res = "path/to/package"
				dir = filepath.Clean(build.Default.SrcDirs()[0] + "/" + res)
			})

			It("should not return an error", func() {
				Expect(err).To(BeNil())
			})

			It("should return corresponding package path", func() {
				Expect(path).To(Equal(res))
			})
		})

		Context("With dir outside of the Go SRC directories", func() {
			BeforeEach(func() {
				dir = "/path/to/package"
			})

			It("should return an error", func() {
				Expect(err).ToNot(BeNil())
			})
		})
	})

	Describe("IsInSrcDirs", func() {
		var (
			res bool
		)

		JustBeforeEach(func() {
			res = IsInSrcDirs(dir)
		})

		Context("With valid dir in Go SRC directories", func() {
			BeforeEach(func() {
				dir = filepath.Clean(build.Default.SrcDirs()[0] + "/test")
			})

			It("should return true", func() {
				Expect(res).To(BeTrue())
			})
		})

		Context("With dir outside of the Go SRC directories", func() {
			BeforeEach(func() {
				dir = "/path/to/package"
			})

			It("should return false", func() {
				Expect(res).To(BeFalse())
			})
		})
	})
})
