// Package parser contain the parsing code
package parser

// This code is heavily based on the po-docgen tool which can be found at:
// https://github.com/prometheus-operator/prometheus-operator/tree/master/cmd/po-docgen

import (
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
)

// GetKubeTypes return the k8s types into a slice
func GetKubeTypes(filePaths []string) ([]KubeTypes, error) {
	fset := token.NewFileSet()
	m := make(map[string]*ast.File)
	for _, filePath := range filePaths {
		f, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
		if err != nil {
			return nil, err
		}
		m[filePath] = f
	}

	// The errors raised by the creation of the Package AST are not considered as
	// we don't need to fully type-check the code and we don't have access to all
	// the types reachable by the code
	apkg, _ := ast.NewPackage(fset, m, nil, nil)

	n := doc.New(apkg, "", 0)

	var docForTypes []KubeTypes

	for _, kubType := range n.Types {
		if structType, ok := kubType.Decl.Specs[0].(*ast.TypeSpec).Type.(*ast.StructType); ok {
			var ks KubeTypes
			ks = append(ks, Pair{kubType.Name, fmtRawDoc(kubType.Doc), "", false})

			for _, field := range structType.Fields.List {
				// Treat inlined fields separately as we don't want the original types to appear in the doc.
				if isInlined(field) {
					// Skip external types, as we don't want their content to be part of the API documentation.
					if isInternalType(field.Type) {
						ks = append(ks, typesDoc[fieldType(field.Type)]...)
					}
					continue
				}

				typeString := fieldType(field.Type)
				fieldMandatory := fieldRequired(field)
				if n := fieldName(field); n != "-" {
					fieldDoc := fmtRawDoc(field.Doc.Text())
					ks = append(ks, Pair{n, fieldDoc, typeString, fieldMandatory})
				}
			}
			docForTypes = append(docForTypes, ks)
		}
	}

	return docForTypes, nil
}
