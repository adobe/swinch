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

package stages

import (
	"encoding/json"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"swinch/domain/datastore"
)

const pipeline = "pipeline"

type Pipeline struct {
	Metadata `mapstructure:",squash"`
	Common   `mapstructure:",squash"`

	Application       string `yaml:"application" json:"application"`
	IsNew             bool   `yaml:"isNew,omitempty" json:"isNew,omitempty"`
	Pipeline          string `yaml:"pipeline" json:"pipeline"`
	WaitForCompletion bool   `yaml:"waitForCompletion" json:"waitForCompletion"`

	// Overriding the field from Common struct without "omitempty" as it's required by the Pipeline Stage
	FailPipeline bool `yaml:"failPipeline" json:"failPipeline"`
}

func (pp Pipeline) GetStageType() string {
	return pipeline
}

func (pp Pipeline) MakeStage(stage *Stage) *map[string]interface{} {
	pp.decode(stage)
	return pp.encode()
}

func (pp *Pipeline) decode(stage *Stage) {
	decoderConfig := mapstructure.DecoderConfig{WeaklyTypedInput: true, Result: &pp}
	decoder, err := mapstructure.NewDecoder(&decoderConfig)
	if err != nil {
		log.Fatalf("err: %v", err)
	}

	err = decoder.Decode(stage.Metadata)
	if err != nil {
		log.Fatalf("error decoding stage metadata: %v", err)
	}
	err = decoder.Decode(stage.Spec)
	if err != nil {
		log.Fatalf("error decoding stage spec: %v", err)
	}
}

func (pp *Pipeline) encode() *map[string]interface{} {
	d := datastore.Datastore{}
	stage := new(map[string]interface{})
	err := json.Unmarshal(d.MarshalJSON(pp), stage)
	if err != nil {
		log.Fatalf("Failed to unmarshal JSON:  %v", err)
	}
	return stage
}
