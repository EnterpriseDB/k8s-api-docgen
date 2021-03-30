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

// Package md contain the code exporting the internal data to markdown
package md

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"text/template"

	"github.com/EnterpriseDB/k8s-api-docgen/pkg/parser"
)

const (
	tableFieldName    string = "Name"
	tableFieldDoc     string = "Doc"
	tableFieldRawType string = "Type"

	docPrefix = "https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.20/"
)

var (
	kubernetesLinks = map[string]string{
		"metav1.ObjectMeta":                docPrefix + "#objectmeta-v1-meta",
		"metav1.ListMeta":                  docPrefix + "#listmeta-v1-meta",
		"metav1.LabelSelector":             docPrefix + "#labelselector-v1-meta",
		"metav1.Time":                      docPrefix + "#time-v1-meta",
		"v1.ResourceRequirements":          docPrefix + "#resourcerequirements-v1-core",
		"v1.LocalObjectReference":          docPrefix + "#localobjectreference-v1-core",
		"v1.SecretKeySelector":             docPrefix + "#secretkeyselector-v1-core",
		"v1.PersistentVolumeClaim":         docPrefix + "#persistentvolumeclaim-v1-core",
		"v1.EmptyDirVolumeSource":          docPrefix + "#emptydirvolumesource-v1-core",
		"apiextensionsv1.JSON":             docPrefix + "#json-v1-apiextensions-k8s-io",
		"corev1.LocalObjectReference":      docPrefix + "#localobjectreference-v1-core",
		"corev1.ResourceRequirements":      docPrefix + "#resourcerequirements-v1-core",
		"corev1.PersistentVolumeClaimSpec": docPrefix + "#persistentvolumeclaim-v1-core",
		"corev1.SecretKeySelector":         docPrefix + "#secretkeyselector-v1-core",
		"corev1.ConfigMapKeySelector":      docPrefix + "#configmapkeyselector-v1-core",
	}
)

// k8s types for generation of docs
type kubeType struct {
	Name                      string
	NameWithAnchor            string
	Doc                       string
	Items                     []kubeItem
	TableFieldName            string
	TableFieldNameDashSize    string
	TableFieldDoc             string
	TableFieldDocDashSize     string
	TableFieldRawType         string
	TableFieldRawTypeDashSize string

	maxSizeOfName    int
	maxSizeOfDoc     int
	maxSizeOfRawType int
}

// k8s items
type kubeItem struct {
	Name      string
	Doc       string
	Type      string
	RawType   string
	Mandatory bool
}

// ToMd get a slice of KubeTypes as input and return the Markdown documentation.
func ToMd(kt parser.KubeTypes) (string, error) {
	kubeDocs := convertToKubeTypes(kt)
	format(kubeDocs)

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

func convertToKubeTypes(kt parser.KubeTypes) []kubeType {
	internalTypes := make(map[string]bool)
	for _, kubeStructure := range kt {
		internalTypes[kubeStructure.Name] = true
	}

	kubeDocs := make([]kubeType, len(kt))
	for idx, kubeStructure := range kt {
		k := kubeType{
			Name:                      kubeStructure.Name,
			NameWithAnchor:            applyAnchor(kubeStructure.Name),
			Doc:                       kubeStructure.Doc,
			Items:                     nil,
			TableFieldName:            "",
			TableFieldNameDashSize:    "",
			TableFieldDoc:             "",
			TableFieldDocDashSize:     "",
			TableFieldRawType:         "",
			TableFieldRawTypeDashSize: "",
		}

		var kItems []kubeItem
		for _, item := range kubeStructure.Fields {
			typeField := wrapInLink(item.Type, internalTypes)
			kItems = append(kItems, kubeItem{
				Name:      item.Name,
				Doc:       item.Doc,
				Type:      item.Type.Name,
				RawType:   typeField,
				Mandatory: item.Mandatory,
			})

			k.maxSizeOfName = max(k.maxSizeOfName, len(item.Name))
			k.maxSizeOfDoc = max(k.maxSizeOfDoc, len(item.Doc))
			k.maxSizeOfRawType = max(k.maxSizeOfRawType, len(typeField))
		}
		k.Items = kItems
		kubeDocs[idx] = k
	}
	return kubeDocs
}

// format applies proper formats to tables and documentation
func format(kubeDocs []kubeType) {
	for i, k := range kubeDocs {
		kubeDocs[i].Doc = strings.Trim(k.Doc, "\n")
		nameMaxLength := k.maxSizeOfName
		docMaxLength := k.maxSizeOfDoc
		rawTypeMaxLength := k.maxSizeOfRawType
		kubeDocs[i].TableFieldName = rightPad(tableFieldName, abs(nameMaxLength-len(tableFieldName)))
		kubeDocs[i].TableFieldNameDashSize = strings.Repeat("-", nameMaxLength)
		kubeDocs[i].TableFieldDoc = rightPad(tableFieldDoc, abs(docMaxLength-len(tableFieldDoc)))
		kubeDocs[i].TableFieldDocDashSize = strings.Repeat("-", docMaxLength)
		kubeDocs[i].TableFieldRawType = rightPad(tableFieldRawType, abs(rawTypeMaxLength-len(tableFieldRawType)))
		kubeDocs[i].TableFieldRawTypeDashSize = strings.Repeat("-", rawTypeMaxLength)
		for j, item := range k.Items {
			kubeDocs[i].Items[j].Name = rightPad(item.Name, nameMaxLength-len(item.Name))
			kubeDocs[i].Items[j].Doc = rightPad(item.Doc, docMaxLength-len(item.Doc))
			// adding hyperlinks to documented keys
			kubeDocs[i].Items[j].RawType = rightPad(kubeDocs[i].Items[j].RawType, rawTypeMaxLength-len(item.RawType))
		}
	}
}

// runTemplate execute the template, fed by docs values
func runTemplate(aTemplate []byte, docs []kubeType) (string, error) {
	var w bytes.Buffer
	tmpl, err := template.New("KubeTypes").Parse(string(aTemplate))
	if err != nil {
		return "", err
	}
	err = tmpl.Execute(&w, docs)
	if err != nil {
		return "", err
	}
	return w.String(), nil
}

// applyAnchor applies an anchor to name, in order to be compliant with MarkDown output
func applyAnchor(name string) string {
	return fmt.Sprintf("<a name='%v'></a> `%v`", name, name)
}

// wrapInLink generate a Markdown link tag from a type
func wrapInLink(info parser.TypeInfo, internalTypes map[string]bool) string {
	if info.Internal && strings.Title(info.BaseType) == info.BaseType {
		// This is an internal type exported, so it is user-defined.
		// Is this a documented type or not?
		_, documented := internalTypes[info.BaseType]

		if documented {
			// Let's use an internal link for that
			return fmt.Sprintf("[%v](#%v)", info.Name, info.BaseType)
		}

		// We don't have documentation for this type, so we are leaving
		// it unlinked
		return info.Name
	}

	if !info.Internal {
		// This is an external type so let's hope it is a Kubernetes native one
		link, ok := kubernetesLinks[info.BaseType]
		if ok {
			return fmt.Sprintf("[%v](%v)", info.Name, link)
		}
	}

	return info.Name
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func rightPad(s string, pLen int) string {
	return s + strings.Repeat(" ", pLen)
}
