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

package pipeline

import (
	"encoding/json"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"swinch/domain/datastore"
	"swinch/domain/stage"
)

type Wait struct {
	Name                 string   `yaml:"name" json:"name"`
	Type                 string   `yaml:"type,omitempty" json:"type,omitempty"`
	RefId                string   `yaml:"refId,omitempty" json:"refId"`
	RequisiteStageRefIds []string `yaml:"requisiteStageRefIds" json:"requisiteStageRefIds"`

	IsNew        bool   `yaml:"isNew,omitempty" json:"isNew,omitempty"`
	SkipWaitText string `yaml:"skipWaitText" json:"skipWaitText"`
	WaitTime     int    `yaml:"waitTime" json:"waitTime"`

	ContinuePipeline              bool `yaml:"continuePipeline,omitempty" json:"continuePipeline,omitempty"`
	FailPipeline                  bool `yaml:"failPipeline,omitempty" json:"failPipeline,omitempty"`
	CompleteOtherBranchesThenFail bool `yaml:"completeOtherBranchesThenFail,omitempty" json:"completeOtherBranchesThenFail,omitempty"`
}

func (wt Wait) ProcessWait(stageMap *map[string]interface{}, metadata *stage.Stage) {
	wt.decode(stageMap)
	wt.RefId = metadata.RefId
	wt.update(stageMap)
}

func (wt *Wait) decode(stageMap *map[string]interface{}) {
	decoderConfig := mapstructure.DecoderConfig{WeaklyTypedInput: true, Result: &wt}
	decoder, err := mapstructure.NewDecoder(&decoderConfig)
	if err != nil {
		log.Fatalf("err: %v", err)
	}

	err = decoder.Decode(stageMap)
	if err != nil {
		log.Fatalf("err: %v", err)
	}
}

func (wt *Wait) update(stageMap *map[string]interface{}) {
	d := datastore.Datastore{}
	buffer := new(map[string]interface{})
	err := json.Unmarshal(d.MarshalJSON(wt), buffer)
	if err != nil {
		log.Fatalf("Failed to unmarshal JSON:  %v", err)
	}
	*stageMap = *buffer
}
