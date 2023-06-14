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

package parser

import (
	"bytes"
	"fmt"
	"go/ast"
	"reflect"
	"strings"
)

// -------------------------------------------------------------------------------------------------------------------
// Auxiliary types and functions. Parsing system reuse code from
// https://github.com/prometheus-operator/prometheus-operator/tree/master/cmd/po-docgen
// ------------------------------------------------------------------------------------------------------------------

// KubeField is a struct with all the types we need to generate docs
type KubeField struct {
	// The field name (in the JSON representation of this object)
	Name string

	// The field type
	Type TypeInfo

	// The normalized documentation
	Doc string

	// Mandatory flag
	Mandatory bool
}

// TypeInfo is a struct representing a type with a given name and it's base type name.
// I.e. a type named `[]Pod` has `Pod` as a base type. Atomic types have `Name == BaseName`.
//
// We are adopting a simplification here: we consider only type constructors with 1 parameter.
// The only multiple-arity type constructor we have is `map[T1]T2` and, since kubernetes
// resources must be JSON-serializable, T1 == string. Given that, T1 is not interesting and
// we are only using T2 as base type.
type TypeInfo struct {
	// The type name (i.e. `[]Pod`)
	Name string

	// The base type name (i.e. `Pod`)
	BaseType string

	// The type-constructor who generated the type (i.e. `[]`)
	Constructor string

	// True if the type is internal to this package and false otherwise
	Internal bool
}

// KubeStructure represent a structure that we need to document
type KubeStructure struct {
	// The structure name
	Name string

	// The normalized documentation
	Doc string

	// The structure fields
	Fields []KubeField
}

// KubeTypes is an array to represent all available types in a parsed file. [0] is for the type itself
type KubeTypes []KubeStructure

func fmtRawDoc(rawDoc string) string {
	var buffer bytes.Buffer
	delPrevChar := func() {
		if buffer.Len() > 0 {
			buffer.Truncate(buffer.Len() - 1) // Delete the last " " or "\n"
		}
	}

	// Ignore all lines after ---
	rawDoc = strings.Split(rawDoc, "---")[0]

	for _, line := range strings.Split(rawDoc, "\n") {
		line = strings.TrimRight(line, " ")
		leading := strings.TrimLeft(line, " ")
		switch {
		// Keep paragraphs
		case len(line) == 0:
			delPrevChar()
			buffer.WriteString("\n\n")

			// Ignore one line TODOs
		case strings.HasPrefix(leading, "TODO"):
			// Ignore instructions to go2idl
		case strings.HasPrefix(leading, "+"):
		default:
			if strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t") {
				delPrevChar()

				// Replace it with newline. This is useful when we have a line with: "Example:\n\tJSON-someting..."
				line = fmt.Sprintf("\n%v\n", line)
			} else {
				line += " "
			}
			buffer.WriteString(line)
		}
	}

	return strings.TrimRight(buffer.String(), "\n")
}

func isInlined(field *ast.Field) bool {
	if field.Tag != nil {
		jsonTag := reflect.StructTag(
			field.Tag.Value[1 : len(field.Tag.Value)-1]).Get("json") // Delete first and last quotation
		return strings.Contains(jsonTag, "inline")
	}
	return false
}

// fieldName returns the name of the field as it should appear in JSON format
// "-" indicates that this field is not part of the JSON representation
func fieldName(field *ast.Field) string {
	// If json tag does not exists, we use the name instead
	jsonTag := ""
	if field.Tag != nil {
		// Delete first and last quotation
		jsonTag = reflect.StructTag(
			field.Tag.Value[1 : len(field.Tag.Value)-1]).Get("json")
	}

	// This can return "-"
	jsonTag = strings.Split(jsonTag, ",")[0]
	if jsonTag == "" {
		if field.Names != nil {
			return field.Names[0].Name
		}

		// when there is an inline field, the type is different from '*ast.Ident'(actually it is '*ast.SelectorExpr')
		// we can return empty json tag if not '*ast.Ident', as default, when switching over the field type
		switch t := (field.Type).(type) {
		case *ast.Ident:
			return t.Name
		default:
			return ""
		}
	}
	return jsonTag
}

// fieldRequired returns whether a field is a required field.
func fieldRequired(field *ast.Field) bool {
	jsonTag := ""
	if field.Tag != nil {
		jsonTag = reflect.StructTag(
			field.Tag.Value[1 : len(field.Tag.Value)-1]).Get("json") // Delete first and last quotation
		return !strings.Contains(jsonTag, "omitempty")
	}

	return false
}

func fieldType(typ ast.Expr) TypeInfo {
	switch ft := typ.(type) {
	case *ast.Ident:
		return TypeInfo{
			Name:        ft.Name,
			BaseType:    ft.Name,
			Constructor: "",
			Internal:    true,
		}
	case *ast.StarExpr:
		return TypeInfo{
			Name:        "*" + fieldType(ft.X).Name,
			BaseType:    fieldType(ft.X).Name,
			Constructor: "*",
			Internal:    fieldType(ft.X).Internal,
		}
	case *ast.SelectorExpr:
		pkg := ft.X.(*ast.Ident)
		return TypeInfo{
			Name:        pkg.Name + "." + ft.Sel.Name,
			BaseType:    pkg.Name + "." + ft.Sel.Name,
			Constructor: "",
			Internal:    false,
		}
	case *ast.ArrayType:
		return TypeInfo{
			Name:        "[]" + fieldType(ft.Elt).Name,
			BaseType:    fieldType(ft.Elt).Name,
			Constructor: "[]",
			Internal:    fieldType(ft.Elt).Internal,
		}
	case *ast.MapType:
		return TypeInfo{
			Name:        fmt.Sprintf("map[%v]%v", fieldType(ft.Key).Name, fieldType(ft.Value).Name),
			BaseType:    fieldType(ft.Value).Name,
			Constructor: "map[]",
			Internal:    fieldType(ft.Value).Internal,
		}
	default:
		return TypeInfo{}
	}
}
