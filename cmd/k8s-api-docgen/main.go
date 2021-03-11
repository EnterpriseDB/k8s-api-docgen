package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/EnterpriseDB/k8s-api-docgen/internal/docgen"
	"github.com/EnterpriseDB/k8s-api-docgen/internal/log"
	"github.com/EnterpriseDB/k8s-api-docgen/pkg/parser"
)

func main() {
	format := flag.String("t", "json", `Output format. The only supported one is "json"`)
	out := flag.String("o", "", "Write output to the given named file. By default "+
		"the output will be written to stdout")

	var CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	flag.Usage = func() {
		fmt.Fprintf(CommandLine.Output(), "Usage:\n  k8s-api-docgen [flags] path\n\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if len(os.Args) <= 1 {
		flag.Usage()
		return
	}

	if *format != "json" {
		fmt.Printf("Error: %v\n", docgen.ErrorWrongOutputFormat)
		flag.Usage()
		return
	}

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
