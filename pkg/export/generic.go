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
	"math"
	"strings"

	"github.com/EnterpriseDB/k8s-api-docgen/pkg/parser"
)

const (
	k8sAPIDocgen string = "K8s Api Docgen"
)

// k8s types for generation of docs
type kubeTypes struct {
	Name                      string `json:"name"`
	NameWithAnchor            string `json:"-"`
	Doc                       string     `json:"description"`
	Items                     []kubeItem `json:"items"`
	maxSizeOfTableFields      [4]float64
	TableFieldName            string `json:"-"`
	TableFieldNameDashSize    string `json:"-"`
	TableFieldDoc             string `json:"-"`
	TableFieldDocDashSize     string `json:"-"`
	TableFieldRawType         string `json:"-"`
	TableFieldRawTypeDashSize string `json:"-"`
}

// k8s items
type kubeItem struct {
	Name      string `json:"field"`
	Doc       string `json:"description"`
	Type      string `json:"schema"`
	RawType   string `json:"-"`
	Mandatory bool `json:"required"`
}

const (
	indexName int = iota
	indexDoc
	indexRawType
	indexMandatory
)
const (
	// Header "Mandatory" is 9 characters, always bigger than possible boolean values (true or false)
	sizeMandatory float64 = 9
)

func convertToKubeTypes(kt []parser.KubeTypes) []kubeTypes {
	kubeDocs := make([]kubeTypes, len(kt))
	for idx, kubeType := range kt {
		var k kubeTypes
		var kItems []kubeItem
		for i, item := range kubeType {
			if i != 0 {
				kItems = append(kItems, kubeItem{
					Name:      item.Name,
					Doc:       item.Doc,
					Type:      item.Type,
					RawType:   item.RawType,
					Mandatory: item.Mandatory,
				})
				k.maxSizeOfTableFields[indexName] = math.Max(k.maxSizeOfTableFields[0], float64(len(item.Name)))
				k.maxSizeOfTableFields[indexDoc] = math.Max(k.maxSizeOfTableFields[1], float64(len(item.Doc)))
				k.maxSizeOfTableFields[indexRawType] = math.Max(k.maxSizeOfTableFields[2], float64(len(item.RawType)))
			} else {
				k.Name = item.Name
				k.NameWithAnchor = item.NameWithAnchor
				k.Doc = item.Doc
				k.Items = kItems
			}
		}
		k.Items = kItems
		k.maxSizeOfTableFields[indexMandatory] = sizeMandatory
		kubeDocs[idx] = k
	}
	return kubeDocs
}

func rightPad(s string, pLen int) string {
	return s + strings.Repeat(" ", pLen)
}
