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

const wait = "wait"

type Wait struct {
	Stage `mapstructure:",squash"`

	IsNew        bool   `yaml:"isNew,omitempty" json:"isNew,omitempty"`
	SkipWaitText string `yaml:"skipWaitText" json:"skipWaitText"`
	WaitTime     int    `yaml:"waitTime" json:"waitTime"`

	ContinuePipeline              bool `yaml:"continuePipeline,omitempty" json:"continuePipeline,omitempty"`
	FailPipeline                  bool `yaml:"failPipeline,omitempty" json:"failPipeline,omitempty"`
	CompleteOtherBranchesThenFail bool `yaml:"completeOtherBranchesThenFail,omitempty" json:"completeOtherBranchesThenFail,omitempty"`
}

func (wt Wait) GetStageType() string {
	return wait
}

func (wt Wait) Process(stage *Stage) {
	wt.decode(stage)
	wt.update(stage)
}

func (wt *Wait) decode(stage *Stage) {
	decoderConfig := mapstructure.DecoderConfig{WeaklyTypedInput: true, Result: &wt}
	decoder, err := mapstructure.NewDecoder(&decoderConfig)

	err = decoder.Decode(stage.Metadata)
	if err != nil {
		log.Fatalf("error decoding stage metadata: %v", err)
	}
	err = decoder.Decode(stage.Spec)
	if err != nil {
		log.Fatalf("error decoding stage spec: %v", err)
	}
}

func (wt *Wait) update(stage *Stage) {
	d := datastore.Datastore{}
	tmpStage := new(map[string]interface{})
	err := json.Unmarshal(d.MarshalJSON(wt), tmpStage)
	if err != nil {
		log.Fatalf("Failed to unmarshal JSON:  %v", err)
	}
	*stage.RawStage = *tmpStage
}
