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

// Pair of strings. We need the name of fields and the doc
type Pair struct {
	Name, Doc, Type string
	Mandatory       bool
}

// KubeTypes is an array to represent all available types in a parsed file. [0] is for the type itself
type KubeTypes []Pair

const (
	docPrefix = "https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.20/"
)

var (
	links = map[string]string{
		"metav1.ObjectMeta":        docPrefix + "#objectmeta-v1-meta",
		"metav1.ListMeta":          docPrefix + "#listmeta-v1-meta",
		"metav1.LabelSelector":     docPrefix + "#labelselector-v1-meta",
		"v1.ResourceRequirements":  docPrefix + "#resourcerequirements-v1-core",
		"v1.LocalObjectReference":  docPrefix + "#localobjectreference-v1-core",
		"v1.SecretKeySelector":     docPrefix + "#secretkeyselector-v1-core",
		"v1.PersistentVolumeClaim": docPrefix + "#persistentvolumeclaim-v1-core",
		"v1.EmptyDirVolumeSource":  docPrefix + "#emptydirvolumesource-v1-core",
		"apiextensionsv1.JSON":     docPrefix + "#json-v1-apiextensions-k8s-io",
	}

	selfLinks = map[string]string{}
	typesDoc  = map[string]KubeTypes{}
)

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

func toLink(typeName string) string {
	selfLink, hasSelfLink := selfLinks[typeName]
	if hasSelfLink {
		return wrapInLink(typeName, selfLink)
	}

	link, hasLink := links[typeName]
	if hasLink {
		return wrapInLink(typeName, link)
	}

	return typeName
}

func wrapInLink(text, link string) string {
	return fmt.Sprintf("[%s](%s)", text, link)
}

func isInlined(field *ast.Field) bool {
	jsonTag := reflect.StructTag(
		field.Tag.Value[1 : len(field.Tag.Value)-1]).Get("json") // Delete first and last quotation
	return strings.Contains(jsonTag, "inline")
}

func isInternalType(typ ast.Expr) bool {
	switch it := typ.(type) {
	case *ast.SelectorExpr:
		pkg := it.X.(*ast.Ident)
		return strings.HasPrefix(pkg.Name, "monitoring")
	case *ast.StarExpr:
		return isInternalType(it.X)
	case *ast.ArrayType:
		return isInternalType(it.Elt)
	case *ast.MapType:
		return isInternalType(it.Key) && isInternalType(it.Value)
	default:
		return true
	}
}

// fieldName returns the name of the field as it should appear in JSON format
// "-" indicates that this field is not part of the JSON representation
func fieldName(field *ast.Field) string {
	// Delete first and last quotation
	jsonTag := reflect.StructTag(
		field.Tag.Value[1 : len(field.Tag.Value)-1]).Get("json")

	// This can return "-"
	jsonTag = strings.Split(jsonTag, ",")[0]
	if jsonTag == "" {
		if field.Names != nil {
			return field.Names[0].Name
		}
		return field.Type.(*ast.Ident).Name
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

func fieldType(typ ast.Expr) string {
	switch ft := typ.(type) {
	case *ast.Ident:
		return toLink(ft.Name)
	case *ast.StarExpr:
		return "*" + toLink(fieldType(ft.X))
	case *ast.SelectorExpr:
		pkg := ft.X.(*ast.Ident)
		return toLink(pkg.Name + "." + ft.Sel.Name)
	case *ast.ArrayType:
		return "[]" + toLink(fieldType(ft.Elt))
	case *ast.MapType:
		return "map[" + toLink(fieldType(ft.Key)) + "]" + toLink(fieldType(ft.Value))
	default:
		return ""
	}
}