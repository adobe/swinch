/*
Copyright 2021 Adobe. All rights reserved.
This file is licensed to you under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License. You may obtain a copy
of the License at http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under
the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR REPRESENTATIONS
OF ANY KIND, either express or implied. See the License for the specific language
governing permissions and limitations under the License.
*/

package domain

import (
	"fmt"
)

var Kinds = map[string]string{
	ApplicationKind: ApplicationAPI,
	PipelineKind:    PipelineAPI,
}

type Manifest struct {
	ApiVersion string      `yaml:"apiVersion" json:"apiVersion"`
	Kind       string      `yaml:"kind" json:"kind"`
	Metadata   interface{} `yaml:"metadata" json:"metadata"`
	Spec       interface{} `yaml:"spec" json:"spec"`
}

func (m Manifest) GetManifests(filePath string) ([]ApplicationManifest, []PipelineManifest) {
	d := Datastore{}
	discoveredYAMLDocs := d.DiscoverYAMLFiles(filePath)
	return d.ReadYAMLDocs(discoveredYAMLDocs)
}

func (m Manifest) Validate() error {
	_, ok := Kinds[m.Kind]
	if ok {
		kindApiVersion := Kinds[m.Kind]
		if m.ApiVersion != kindApiVersion {
			return fmt.Errorf("bad api version, expected: %v, got: %v", kindApiVersion, m.ApiVersion)
		} else {
			return nil
		}
	} else {
		return fmt.Errorf("unknown manifest Kind: %v", m.Kind)
	}
}
