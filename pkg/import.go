package pkg

import "go/ast"

var (
	pseudoPackages = []string{
		"C",
	}
)

func GetPackagePathFromImport(imp *ast.ImportSpec) string {
	return imp.Path.Value[1 : len(imp.Path.Value)-1]
}

func IsPseudoPackage(pkgPath string) bool {
	for _, p := range pseudoPackages {
		if pkgPath == p {
			return true
		}
	}
	return false
}
