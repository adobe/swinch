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

type ManualJudgment struct {
	Name                 string   `yaml:"name" json:"name"`
	Type                 string   `yaml:"type,omitempty" json:"type,omitempty"`
	RefId                string   `yaml:"refId,omitempty" json:"refId"`
	RequisiteStageRefIds []string `yaml:"requisiteStageRefIds" json:"requisiteStageRefIds"`

	IsNew                          bool          `yaml:"isNew,omitempty" json:"isNew,omitempty"`
	JudgmentInputs                 []interface{} `yaml:"judgmentInputs" json:"judgmentInputs"`
	PropagateAuthenticationContext bool          `yaml:"propagateAuthenticationContext" json:"propagateAuthenticationContext"`
	SelectedStageRoles             []string      `yaml:"selectedStageRoles" json:"selectedStageRoles"`
	Instructions                   string        `yaml:"instructions" json:"instructions"`

	ContinuePipeline              bool          `yaml:"continuePipeline,omitempty" json:"continuePipeline,omitempty"`
	FailPipeline                  bool          `yaml:"failPipeline,omitempty" json:"failPipeline,omitempty"`
	CompleteOtherBranchesThenFail bool          `yaml:"completeOtherBranchesThenFail,omitempty" json:"completeOtherBranchesThenFail,omitempty"`
	StageTimeoutMs                int           `yaml:"stageTimeoutMs" json:"stageTimeoutMs"`
	Notifications                 []interface{} `yaml:"notifications" json:"notifications"`
}

func (mj ManualJudgment) ProcessManualJudgment(stageMap *map[string]interface{}, metadata *stage.Stage) {
	mj.decode(stageMap)
	mj.RefId = metadata.RefId
	mj.update(stageMap)
}

func (mj *ManualJudgment) decode(stageMap *map[string]interface{}) {
	decoderConfig := mapstructure.DecoderConfig{WeaklyTypedInput: true, Result: &mj}
	decoder, err := mapstructure.NewDecoder(&decoderConfig)
	if err != nil {
		log.Fatalf("err: %v", err)
	}

	err = decoder.Decode(stageMap)
	if err != nil {
		log.Fatalf("err: %v", err)
	}
}

func (mj *ManualJudgment) update(stageMap *map[string]interface{}) {
	d := datastore.Datastore{}
	buffer := new(map[string]interface{})
	err := json.Unmarshal(d.MarshalJSON(mj), buffer)
	if err != nil {
		log.Fatalf("Failed to unmarshal JSON:  %v", err)
	}
	*stageMap = *buffer
}
