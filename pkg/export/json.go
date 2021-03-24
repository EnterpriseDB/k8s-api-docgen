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
	"encoding/json"

	"github.com/EnterpriseDB/k8s-api-docgen/pkg/parser"
)

// ToJSON get a slice of KubeTypes as input and return the JSON documentation.
// JSON fields are the ones defined in kubeTypes (and kubeItem) definition.
func ToJSON(kt []parser.KubeTypes) (string, error) {
	kubeDocs := convertToKubeTypes(kt)

	j, err := json.MarshalIndent(kubeDocs, "", "\t")
	return string(j), err
}
