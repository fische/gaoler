package dependency_test

import (
	"fmt"
	"os"
	"path/filepath"

	. "github.com/fische/gaoler/project/dependency"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Clean", func() {
	Describe("Checker", func() {
		var (
			opt CleanOption

			info fileInfo
		)

		BeforeEach(func() {
			info = fileInfo{}
		})

		Describe("RemoveGoTestFiles", func() {
			JustBeforeEach(func() {
				opt = RemoveGoTestFiles(info)
			})

			Context("When file is not a Go test", func() {
				BeforeEach(func() {
					info.name = "notATest.go"
				})

				It("should return Pass", func() {
					Expect(opt).To(Equal(Pass))
				})
			})

			Context("When file is a Go test file", func() {
				BeforeEach(func() {
					info.name = "file_test.go"
				})

				It("should return Remove", func() {
					Expect(opt).To(Equal(Remove))
				})
			})

			Context("When directory is testdata", func() {
				BeforeEach(func() {
					info.name = "testdata"
					info.isDir = true
				})

				It("should return Remove", func() {
					Expect(opt).To(Equal(Remove))
				})
			})
		})

		Describe("KeepGoTestFiles", func() {
			JustBeforeEach(func() {
				opt = KeepGoTestFiles(info)
			})

			Context("When file is not a Go test", func() {
				BeforeEach(func() {
					info.name = "notATest.go"
				})

				It("should return Pass", func() {
					Expect(opt).To(Equal(Pass))
				})
			})

			Context("When file is a Go test file", func() {
				BeforeEach(func() {
					info.name = "file_test.go"
				})

				It("should return Pass", func() {
					Expect(opt).To(Equal(Pass))
				})
			})

			Context("When directory is testdata", func() {
				BeforeEach(func() {
					info.name = "testdata"
					info.isDir = true
				})

				It("should return SkipDir", func() {
					Expect(opt).To(Equal(SkipDir))
				})
			})
		})
	})

	Describe("CleanVendor", func() {
		var (
			err error
			dep Dependency

			vendorRoot string
			checkers   []func(info os.FileInfo) CleanOption
		)

		BeforeEach(func() {
			dep = Dependency{}
			checkers = []func(info os.FileInfo) CleanOption{}
		})

		JustBeforeEach(func() {
			err = dep.CleanVendor(vendorRoot, checkers...)
		})

		Context("With non existing vendor directory", func() {
			BeforeEach(func() {
				vendorRoot = "nonExistingDirectory"
			})

			It("should return an error", func() {
				Expect(err).ToNot(BeNil())
			})
		})

		Context("With valid vendor directory", func() {
			BeforeEach(func() {
				vendorRoot = "testdata/vendor"
			})

			Context("With non existing package", func() {
				BeforeEach(func() {
					dep.RootPackage = "nonExistingPackage"
				})

				It("should return an error", func() {
					Expect(err).ToNot(BeNil())
				})
			})

			Context("With valid package", func() {
				var root string

				BeforeEach(func() {
					dep.RootPackage = "package"
					root = fmt.Sprintf("%s/%s", vendorRoot, dep.RootPackage)
					resetDirectory(root)
				})

				AfterEach(func() {
					removeDirectory(root)
				})

				Context("Without root package in package dependencies", func() {
					var file string

					BeforeEach(func() {
						file = fmt.Sprintf("%s/main.go", root)
						resetFile(file)
					})

					AfterEach(func() {
						removeFile(file)
					})

					It("should not return an error", func() {
						Expect(err).To(BeNil())
					})

					It("should have removed root package files without root package directory", func() {
						_, err := os.Lstat(file)
						Expect(os.IsNotExist(err)).To(BeTrue())
						_, err = os.Lstat(root)
						Expect(err).To(BeNil())
					})
				})

				Context("Without subdirectory in package dependencies", func() {
					var dir string

					BeforeEach(func() {
						dir = fmt.Sprintf("%s/subdir", root)
						resetDirectory(dir)
					})

					AfterEach(func() {
						removeDirectory(dir)
					})

					It("should not return an error", func() {
						Expect(err).To(BeNil())
					})

					It("should have removed subdirectory", func() {
						_, err := os.Lstat(dir)
						Expect(os.IsNotExist(err)).To(BeTrue())
					})
				})

				Context("With checkers", func() {
					var (
						dirToRemove  string
						dirToSkip    string
						fileToSkip   string
						fileToRemove string
					)

					BeforeEach(func() {
						dirToRemove = fmt.Sprintf("%s/toRemove", root)
						dirToSkip = fmt.Sprintf("%s/toSkip", root)
						fileToSkip = fmt.Sprintf("%s/skipped.go", dirToSkip)
						fileToRemove = fmt.Sprintf("%s/removed.go", root)

						resetDirectory(dirToRemove)
						resetDirectory(dirToSkip)
						resetFile(fileToSkip)
						resetFile(fileToRemove)

						checkers = append(checkers, func(info os.FileInfo) CleanOption {
							if filepath.Base(dirToRemove) == info.Name() || filepath.Base(fileToRemove) == info.Name() {
								return Remove
							} else if filepath.Base(dirToSkip) == info.Name() {
								return SkipDir
							} else if filepath.Base(fileToSkip) == info.Name() {
								Fail("This file should have been removed")
							}
							return Pass
						})
					})

					AfterEach(func() {
						removeDirectory(dirToRemove)
						removeDirectory(dirToSkip)
						removeFile(fileToSkip)
						removeFile(fileToRemove)
					})

					It("should not return an error", func() {
						Expect(err).To(BeNil())
					})

					It("should have removed toRemove directory", func() {
						_, err := os.Lstat(dirToRemove)
						Expect(os.IsNotExist(err)).To(BeTrue())
					})

					It("should have removed removed.go file", func() {
						_, err := os.Lstat(fileToRemove)
						Expect(os.IsNotExist(err)).To(BeTrue())
					})

					It("should not have removed toSkip directory and its content", func() {
						_, err := os.Lstat(dirToSkip)
						Expect(err).To(BeNil())
						_, err = os.Lstat(fileToSkip)
						Expect(err).To(BeNil())
					})
				})
			})
		})
	})
})
