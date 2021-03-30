/*
Copyright 2021 EnterpriseDB Corporation
Some portions Copyright 2016 The prometheus-operator Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// This code is heavily based on the po-docgen tool which can be found at:
// https://github.com/prometheus-operator/prometheus-operator/tree/master/cmd/po-docgen

// Package parser contain the parsing code
package parser

import (
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
)

// GetKubeTypes return the k8s types into a slice
func GetKubeTypes(filePaths []string) (KubeTypes, error) {
	// Parse the input files or exit with an error state
	fSet := token.NewFileSet()
	m := make(map[string]*ast.File)
	for _, filePath := range filePaths {
		f, err := parser.ParseFile(fSet, filePath, nil, parser.ParseComments)
		if err != nil {
			return nil, err
		}
		m[filePath] = f
	}

	// The errors raised by the creation of the Package AST are not considered as
	// we don't need to fully type-check the code and we don't have access to all
	// the types reachable by the code
	apkg, _ := ast.NewPackage(fSet, m, nil, nil)

	n := doc.New(apkg, "", 0)

	var docForTypes KubeTypes

	for _, kubType := range n.Types {
		if structType, ok := kubType.Decl.Specs[0].(*ast.TypeSpec).Type.(*ast.StructType); ok {
			kubeStructure := KubeStructure{
				Name: kubType.Name,
				Doc:  fmtRawDoc(kubType.Doc),
			}

			for _, field := range structType.Fields.List {
				// TODO: manage inlined fields
				if isInlined(field) {
					// TODO: Ask Gabriele why typesDoc was not used
					continue
				}

				typeInfo := fieldType(field.Type)
				fieldMandatory := fieldRequired(field)
				if n := fieldName(field); n != "-" {
					fieldDoc := fmtRawDoc(field.Doc.Text())
					kubeStructure.Fields = append(kubeStructure.Fields,
						KubeField{
							Name:      n,
							Type:      typeInfo,
							Doc:       fieldDoc,
							Mandatory: fieldMandatory,
						})
				}
			}
			docForTypes = append(docForTypes, kubeStructure)
		}
	}
	return docForTypes, nil
}
