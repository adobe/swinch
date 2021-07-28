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

type PipelineSpec struct {
	Application          string        `yaml:"-" json:"application"`
	Name                 string        `yaml:"-" json:"name"`
	Index                int           `yaml:"index" json:"index"`
	KeepWaitingPipelines bool          `yaml:"keepWaitingPipelines,omitempty" json:"keepWaitingPipelines,omitempty"`
	LimitConcurrent      bool          `yaml:"limitConcurrent,omitempty" json:"limitConcurrent,omitempty"`
	SpelEvaluator        string        `yaml:"spelEvaluator,omitempty" json:"spelEvaluator,omitempty"`
	Stages               []map[string]interface{} `yaml:"stages" json:"stages"`
	Triggers             []interface{} `yaml:"triggers,omitempty" json:"triggers,omitempty""`
}


func (s PipelineSpec) LoadSpec(spec []byte) PipelineSpec {
	err := json.Unmarshal(spec, &s)

	if err != nil {
		log.Fatalf("Error LoadSpec: %v", err)
	}
	return s
}
