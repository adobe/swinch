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

const manualJudgment = "manualJudgment"

type ManualJudgment struct {
	Metadata  `mapstructure:",squash"`
	Common `mapstructure:",squash"`

	IsNew                          bool          `yaml:"isNew,omitempty" json:"isNew,omitempty"`
	JudgmentInputs                 []interface{} `yaml:"judgmentInputs" json:"judgmentInputs"`
	PropagateAuthenticationContext bool          `yaml:"propagateAuthenticationContext" json:"propagateAuthenticationContext"`
	SelectedStageRoles             []string      `yaml:"selectedStageRoles,omitempty" json:"selectedStageRoles,omitempty"`
	Instructions                   string        `yaml:"instructions" json:"instructions"`

	ContinuePipeline              bool          `yaml:"continuePipeline,omitempty" json:"continuePipeline,omitempty"`
	FailPipeline                  bool          `yaml:"failPipeline,omitempty" json:"failPipeline,omitempty"`
	CompleteOtherBranchesThenFail bool          `yaml:"completeOtherBranchesThenFail,omitempty" json:"completeOtherBranchesThenFail,omitempty"`
	StageTimeoutMs                int           `yaml:"stageTimeoutMs" json:"stageTimeoutMs"`
	Notifications                 []interface{} `yaml:"notifications" json:"notifications"`
}

func (mj ManualJudgment) GetStageType() string {
	return manualJudgment
}

func (mj ManualJudgment) MakeStage(stage *Stage) *map[string]interface{} {
	mj.decode(stage)
	return mj.encode()
}

func (mj *ManualJudgment) decode(stage *Stage) {
	decoderConfig := mapstructure.DecoderConfig{WeaklyTypedInput: true, Result: &mj}
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

func (mj *ManualJudgment) encode() *map[string]interface{} {
	d := datastore.Datastore{}
	stage := new(map[string]interface{})
	err := json.Unmarshal(d.MarshalJSON(mj), stage)
	if err != nil {
		log.Fatalf("Failed to unmarshal JSON:  %v", err)
	}
	return stage
}
