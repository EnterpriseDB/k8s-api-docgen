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
	if fileName != "" {
		outputStream, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600) // #nosec
		if err != nil {
			return err
		}

		defer func() {
			err := outputStream.Close()
			if err != nil {
				log.Log.Error(err, "Cannot close output file",
					"fileName", fileName)
			}
		}()
	}

	_, err := outputStream.Write([]byte(content))
	return err
}
