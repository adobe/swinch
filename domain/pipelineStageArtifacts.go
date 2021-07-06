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

// ExpectedArtifacts is part of Stages
type ExpectedArtifacts struct {
	DefaultArtifact    *DefaultArtifact `yaml:"defaultArtifact,omitempty" json:"defaultArtifact,omitempty"`
	DisplayName        string           `yaml:"displayName" json:"displayName"`
	Id                 string           `yaml:"-" json:"id"` // swinch generated
	MatchArtifact      MatchArtifact    `yaml:"matchArtifact,omitempty" json:"matchArtifact,omitempty"`
	UseDefaultArtifact *bool            `yaml:"useDefaultArtifact,omitempty" json:"useDefaultArtifact,omitempty"`
	UsePriorArtifact   *bool            `yaml:"usePriorArtifact,omitempty" json:"usePriorArtifact,omitempty"`
}

//func (e *ExpectedArtifacts) GetID() {
//	e.Id = NewUUID(e.DisplayName).String()
//}

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
