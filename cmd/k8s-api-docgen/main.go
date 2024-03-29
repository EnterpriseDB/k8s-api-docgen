/*
Copyright 2021 EnterpriseDB Corporation

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
// Package main is the entrypoint for k8s-api-docgen
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
	format := flag.String("t", string(docgen.OutputTypeJSON),
		`Output format. The only supported ones are "json" (JSON) and "md" (Markdown)`)
	out := flag.String("o", "", "Write output to the given named file. By default "+
		"the output will be written to stdout")
	mdConfiguration := flag.String("c", "md-configuration.yaml",
		"Path of the YAML file containing Markdown configuration. By default the "+
			"configuration will be read from 'md-configuration.yaml'")
	mdTemplate := flag.String("m", "md-template.md",
		"Path of the Markdown template file for generating Markdown documentation. By default the "+
			"Markdown template will be read from 'md-template.md'")

	CommandLine := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	flag.Usage = func() {
		_, _ = fmt.Fprintf(CommandLine.Output(), "Usage:\n  k8s-api-docgen [flags] path\n\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if len(os.Args) <= 1 {
		flag.Usage()
		return
	}

	if *format != string(docgen.OutputTypeJSON) && *format != string(docgen.OutputTypeMD) {
		fmt.Printf("Error: %v\n", docgen.ErrorWrongOutputFormat)
		flag.Usage()
		return
	}

	var kubeTypes parser.KubeTypes
	kubeTypes, err := parser.GetKubeTypes(flag.Args())
	if err != nil {
		log.Log.Error(
			err, "Error while parsing source files",
			"args", flag.Args())
		return
	}

	output, err := docgen.Extract(kubeTypes, docgen.OutputType(*format), *mdConfiguration, *mdTemplate)
	if err != nil {
		log.Log.Error(err, "Error while exporting data")
		return
	}

	if err = docgen.Output(*out, output); err != nil {
		log.Log.Error(err, "Cannot write output file")
	}
}
