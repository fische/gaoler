package pkg_test

import (
	"fmt"
	"go/ast"

	. "github.com/fische/gaoler/pkg"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Import", func() {
	var (
		path string
	)

	Describe("GetPackagePathFromImport", func() {
		var (
			imp *ast.ImportSpec
		)

		JustBeforeEach(func() {
			path = "testpackage"
			imp = &ast.ImportSpec{
				Path: &ast.BasicLit{
					Value: fmt.Sprintf(`"%s"`, path),
				},
			}
		})

		It("should return package path without surrounding quotes", func() {
			Expect(GetPackagePathFromImport(imp)).To(Equal(path))
		})
	})

	Describe("IsPseudoPackage", func() {
		var pseudo bool

		JustBeforeEach(func() {
			pseudo = IsPseudoPackage(path)
		})

		Context("When package is a pseudo one", func() {
			BeforeEach(func() {
				path = "C"
			})

			It("should return true", func() {
				Expect(pseudo).To(BeTrue())
			})
		})

		Context("When package is not a pseudo one", func() {
			BeforeEach(func() {
				path = "normalpackage"
			})

			It("should return false", func() {
				Expect(pseudo).To(BeFalse())
			})
		})
	})
})
