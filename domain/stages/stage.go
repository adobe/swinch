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
	// "squash" will nest keys from Metadata struct directly under Stage
	Metadata `mapstructure:",squash"`
	// Separate maps that will get decoded into proper stage struct and discarded
	ManifestMetadata
	Spec     map[string]interface{} `mapstructure:",remain"`
	RawStage *map[string]interface{}
	Stages   *[]map[string]interface{}
}

type Metadata struct {
	Name                 string   `yaml:"name" json:"name"`
	Type                 string   `yaml:"type,omitempty" json:"type,omitempty"`
	RefId                string   `yaml:"refId,omitempty" json:"refId,omitempty"`
	RequisiteStageRefIds []string `yaml:"requisiteStageRefIds" json:"requisiteStageRefIds"`
}

type ManifestMetadata struct {
	// ManifestMetadata
	Name        string
	Application string
}

type Common struct {
	// From here to the end the fields are common except stageTimeoutMs (fail stage after specified time) and BakeStageRefIds
	ContinuePipeline                  bool `yaml:"continuePipeline,omitempty" json:"continuePipeline,omitempty"`
	FailPipeline                      bool `yaml:"failPipeline,omitempty" json:"failPipeline,omitempty"`
	CompleteOtherBranchesThenFail     bool `yaml:"completeOtherBranchesThenFail,omitempty" json:"completeOtherBranchesThenFail,omitempty"`
	RestrictExecutionDuringTimeWindow bool `yaml:"restrictExecutionDuringTimeWindow,omitempty" json:"restrictExecutionDuringTimeWindow,omitempty"`

	RestrictedExecutionWindow *RestrictedExecutionWindow `yaml:"restrictedExecutionWindow,omitempty" json:"restrictedExecutionWindow,omitempty"`
	SkipWindowText            string                     `yaml:"skipWindowText,omitempty" json:"skipWindowText,omitempty"`
	// StageTimeoutMs applies only to select stages (Deploy, Manual Judgement, Run Job etc.)
	StageTimeoutMs          *int `yaml:"stageTimeoutMs,omitempty" json:"stageTimeoutMs,omitempty"`
	FailOnFailedExpressions bool `yaml:"failOnFailedExpressions,omitempty" json:"failOnFailedExpressions,omitempty"`

	StageEnabled *StageEnabled `yaml:"stageEnabled,omitempty" json:"stageEnabled,omitempty"`

	SendNotifications bool           `yaml:"sendNotifications,omitempty" json:"sendNotifications,omitempty"`
	Notifications     *Notifications `yaml:"notifications,omitempty" json:"notifications,omitempty"`

	Comments string `yaml:"comments,omitempty" json:"comments,omitempty"`
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

func (s *Stage) GetStage(stage *map[string]interface{}) Stage {
	s.RawStage = stage
	s.decode(stage)
	return *s
}

func (s *Stage) decode(stage *map[string]interface{}) {
	decoderConfig := mapstructure.DecoderConfig{WeaklyTypedInput: true, Result: &s}
	decoder, err := mapstructure.NewDecoder(&decoderConfig)
	if err != nil {
		log.Fatalf("err: %v", err)
	}

	err = decoder.Decode(stage)
	if err != nil {
		log.Fatalf("err: %v", err)
	}
}
