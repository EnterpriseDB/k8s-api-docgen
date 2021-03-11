// Package export contain the code exporting the internal data to the user
package export

import (
	"encoding/json"

	"github.com/EnterpriseDB/k8s-api-docgen/pkg/parser"
)

// k8s types for generation of docs
type kubeTypes struct {
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

// ToJSON get a slice of KubeTypes as input and return the JSON documentation.
// JSON fields are the ones defined in kubeTypes (and kubeItem) definition.
func ToJSON(kt []parser.KubeTypes) (string, error) {
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
					Mandatory: item.Mandatory,
				})
			} else {
				k.Name = item.Name
				k.Doc = item.Doc
				k.Items = kItems
			}
		}
		k.Items = kItems
		kubeDocs[idx] = k
	}

	j, err := json.MarshalIndent(kubeDocs, "", "\t")
	return string(j), err
}
