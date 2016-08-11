package project

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

//GetImports lists imports from file at `filepath`
func GetImports(filepath string) ([]*ast.ImportSpec, error) {
	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, filepath, nil, parser.ImportsOnly)
	if err != nil {
		return nil, err
	}

	var imports []*ast.ImportSpec
	for _, i := range f.Imports {
		imports = append(imports, i)
	}
	return imports, nil
}

//GetName returns clean name of the imported package by removing surrounding `"`
func GetName(i *ast.ImportSpec) string {
	return strings.Replace(i.Path.Value, "\"", "", -1)
}
