package project

import (
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"strings"
)

//Import represents a Go package import
type Import struct {
	*ast.ImportSpec
	*build.Package
}

//NewImport creates a new `Import` and retrieves `build.Package` from `i`
func NewImport(i *ast.ImportSpec) (*Import, error) {
	p, err := build.Import(strings.Replace(i.Path.Value, "\"", "", -1), "", build.FindOnly)
	if err != nil {
		return nil, err
	}
	return &Import{
		ImportSpec: i,
		Package:    p,
	}, nil
}

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
