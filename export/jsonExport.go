package export

import (
	"cnp-docgen/parser"
	"encoding/json"
	"strings"
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

// Get a slice of KubeTypes as input and return the JSON documentation.
// JSON fields are the ones defined in kubeTypes (and kubeItem) definition.
func ToJSON(kt []parser.KubeTypes) (string, error) {
	var kubeDocs []kubeTypes
	for _, kubeType := range kt {
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
		kubeDocs = append(kubeDocs, k)
	}
	// to JSON with indentation
	j, err := json.MarshalIndent(kubeDocs, "", "\t")
	// trim "\n"
	return strings.Replace(string(j), `\n`, "", -1), err
}
