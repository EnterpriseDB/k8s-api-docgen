package main

import (
	"cnp-docgen/export"
	"cnp-docgen/parser"
	"cnp-docgen/utils"
	"flag"
	"fmt"
	"io/ioutil"
)

// Run the program giving parameters, e.g.
//  ./cnp_docgen -t json -path "../cloud-native-postgresql/api/v1/*types.go" -print true
func main() {
	var outputDoc interface{}

	// Please provide parameters "-t", "-path", "-pattern" and "-print".
	// Default are "json" for output type, current directory for path and do not print out.
	outputFormat := flag.String("t", "json", `Output format, e.g. "json","md"...`)
	path := flag.String("path", ".",
		`Path to files of interest, e.g. ../cloud-native-postgresql/api/v1. `)
	pattern := flag.String("pattern", "types.go",
		`Pattern to match the target files, e.g. "types.go". The pattern to match will be "*types.go*"`)
	out := flag.String("out", "print", `out can be "print", "file" or "none"`)

	// Parse main input parameters
	flag.Parse()

	// Get filenames in a given path, the related k8s types documentation and export it into outputDoc variable
	if filenames, err := utils.GetFilenames(*path, *pattern); err == nil {
		var kubeTypes []parser.KubeTypes
		kubeTypes, _ = parser.GetKubeTypes(filenames)
		// TODO add "md" output format
		switch *outputFormat {
		case "json":
			outputDoc, err = export.ToJSON(kubeTypes)
		default:
			// "json" is default
			outputDoc, err = export.ToJSON(kubeTypes)
		}

		// print the output documentation if main parameter "-out" is "print"
		if *out == "print" && err == nil {
			fmt.Print(outputDoc.(string))
		}
		// save the output documentation to file if main parameter "-out" is "file"
		if *out == "file" && err == nil {
			_ = ioutil.WriteFile("doc.json", []byte(outputDoc.(string)), 0644)
		}
	}

}
