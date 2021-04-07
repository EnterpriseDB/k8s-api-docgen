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

	"gopkg.in/yaml.v2"

	"github.com/EnterpriseDB/k8s-api-docgen/pkg/parser"
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
	TableFieldMandatory       string

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

// Markdown configuration to be provided via YAML file
type mdConfiguration struct {
	TableFieldName      string            `yaml:"name,omitempty"`
	TableFieldDoc       string            `yaml:"doc,omitempty"`
	TableFieldRawType   string            `yaml:"type,omitempty"`
	TableFieldMandatory string            `yaml:"mandatory,omitempty"`
	K8sURL              string            `yaml:"k8s_url,omitempty"`
	Version             string            `yaml:"version,omitempty"`
	Sections            map[string]string `yaml:"sections,omitempty"`
}

var conf mdConfiguration

// ToMd gets a slice of KubeTypes and the path to YAML file of the Markdown configuration.
// It returns the Markdown documentation.
func ToMd(kt parser.KubeTypes, mdConfiguration string) (string, error) {
	if mdConfiguration != "" {
		yamlFile, err := ioutil.ReadFile(mdConfiguration) // #nosec
		if err != nil {
			return "", err
		}

		if err = yaml.Unmarshal(yamlFile, &conf); err != nil {
			return "", err
		}
	} else {
		conf.TableFieldName = "Name"
		conf.TableFieldDoc = "Doc"
		conf.TableFieldRawType = "Type"
		conf.TableFieldMandatory = "Mandatory"
	}

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
		kubeDocs[i].TableFieldName = rightPad(conf.TableFieldName, abs(nameMaxLength-len(conf.TableFieldName)))
		kubeDocs[i].TableFieldNameDashSize = strings.Repeat("-", nameMaxLength)
		kubeDocs[i].TableFieldDoc = rightPad(conf.TableFieldDoc, abs(docMaxLength-len(conf.TableFieldDoc)))
		kubeDocs[i].TableFieldDocDashSize = strings.Repeat("-", docMaxLength)
		kubeDocs[i].TableFieldRawType = rightPad(conf.TableFieldRawType, abs(rawTypeMaxLength-len(conf.TableFieldRawType)))
		kubeDocs[i].TableFieldRawTypeDashSize = strings.Repeat("-", rawTypeMaxLength)
		kubeDocs[i].TableFieldMandatory = rightPad(conf.TableFieldMandatory, abs(nameMaxLength-len(conf.TableFieldMandatory)))
		for j, item := range k.Items {
			kubeDocs[i].Items[j].Name = rightPad(item.Name, nameMaxLength-len(item.Name))
			kubeDocs[i].Items[j].Doc = rightPad(item.Doc, docMaxLength-len(item.Doc))
			// adding hyperlinks to documented keys
			kubeDocs[i].Items[j].RawType = rightPad(kubeDocs[i].Items[j].RawType, rawTypeMaxLength-len(item.RawType))
			kubeDocs[i].Items[j].Mandatory = item.Mandatory
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
		section, ok := conf.Sections[info.BaseType]
		if ok {
			return fmt.Sprintf(`[%v](%v/%v/%v)`, info.Name, conf.K8sURL, conf.Version, section)
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
