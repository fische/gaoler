package pkg_test

import (
	"fmt"
	"go/ast"
	"os"
	"path/filepath"

	. "github.com/fische/gaoler/pkg"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Package", func() {
	var (
		path string
		p    *Package
	)

	Describe("New", func() {
		BeforeEach(func() {
			path = "testpath"
		})

		JustBeforeEach(func() {
			p = New(path)
		})

		It("should create a new instance of Package with the given path", func() {
			Expect(p.Path).To(Equal(path))
		})
	})

	Describe("NewFromImport", func() {
		var imp *ast.ImportSpec

		BeforeEach(func() {
			path = "testpath"
			imp = &ast.ImportSpec{
				Path: &ast.BasicLit{
					Value: fmt.Sprintf(`"%s"`, path),
				},
			}
		})

		JustBeforeEach(func() {
			p = NewFromImport(imp)
		})

		It("should create a new instance of Package with the path from the import", func() {
			Expect(p.Path).To(Equal(path))
		})
	})

	Describe("Methods", func() {
		JustBeforeEach(func() {
			p = &Package{
				Path: path,
			}
		})

		Describe("Import", func() {
			var (
				ignoreVendor bool
				srcPath      string

				err error
			)

			JustBeforeEach(func() {
				err = p.Import(srcPath, ignoreVendor)
			})

			Context("With valid package", func() {
				BeforeEach(func() {
					srcPath = ""
					path = "os"
					ignoreVendor = true
				})

				It("should not return an error", func() {
					Expect(err).To(BeNil())
				})

				It("should set Root to true", func() {
					Expect(p.Root).To(BeTrue())
				})
			})
			Context("With vendored package", func() {
				BeforeEach(func() {
					wd, _ := os.Getwd()
					srcPath = filepath.Clean(wd + "/testdata")
					path = "package"
					ignoreVendor = false
				})

				It("should not return an error", func() {
					Expect(err).To(BeNil())
				})

				It("should set Vendored to true", func() {
					Expect(p.Vendored).To(BeTrue())
				})
			})
			Context("With non existing package", func() {
				BeforeEach(func() {
					srcPath = ""
					path = "__NON_EXISTING__"
					ignoreVendor = true
				})

				It("should return error", func() {
					Expect(err).ToNot(BeNil())
				})
			})
		})

		Describe("IsPseudoPackage", func() {
			var pseudo bool

			JustBeforeEach(func() {
				pseudo = p.IsPseudoPackage()
			})

			Context("With a pseudo package", func() {
				BeforeEach(func() {
					path = "C"
				})

				It("should return true", func() {
					Expect(pseudo).To(BeTrue())
				})
			})

			Context("Without a pseudo package", func() {
				BeforeEach(func() {
					path = "noPseudo"
				})

				It("should return false", func() {
					Expect(pseudo).To(BeFalse())
				})
			})
		})
	})
})
