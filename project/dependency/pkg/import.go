package pkg

import "go/ast"

var (
	pseudoPackages = []string{
		"C",
	}
)

func GetNameFromImport(imp *ast.ImportSpec) string {
	return imp.Path.Value[1 : len(imp.Path.Value)-1]
}

func IsPseudoPackage(imp *ast.ImportSpec) bool {
	n := GetNameFromImport(imp)
	for _, p := range pseudoPackages {
		if n == p {
			return true
		}
	}
	return false
}
