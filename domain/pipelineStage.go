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

import (
	"github.com/google/uuid"
	"strconv"
)

// StageSpec is part of Pipeline
type StageSpec struct {
	ExpectedArtifacts        []ExpectedArtifacts `yaml:"expectedArtifacts,omitempty" json:"expectedArtifacts,omitempty"`
	InputArtifacts           []InputArtifacts    `yaml:"inputArtifacts,omitempty" json:"inputArtifacts,omitempty"`
	Name                     string              `yaml:"name" json:"name"`
	Namespace                string              `yaml:"namespace,omitempty" json:"namespace,omitempty"`
	OutputName               string              `yaml:"outputName,omitempty" json:"outputName,omitempty"`
	Overrides                struct{}            `yaml:"overrides,omitempty" json:"overrides,omitempty"`
	RefId                    string              `yaml:"refId,omitempty" json:"refId,omitempty"`
	RequisiteStageRefIds     []string            `yaml:"requisiteStageRefIds" json:"requisiteStageRefIds"`
	TemplateRenderer         string              `yaml:"templateRenderer,omitempty" json:"templateRenderer,omitempty"`
	Type                     string              `yaml:"type,omitempty" json:"type,omitempty"`
	Account                  string              `yaml:"account,omitempty" json:"account,omitempty"`
	CloudProvider            string              `yaml:"cloudProvider,omitempty" json:"cloudProvider,omitempty"`
	ManifestArtifactId       string              `yaml:"-" json:"manifestArtifactId"`
	Moniker                  *Moniker            `yaml:"moniker,omitempty" json:"moniker,omitempty"`
	NamespaceOverride        string              `yaml:"namespaceOverride,omitempty" json:"namespaceOverride,omitempty"`
	SkipExpressionEvaluation bool                `yaml:"skipExpressionEvaluation,omitempty" json:"skipExpressionEvaluation,omitempty"`
	Source                   string              `yaml:"source,omitempty" json:"source,omitempty"`
	TrafficManagement        *TrafficManagement  `yaml:"trafficManagement,omitempty" json:"trafficManagement,omitempty"`
	App                      string              `yaml:"-" json:"app,omitempty"`
	Kinds                    []string            `yaml:"kinds,omitempty" json:"kinds,omitempty"`
	LabelSelectors           *LabelSelectors     `yaml:"labelSelectors,omitempty" json:"labelSelectors,omitempty"`
	Location                 string              `yaml:"-" json:"location,omitempty"`
	Mode                     string              `yaml:"mode,omitempty" json:"mode,omitempty"`
	Options                  *Options            `yaml:"options,omitempty" json:"options,omitempty"`
	// swinch extra fields
	BakeStageRefIds string `yaml:"bakeStageRefIds,omitempty" json:"-"`
}

func (s *StageSpec) bakeManifest() {
	//TODO check that index on ExpectedArtifacts is always 0
	expectArtifacts := &s.ExpectedArtifacts[0]

	//expectArtifacts ID used in deploy stages
	expectArtifacts.Id = s.newUUID(expectArtifacts.DisplayName + s.RefId).String()

	//TODO check that MatchArtifact ID not used
	//expectArtifacts.MatchArtifact.Id = NewUUID(expectArtifacts.MatchArtifact.Name+expectArtifacts.MatchArtifact.Type).String()

	//TODO check InputArtifacts ID not used
	//TODO check that index on InputArtifacts is always 0
	//Deduplicate ArtifactAccount name
	s.InputArtifacts[0].Artifact.ArtifactAccount = s.InputArtifacts[0].Account
}

func (s *StageSpec) deployManifest(p *Pipeline) {
	s.Moniker = new(Moniker)
	s.Moniker.App = p.Metadata.Application
	bakeStageIndex := new(int)
	if s.BakeStageRefIds == "" {
		// Presume a deploy stage has the bake stage as the first element in RequisiteStageRefIds
		*bakeStageIndex, _ = strconv.Atoi(s.RequisiteStageRefIds[0])
	} else {
		*bakeStageIndex, _ = strconv.Atoi(s.BakeStageRefIds)
	}

	*bakeStageIndex -= 1
	s.ManifestArtifactId = p.Spec.Stages[*bakeStageIndex].ExpectedArtifacts[0].Id
}

func (s *StageSpec) deleteManifest(p *Pipeline) {
	s.App = p.Metadata.Application
	s.Location = s.Namespace
}

func (s *StageSpec) newUUID(data string) uuid.UUID {
	// Just a rand root uuid
	namespace, _ := uuid.Parse("e8b764da-5fe5-51ed-8af8-c5c6eca28d7a")
	return uuid.NewSHA1(namespace, []byte(data))
}
