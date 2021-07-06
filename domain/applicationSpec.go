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
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

type ApplicationSpec struct {
	CloudProviders string      `yaml:"cloudProviders" json:"cloudProviders"`
	Email          string      `yaml:"email" json:"email"`
	Name           string      `yaml:"-" json:"name"`
	Permissions    Permissions `yaml:"permissions" json:"permissions"`
}

type Permissions struct {
	EXECUTE []string `yaml:"EXECUTE" json:"EXECUTE"`
	READ    []string `yaml:"READ" json:"READ"`
	WRITE   []string `yaml:"WRITE" json:"WRITE"`
}

func (a ApplicationSpec) LoadSpec(spec []byte) ApplicationSpec {
	err := json.Unmarshal(spec, &a)
	if err != nil {
		log.Fatalf("Error LoadSpec: %v", err)
	}
	return a
}
