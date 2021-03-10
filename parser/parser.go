package parser

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"reflect"
	"strings"
)

// Return the k8s types into a slice
func GetKubeTypes(filePaths []string) ([]KubeTypes, error) {
	fset := token.NewFileSet()
	m := make(map[string]*ast.File)
	for _, filePath := range filePaths {
		f, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
		if err != nil {
			fmt.Println(err)
		}
		m[filePath] = f
	}
	apkg, err := ast.NewPackage(fset, m, nil, nil)
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
	return docForTypes, err
}

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

var (
	links = map[string]string{
		"metav1.ObjectMeta":        "https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#objectmeta-v1-meta",
		"metav1.ListMeta":          "https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#listmeta-v1-meta",
		"metav1.LabelSelector":     "https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#labelselector-v1-meta",
		"v1.ResourceRequirements":  "https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#resourcerequirements-v1-core",
		"v1.LocalObjectReference":  "https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#localobjectreference-v1-core",
		"v1.SecretKeySelector":     "https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#secretkeyselector-v1-core",
		"v1.PersistentVolumeClaim": "https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#persistentvolumeclaim-v1-core",
		"v1.EmptyDirVolumeSource":  "https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#emptydirvolumesource-v1-core",
		"apiextensionsv1.JSON":     "https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#json-v1-apiextensions-k8s-io",
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
		case len(line) == 0: // Keep paragraphs
			delPrevChar()
			buffer.WriteString("\n\n")
		case strings.HasPrefix(leading, "TODO"): // Ignore one line TODOs
		case strings.HasPrefix(leading, "+"): // Ignore instructions to go2idl
		default:
			if strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t") {
				delPrevChar()
				line = "\n" + line + "\n" // Replace it with newline. This is useful when we have a line with: "Example:\n\tJSON-someting..."
			} else {
				line += " "
			}
			buffer.WriteString(line)
		}
	}

	postDoc := strings.TrimRight(buffer.String(), "\n")
	postDoc = strings.Replace(postDoc, "\\\"", "\"", -1) // replace user's \" to "
	postDoc = strings.Replace(postDoc, "\"", "\\\"", -1) // Escape "
	postDoc = strings.Replace(postDoc, "\n", "\\n", -1)
	postDoc = strings.Replace(postDoc, "\t", "\\t", -1)
	postDoc = strings.Replace(postDoc, "|", "\\|", -1)

	return postDoc
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
	jsonTag := reflect.StructTag(field.Tag.Value[1 : len(field.Tag.Value)-1]).Get("json") // Delete first and last quotation
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
	jsonTag := reflect.StructTag(field.Tag.Value[1 : len(field.Tag.Value)-1]).Get("json") // Delete first and last quotation
	jsonTag = strings.Split(jsonTag, ",")[0]                                              // This can return "-"
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
		jsonTag = reflect.StructTag(field.Tag.Value[1 : len(field.Tag.Value)-1]).Get("json") // Delete first and last quotation
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
