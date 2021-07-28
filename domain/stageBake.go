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

package domain

import "github.com/google/uuid"

func (b *Bake) bakeManifest() {
	//TODO check that index on ExpectedArtifacts is always 0
	expectArtifacts := &b.ExpectedArtifacts[0]

	//expectArtifacts ID used in deploy stages
	expectArtifacts.Id = b.newUUID(expectArtifacts.DisplayName + b.Name).String()
	//expectArtifacts.Id = b.newUUID(expectArtifacts.DisplayName + b.RefId).String()

	//TODO check that MatchArtifact ID not used
	//expectArtifacts.MatchArtifact.Id = NewUUID(expectArtifacts.MatchArtifact.Name+expectArtifacts.MatchArtifact.Type).String()

	//TODO check InputArtifacts ID not used
	//TODO check that index on InputArtifacts is always 0
	//Deduplicate ArtifactAccount name
	b.InputArtifacts[0].Artifact.ArtifactAccount = b.InputArtifacts[0].Account
}

func (b *Bake) newUUID(data string) uuid.UUID {
	// Just a rand root uuid
	namespace, _ := uuid.Parse("e8b764da-5fe5-51ed-8af8-c5c6eca28d7a")
	return uuid.NewSHA1(namespace, []byte(data))
}

//Stage is part of Pipeline
type Stage struct {
	Name                 string `yaml:"name" json:"name"`
	Type                 string `yaml:"type,omitempty" json:"type,omitempty"`
	RefId                string `yaml:"refId,omitempty" json:"refId,omitempty"`
	RequisiteStageRefIds []int  `yaml:"requisiteStageRefIds" json:"requisiteStageRefIds"`
}

type Bake struct {
	Name                 string `yaml:"name" json:"name"`
	Type                 string `yaml:"type,omitempty" json:"type,omitempty"`
	RefId                string `yaml:"refId,omitempty" json:"refId,omitempty"`
	RequisiteStageRefIds []int  `yaml:"requisiteStageRefIds" json:"requisiteStageRefIds"`

	OutputName         string              `json:"outputName"`
	ExpectedArtifacts  []ExpectedArtifacts `yaml:"expectedArtifacts,omitempty" json:"expectedArtifacts,omitempty"`
	InputArtifacts     []InputArtifacts    `yaml:"inputArtifacts,omitempty" json:"inputArtifacts,omitempty"`
	ManifestArtifactId string              `json:"manifestArtifactId"`
	Namespace          string              `json:"namespace"`
	TemplateRenderer   string              `json:"templateRenderer"`
	Overrides          struct{}
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
		ArtifactAccount string  `yaml:"artifactAccount" json:"artifactAccount"`
		Id              *string `yaml:"-" json:"id,omitempty"`
		Name            string  `yaml:"name" json:"name"`
		Type            string  `yaml:"type" json:"type"`
		Version         string  `yaml:"version" json:"version"`
	} `yaml:"artifact" json:"artifact"`
}
