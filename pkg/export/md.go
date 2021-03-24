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

// Package export contain the code exporting the internal data to the user
package export

import (
	"bytes"
	"io/ioutil"
	"strings"
	"text/template"

	"github.com/EnterpriseDB/k8s-api-docgen/pkg/parser"
)

const (
	tableFieldName    string = "Name"
	tableFieldDoc     string = "Doc"
	tableFieldRawType string = "Type"
)

// ToMd get a slice of KubeTypes as input and return the Markdown documentation.
func ToMd(kt []parser.KubeTypes) (string, error) {
	kubeDocs := convertToKubeTypes(kt)
	format(&kubeDocs)

	mdTemplate, err := ioutil.ReadFile("md-template.md")
	if err != nil {
		return "", err
	}
	md, err := runTemplate(mdTemplate, kubeDocs)
	if err != nil {
		return "", err
	}
	return md, err
}

// format applies proper formats to tables and documentation
func format(kubeDocs *[]kubeTypes) {
	for i, k := range *kubeDocs {
		(*kubeDocs)[i].Doc = strings.Trim(k.Doc, "\n")
		nameMaxLength := int(k.maxSizeOfTableFields[indexName])
		docMaxLength := int(k.maxSizeOfTableFields[indexDoc])
		rawTypeMaxLength := int(k.maxSizeOfTableFields[indexRawType])
		(*kubeDocs)[i].TableFieldName = rightPad(tableFieldName, nameMaxLength-len(tableFieldName))
		(*kubeDocs)[i].TableFieldNameDashSize = strings.Repeat("-", nameMaxLength)
		(*kubeDocs)[i].TableFieldDoc = rightPad(tableFieldDoc, docMaxLength-len(tableFieldDoc))
		(*kubeDocs)[i].TableFieldDocDashSize = strings.Repeat("-", docMaxLength)
		(*kubeDocs)[i].TableFieldRawType = rightPad(tableFieldRawType, rawTypeMaxLength-len(tableFieldRawType))
		(*kubeDocs)[i].TableFieldRawTypeDashSize = strings.Repeat("-", rawTypeMaxLength)
		for j, item := range k.Items {
			(*kubeDocs)[i].Items[j].Name = rightPad(item.Name, nameMaxLength-len(item.Name))
			(*kubeDocs)[i].Items[j].Doc = rightPad(item.Doc, docMaxLength-len(item.Doc))
			// adding hyperlinks to documented keys
			(*kubeDocs)[i].Items[j].RawType = rightPad((*kubeDocs)[i].Items[j].RawType, rawTypeMaxLength-len(item.RawType))
		}
	}
}

// runTemplate execute the template, fed by docs values
func runTemplate(aTemplate []byte, docs []kubeTypes) (string, error) {
	var w bytes.Buffer
	tmpl, err := template.New("KubeTypes").Parse(string(aTemplate))
	if err != nil {
		return "", err
	}
	_, err = w.WriteString("# " + k8sAPIDocgen + "\n")
	if err != nil {
		return "", err
	}
	err = tmpl.Execute(&w, docs)
	if err != nil {
		return "", err
	}
	return w.String(), nil
}
