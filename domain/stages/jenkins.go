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

const jenkins StageType = "jenkins"

type Jenkins struct {
	Metadata `mapstructure:",squash"`
	Common   `mapstructure:",squash"`

	IsNew                    bool     `yaml:"isNew,omitempty" json:"isNew,omitempty"`
	Master                   string   `yaml:"master" json:"master"`
	Job                      string   `yaml:"job" json:"job"`
	Parameters               struct{} `yaml:"parameters" json:"parameters"`
	MarkUnstableAsSuccessful bool     `yaml:"markUnstableAsSuccessful" json:"markUnstableAsSuccessful"`
	WaitForCompletion        bool     `yaml:"waitForCompletion" json:"waitForCompletion"`

	// Overriding the field from Common struct without "omitempty" as it's required by the Jenkins Stage
	ContinuePipeline bool `yaml:"continuePipeline" json:"continuePipeline"`
}

func (jks Jenkins) MakeStage(stage *Stage) *map[string]interface{} {
	jks.decode(stage)
	return jks.encode()
}

func (jks *Jenkins) decode(stage *Stage) {
	decoderConfig := mapstructure.DecoderConfig{WeaklyTypedInput: true, Result: &jks}
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

func (jks *Jenkins) encode() *map[string]interface{} {
	d := datastore.Datastore{}
	stage := new(map[string]interface{})
	err := json.Unmarshal(d.MarshalJSON(jks), stage)
	if err != nil {
		log.Fatalf("Failed to unmarshal JSON:  %v", err)
	}
	return stage
}
