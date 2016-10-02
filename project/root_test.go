package project_test

import (
	"fmt"
	"go/build"
	"log"
	"path/filepath"

	. "github.com/fische/gaoler/project"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Root", func() {
	Describe("GetProjectRootFromDir", func() {
		var (
			path string
			err  error

			dir string
		)

		JustBeforeEach(func() {
			path, err = GetProjectRootFromDir(dir)
		})

		Context("When dir is outside GO environment", func() {
			BeforeEach(func() {
				dir = "/nowhere"
			})

			It("should return an error", func() {
				Expect(err).ToNot(BeNil())
			})
		})

		Context("When dir is in the GO environment", func() {
			Context("When dir is not a subdirectory of a main package", func() {
				BeforeEach(func() {
					dir = filepath.Clean(build.Default.SrcDirs()[0] + "/")
				})

				It("should return an error", func() {
					Expect(err).ToNot(BeNil())
				})
			})

			Context("When dir is a valid subdirectory of a main package", func() {
				BeforeEach(func() {
					dir, err = filepath.Abs("testdata/subdir/subpackage")
					if err != nil {
						log.Fatalf("Could not get absolute path : %v", err)
					}
				})

				It("should not return an error", func() {
					Expect(err).To(BeNil())
				})

				It("should return a valid path", func() {
					abs, err := filepath.Abs("testdata")
					if err != nil {
						Fail(fmt.Sprintf("Could not get absolute path : %v", err))
					}
					Expect(path).To(Equal(abs))
				})
			})

			Context("When parsing go files fails", func() {
				BeforeEach(func() {
					dir, err = filepath.Abs("testdata/invalid")
					if err != nil {
						log.Fatalf("Could not get absolute path : %v", err)
					}
				})

				It("should not return an error", func() {
					Expect(err).To(BeNil())
				})

				It("should skip directory", func() {
					abs, err := filepath.Abs("testdata")
					if err != nil {
						Fail(fmt.Sprintf("Could not get absolute path : %v", err))
					}
					Expect(path).To(Equal(abs))
				})
			})
		})
	})
})
