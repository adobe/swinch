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
	"swinch/domain/util"
)

const StageType = "bakeManifest"

type BakeManifest struct {
	Metadata `mapstructure:",squash"`
	Common   `mapstructure:",squash"`

	Account                     string              `yaml:"account,omitempty" json:"account,omitempty"`
	OutputName                  string              `json:"outputName"`
	ExpectedArtifacts           []ExpectedArtifacts `yaml:"expectedArtifacts,omitempty" json:"expectedArtifacts,omitempty"`
	InputArtifacts              []InputArtifacts    `yaml:"inputArtifacts,omitempty" json:"inputArtifacts,omitempty"`
	ManifestArtifactId          string              `json:"manifestArtifactId,omitempty"`
	Namespace                   string              `json:"namespace"`
	TemplateRenderer            string              `json:"templateRenderer"`
	Overrides                   struct{}            `yaml:"overrides,omitempty" json:"overrides,omitempty"`
	RawOverrides                bool                `yaml:"rawOverrides,omitempty" json:"rawOverrides,omitempty"`
	EvaluateOverrideExpressions bool                `yaml:"evaluateOverrideExpressions,omitempty" json:"evaluateOverrideExpressions,omitempty"`
}

type ExpectedArtifacts struct {
	DefaultArtifact    *DefaultArtifact `yaml:"defaultArtifact,omitempty" json:"defaultArtifact,omitempty"`
	DisplayName        string           `yaml:"displayName" json:"displayName"`
	Id                 string           `yaml:"-" json:"id"` // swinch generated
	MatchArtifact      MatchArtifact    `yaml:"matchArtifact,omitempty" json:"matchArtifact,omitempty"`
	UseDefaultArtifact *bool            `yaml:"useDefaultArtifact,omitempty" json:"useDefaultArtifact,omitempty"`
	UsePriorArtifact   *bool            `yaml:"usePriorArtifact,omitempty" json:"usePriorArtifact,omitempty"`
}

// DefaultArtifact is part of ExpectedArtifacts
type DefaultArtifact struct {
	CustomKind bool   `yaml:"customKind" json:"customKind"`
	Id         string `yaml:"id" json:"id"`
}

// MatchArtifact is part of ExpectedArtifacts
type MatchArtifact struct {
	ArtifactAccount string  `yaml:"artifactAccount" json:"artifactAccount"`
	CustomKind      *bool   `yaml:"customKind,omitempty" json:"customKind,omitempty"`
	Id              *string `yaml:"-" json:"id,omitempty"`
	Name            string  `yaml:"name" json:"name"`
	Type            string  `yaml:"type" json:"type"`
}

// InputArtifacts is part of Stages
type InputArtifacts struct {
	Account  string `yaml:"account" json:"account"`
	Artifact struct {
		ArtifactAccount string `yaml:"artifactAccount" json:"artifactAccount"`
		Id              string `yaml:"-" json:"id,omitempty"`
		Name            string `yaml:"name" json:"name"`
		Type            string `yaml:"type" json:"type"`
		Version         string `yaml:"version" json:"version"`
	} `yaml:"artifact" json:"artifact"`
}

func (bm BakeManifest) GetStageType() string {
	return StageType
}

func (bm BakeManifest) MakeStage(stage *Stage) *map[string]interface{} {
	bm.decode(stage)
	bm.expand()
	return bm.encode()
}

func (bm *BakeManifest) decode(stage *Stage) {
	decoderConfig := mapstructure.DecoderConfig{WeaklyTypedInput: true, Result: &bm}
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

func (bm *BakeManifest) expand() {
	u := util.Util{}

	// TODO check that index on ExpectedArtifacts is always 0
	expectArtifacts := &bm.ExpectedArtifacts[0]
	// expectArtifacts ID is used by the deploy stage
	expectArtifacts.Id = u.GenerateUUID(expectArtifacts.DisplayName + bm.Name).String()

	// TODO make sure MatchArtifact ID is not used
	//expectArtifacts.MatchArtifact.Id = bm.newUUID(expectArtifacts.MatchArtifact.Name+expectArtifacts.MatchArtifact.Type).String()

	// TODO check that index on InputArtifacts is always 0
	inputArtifacts := &bm.InputArtifacts[0]
	//Deduplicate ArtifactAccount name
	inputArtifacts.Artifact.ArtifactAccount = inputArtifacts.Account
	// inputArtifacts.Artifact.Id not mandatory
	//inputArtifacts.Artifact.Id = bm.newUUID(inputArtifacts.Artifact.Name + inputArtifacts.Artifact.Version).String()
}

func (bm *BakeManifest) encode() *map[string]interface{} {
	d := datastore.Datastore{}
	stage := new(map[string]interface{})
	err := json.Unmarshal(d.MarshalJSON(bm), stage)
	if err != nil {
		log.Fatalf("Failed to unmarshal JSON:  %v", err)
	}
	return stage
}
