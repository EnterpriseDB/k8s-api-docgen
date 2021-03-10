package main

import (
	"flag"

	"github.com/EnterpriseDB/k8s-api-docgen/internal/docgen"
	"github.com/EnterpriseDB/k8s-api-docgen/internal/log"
	"github.com/EnterpriseDB/k8s-api-docgen/pkg/parser"
)

func main() {
	format := flag.String("t", "json", `Output format. The only supported one is "json"`)
	out := flag.String("o", "", "Write output to the given named file. By default "+
		"the output will be written to stdout")

	flag.Parse()

	var kubeTypes []parser.KubeTypes
	kubeTypes, err := parser.GetKubeTypes(flag.Args())
	if err != nil {
		log.Log.Error(
			err, "Error while parsing source files",
			"args", flag.Args())
		return
	}

	output, err := docgen.Extract(kubeTypes, *format)
	if err != nil {
		log.Log.Error(err, "Error while exporting data")
		return
	}

	if err = docgen.Output(*out, output); err != nil {
		log.Log.Error(err, "Cannot write JSON output")
	}
}
