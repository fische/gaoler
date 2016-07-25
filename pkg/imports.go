package pkg

import (
	"go/ast"
	"go/parser"
	"go/token"
)

//ListImports lists imports from file at `filepath`
func ListImports(filepath string) ([]*ast.ImportSpec, error) {
	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, filepath, nil, parser.ImportsOnly)
	if err != nil {
		return nil, err
	}

	var imports []*ast.ImportSpec
	for _, s := range f.Imports {
		imports = append(imports, s)
	}
	return imports, nil
}
