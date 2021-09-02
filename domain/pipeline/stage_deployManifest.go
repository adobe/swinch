/*
Copyright 2021 Adobe. All rights reservedm.
This file is licensed to you under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License. You may obtain a copy
of the License at http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under
the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR REPRESENTATIONS
OF ANY KIND, either express or impliedm. See the License for the specific language
governing permissions and limitations under the License.
*/

package pipeline

import (
	"encoding/json"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"strconv"
	"swinch/domain/datastore"
	"swinch/domain/stage"
)

type DeployManifest struct {
	Name                 string   `yaml:"name" json:"name"`
	Type                 string   `yaml:"type,omitempty" json:"type,omitempty"`
	RefId                string   `yaml:"refId,omitempty" json:"refId"`
	RequisiteStageRefIds []string `yaml:"requisiteStageRefIds" json:"requisiteStageRefIds"`

	Account                  string   `yaml:"account,omitempty" json:"account,omitempty"`
	CloudProvider            string   `json:"cloudProvider"`
	ManifestArtifactId       string   `json:"manifestArtifactId"`
	Moniker                  *Moniker `yaml:"moniker,omitempty" json:"moniker,omitempty"`
	NamespaceOverride        string   `json:"namespaceOverride"`
	Overrides                struct{} `yaml:"overrides,omitempty" json:"overrides,omitempty"`
	Source                   string   `json:"source"`
	SkipExpressionEvaluation bool     `yaml:"skipExpressionEvaluation,omitempty" json:"skipExpressionEvaluation,omitempty"`

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

	// Swinch only field
	BakeStageRefIds *int `yaml:"bakeStageRefIds,omitempty" json:"-"`
}

// Moniker is part of Stages
type Moniker struct {
	App string `yaml:"app" json:"app"`
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

func (dm DeployManifest) ProcessDeployManifest(p *Pipeline, stageMap *map[string]interface{}, metadata *stage.Stage) {
	dm.decode(stageMap)
	dm.expand(p, metadata)
	dm.update(stageMap)
}

func (dm *DeployManifest) decode(stageMap *map[string]interface{}) {
	decoderConfig := mapstructure.DecoderConfig{WeaklyTypedInput: true, Result: &dm}
	decoder, err := mapstructure.NewDecoder(&decoderConfig)
	if err != nil {
		log.Fatalf("err: %v", err)
	}

	err = decoder.Decode(stageMap)
	if err != nil {
		log.Fatalf("err: %v", err)
	}
}

func (dm *DeployManifest) expand(p *Pipeline, metadata *stage.Stage) {
	dm.Moniker = new(Moniker)
	dm.Moniker.App = p.Manifest.Metadata.Application
	bakeStageIndex := new(int)

	// Bind deploy stage to a specific bake
	if dm.BakeStageRefIds == nil {
		// Presume a deploy stage has the bake stage as the first element in RequisiteStageRefIds
		*bakeStageIndex, _ = strconv.Atoi(dm.RequisiteStageRefIds[0])
	} else {
		*bakeStageIndex = *dm.BakeStageRefIds
	}

	// Convert from Spinnaker human readable indexing
	*bakeStageIndex -= 1

	//TODO get the bake stage without decoding
	bake := new(BakeManifest)
	err := mapstructure.Decode(p.Manifest.Spec.Stages[*bakeStageIndex], bake)
	if err != nil {
		log.Fatalf("err: %v", err)
	}

	dm.ManifestArtifactId = bake.ExpectedArtifacts[0].Id

	// RefId is either specified by the user or generated based on the stage index
	dm.RefId = metadata.RefId
}

func (dm *DeployManifest) update(stageMap *map[string]interface{}) {
	d := datastore.Datastore{}
	buffer := new(map[string]interface{})
	err := json.Unmarshal(d.MarshalJSON(dm), buffer)
	if err != nil {
		log.Fatalf("Failed to unmarshal JSON:  %v", err)
	}
	*stageMap = *buffer
}
