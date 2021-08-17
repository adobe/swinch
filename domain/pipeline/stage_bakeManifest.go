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
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"swinch/domain/datastore"
	"swinch/domain/stage"
)

type BakeManifest struct {
	Name                 string   `yaml:"name" json:"name"`
	Type                 string   `yaml:"type,omitempty" json:"type,omitempty"`
	RefId                string   `yaml:"refId,omitempty" json:"refId"`
	RequisiteStageRefIds []string `yaml:"requisiteStageRefIds" json:"requisiteStageRefIds"`

	Account            string              `yaml:"account,omitempty" json:"account,omitempty"`
	OutputName         string              `json:"outputName"`
	ExpectedArtifacts  []ExpectedArtifacts `yaml:"expectedArtifacts,omitempty" json:"expectedArtifacts,omitempty"`
	InputArtifacts     []InputArtifacts    `yaml:"inputArtifacts,omitempty" json:"inputArtifacts,omitempty"`
	ManifestArtifactId string             `json:"manifestArtifactId,omitempty"`
	Namespace          string              `json:"namespace"`
	TemplateRenderer   string              `json:"templateRenderer"`
	Overrides          struct{}            `yaml:"overrides,omitempty" json:"overrides,omitempty"`
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

func (bm *BakeManifest) ProcessBakeManifest(p *Pipeline, stage *map[string]interface{}, metadata *stage.Stage) {
	bm.decode(stage)
	bm.expand(metadata)
	bm.updateStage(stage)
}

func (bm *BakeManifest) decode(stage *map[string]interface{}) {
	decoderConfig := mapstructure.DecoderConfig{WeaklyTypedInput: true, Result: &bm}
	decoder, err := mapstructure.NewDecoder(&decoderConfig)

	if err != nil {
		log.Fatalf("err: %v", err)
	}

	err = decoder.Decode(stage)
	if err != nil {
		log.Fatalf("err: %v", err)
	}
}

func (bm *BakeManifest) expand(metadata *stage.Stage) {
	// TODO check that index on ExpectedArtifacts is always 0
	expectArtifacts := &bm.ExpectedArtifacts[0]
	// expectArtifacts ID is used by the deploy stage
	expectArtifacts.Id = bm.newUUID(expectArtifacts.DisplayName + bm.Name).String()

	// TODO make sure MatchArtifact ID is not used
	//expectArtifacts.MatchArtifact.Id = bm.newUUID(expectArtifacts.MatchArtifact.Name+expectArtifacts.MatchArtifact.Type).String()

	// TODO check that index on InputArtifacts is always 0
	inputArtifacts := &bm.InputArtifacts[0]
	//Deduplicate ArtifactAccount name
	inputArtifacts.Artifact.ArtifactAccount = inputArtifacts.Account
	//inputArtifacts.Artifact.Id = bm.newUUID(inputArtifacts.Artifact.Name + inputArtifacts.Artifact.Version).String()

	// RefId is either specified by the user or generated based on the stage index
	bm.RefId = metadata.RefId
}

func (bm BakeManifest) newUUID(data string) uuid.UUID {
	// Just a rand root uuid
	namespace, _ := uuid.Parse("e8b764da-5fe5-51ed-8af8-c5c6eca28d7a")
	return uuid.NewSHA1(namespace, []byte(data))
}

func (bm *BakeManifest) updateStage(stage *map[string]interface{}) {
	d := datastore.Datastore{}
	buffer := d.MarshalJSON(bm)
	stageMap := new(map[string]interface{})
	err := json.Unmarshal(buffer, stageMap)
	if err != nil {
		log.Fatalf("Failed to unmarshal JSON:  %v", err)
	}

	*stage = *stageMap
}
