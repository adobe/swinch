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
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
)

type Stage struct {
	// Metadata Common and Spec will be decoded into the final Stage struct
	// "squash" will nest keys from Metadata and Common struct directly under Stage
	Metadata `mapstructure:",squash"`
	Common   `mapstructure:",squash"`
	// Stage specific fields
	Spec map[string]interface{} `mapstructure:",remain"`

	// Separate maps that will get decoded into proper stage struct and discarded
	ManifestMetadata
	// Map for lookup on other referenced stages
	AllStages *[]map[string]interface{}
	// After processing the stage overwrite it's initial state
	InitStage *map[string]interface{}
}

// ManifestMetadata propagates the metadata from the manifest in the stage
type ManifestMetadata struct {
	Name        string
	Application string
}

type Metadata struct {
	Name                 string   `yaml:"name" json:"name"`
	Type                 string   `yaml:"type,omitempty" json:"type,omitempty"`
	RefId                string   `yaml:"refId,omitempty" json:"refId,omitempty"`
	RequisiteStageRefIds []string `yaml:"requisiteStageRefIds" json:"requisiteStageRefIds"`
}

type Common struct {
	StageFails `mapstructure:",squash"`

	RestrictExecutionDuringTimeWindow bool `yaml:"restrictExecutionDuringTimeWindow,omitempty" json:"restrictExecutionDuringTimeWindow,omitempty"`

	RestrictedExecutionWindow *RestrictedExecutionWindow `yaml:"restrictedExecutionWindow,omitempty" json:"restrictedExecutionWindow,omitempty"`
	SkipWindowText            string                     `yaml:"skipWindowText,omitempty" json:"skipWindowText,omitempty"`
	FailOnFailedExpressions   bool                       `yaml:"failOnFailedExpressions,omitempty" json:"failOnFailedExpressions,omitempty"`

	StageEnabled *StageEnabled `yaml:"stageEnabled,omitempty" json:"stageEnabled,omitempty"`

	SendNotifications bool           `yaml:"sendNotifications,omitempty" json:"sendNotifications,omitempty"`
	Notifications     *Notifications `yaml:"notifications,omitempty" json:"notifications,omitempty"`

	Comments string `yaml:"comments,omitempty" json:"comments,omitempty"`
}

type StageFails struct {
	ContinuePipeline              *bool `yaml:"-" json:"continuePipeline,omitempty"`
	FailPipeline                  *bool `yaml:"-" json:"failPipeline,omitempty"`
	CompleteOtherBranchesThenFail *bool `yaml:"-" json:"completeOtherBranchesThenFail,omitempty"`

	// IfStageFails key only in swinch
	IfStageFails string `yaml:"ifStageFails,omitempty" json:"-"`
}

type RestrictedExecutionWindow struct {
	Days      []int `yaml:"days,omitempty" json:"days,omitempty"`
	Whitelist []struct {
		EndHour   string `yaml:"endHour,omitempty" json:"endHour,omitempty"`
		EndMin    string `yaml:"endMin,omitempty" json:"endMin,omitempty"`
		StartHour string `yaml:"startHour,omitempty" json:"startHour,omitempty"`
		StartMin  string `yaml:"startMin,omitempty" json:"startMin,omitempty"`
	} `yaml:"whitelist,omitempty" json:"whitelist,omitempty"`
	Jitter struct {
		Enabled    bool   `yaml:"enabled,omitempty" json:"enabled,omitempty"`
		MaxDelay   string `yaml:"maxDelay,omitempty" json:"maxDelay,omitempty"`
		MinDelay   string `yaml:"minDelay,omitempty" json:"minDelay,omitempty"`
		SkipManual bool   `yaml:"skipManual,omitempty" json:"skipManual,omitempty"`
	} `yaml:"jitter,omitempty" json:"jitter,omitempty"`
}

type StageEnabled struct {
	Expression string `yaml:"expression,omitempty" json:"expression,omitempty"`
	Type       string `yaml:"type,omitempty" json:"type,omitempty"`
}

type Notifications []struct {
	Address string   `yaml:"address,omitempty" json:"address,omitempty"`
	Level   string   `yaml:"level,omitempty" json:"level,omitempty"`
	Type    string   `yaml:"type,omitempty" json:"type,omitempty"`
	When    []string `yaml:"when,omitempty" json:"when,omitempty"`
	Message *Message `yaml:"message,omitempty" json:"message,omitempty"`
}

type Message struct {
	StageComplete struct {
		Text string `yaml:"text,omitempty" json:"text,omitempty"`
	} `yaml:"stageComplete,omitempty" json:"stage.complete,omitempty"`
	StageFailed struct {
		Text string `yaml:"text,omitempty" json:"text,omitempty"`
	} `yaml:"stageFailed,omitempty" json:"stage.failed,omitempty"`
	StageStarting struct {
		Text string `yaml:"text,omitempty" json:"text,omitempty"`
	} `yaml:"stageStarting,omitempty" json:"stage.starting,omitempty"`
}

func (s Stage) Decode(stage *map[string]interface{}) Stage {
	tmp := new(Stage)
	decoderConfig := mapstructure.DecoderConfig{WeaklyTypedInput: true, Result: &tmp}
	decoder, err := mapstructure.NewDecoder(&decoderConfig)
	if err != nil {
		log.Fatalf("err: %v", err)
	}

	err = decoder.Decode(stage)
	if err != nil {
		log.Fatalf("err: %v", err)
	}
	return *tmp
}

// FailStageSetter method reduces the complexity of "If stage fails" execution option
// "If stage fails" execution option has 4 scenarios as seen in the WebUI; to set one of them a bool combination of the below parameters is needed
// to avoid complexity, the user will use ONLY the ifStageFails parameter (which exists only in the yaml)
func (s *Stage) FailStageSetter() {
	s.ContinuePipeline = new(bool)
	s.FailPipeline = new(bool)
	s.CompleteOtherBranchesThenFail = new(bool)

	switch s.IfStageFails {
	case "halt the entire pipeline":
		*s.ContinuePipeline = false
		*s.FailPipeline = true
		*s.CompleteOtherBranchesThenFail = false
	case "halt this branch of the pipeline":
		*s.ContinuePipeline = false
		*s.FailPipeline = false
		*s.CompleteOtherBranchesThenFail = false
	case "halt this branch and fail the pipeline once other branches complete":
		*s.ContinuePipeline = false
		*s.FailPipeline = false
		*s.CompleteOtherBranchesThenFail = true
	case "ignore the failure":
		*s.ContinuePipeline = true
		*s.FailPipeline = false
		*s.CompleteOtherBranchesThenFail = false
	// without these defaults, if the ifStageFails parameters is not set inside the chart,
	// the default option will be "halt this branch of the pipeline" instead of "halt the entire pipeline"
	default:
		*s.ContinuePipeline = false
		*s.FailPipeline = true
		*s.CompleteOtherBranchesThenFail = false
	}
}
