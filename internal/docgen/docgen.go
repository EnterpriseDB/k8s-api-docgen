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

// Package docgen contain the main functions to preprocess an AST and to write
// the documentation to the output file
package docgen

import (
	"fmt"
	"os"

	"github.com/EnterpriseDB/k8s-api-docgen/internal/log"
	"github.com/EnterpriseDB/k8s-api-docgen/pkg/export"
	"github.com/EnterpriseDB/k8s-api-docgen/pkg/parser"
)

var (
	// ErrorWrongOutputFormat means that the used specified an output format which we don't support
	ErrorWrongOutputFormat = fmt.Errorf("wrong output format")
)

// Extract extract the documentation output from the list of types given the
// output format
func Extract(kubeTypes []parser.KubeTypes, format string) (string, error) {
	switch format {
	case "json":
		return export.ToJSON(kubeTypes)

	default:
		return "", ErrorWrongOutputFormat
	}
}

// Output write the documentation to a certain file. If the filename
// is empty the documentation is written to stdout
func Output(fileName string, content string) error {
	outputStream := os.Stdout
	var err error
	if fileName != "" {
		outputStream, err = os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600) // #nosec
		if err != nil {
			return err
		}

		defer func() {
			err = outputStream.Close()
			if err != nil {
				log.Log.Error(err, "Cannot close output file",
					"fileName", fileName)
			}
		}()
	}

	_, err = outputStream.Write([]byte(content))

	return err
}
