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

// Package json contain the code exporting the internal data to JSON
package json

import (
	"encoding/json"

	"github.com/EnterpriseDB/k8s-api-docgen/pkg/parser"
)

// k8s types for generation of docs
type kubeType struct {
	Name  string     `json:"name"`
	Doc   string     `json:"description"`
	Items []kubeItem `json:"items"`
}

// k8s items
type kubeItem struct {
	Name      string `json:"field"`
	Doc       string `json:"description"`
	Type      string `json:"schema"`
	Mandatory bool   `json:"required"`
}

func convertToKubeTypes(kt parser.KubeTypes) []kubeType {
	kubeDocs := make([]kubeType, len(kt))
	for idx, kubeStructure := range kt {
		k := kubeType{
			Name:  kubeStructure.Name,
			Doc:   kubeStructure.Doc,
			Items: nil,
		}

		for _, item := range kubeStructure.Fields {
			k.Items = append(k.Items, kubeItem{
				Name:      item.Name,
				Doc:       item.Doc,
				Type:      item.Type.Name,
				Mandatory: item.Mandatory,
			})
		}
		kubeDocs[idx] = k
	}
	return kubeDocs
}

// ToJSON get a slice of KubeTypes as input and return the JSON documentation.
// JSON fields are the ones defined in kubeTypes (and kubeItem) definition.
func ToJSON(kt parser.KubeTypes) (string, error) {
	kubeDocs := convertToKubeTypes(kt)

	j, err := json.MarshalIndent(kubeDocs, "", "\t")
	return string(j), err
}
