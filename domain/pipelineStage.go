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

//func (s *Stage) deleteManifest(p *Pipeline) {
//	s.App = p.Metadata.Application
//	//s.Location = s.Namespace
//}



type TypeManifest struct {
	Namespace string `yaml:"namespace,omitempty" json:"namespace,omitempty"`
}

type ManifestStageToFix struct {
	Overrides                struct{}           `yaml:"overrides,omitempty" json:"overrides,omitempty"`
	RefId                    string             `yaml:"refId,omitempty" json:"refId,omitempty"`
	Account                  string             `yaml:"account,omitempty" json:"account,omitempty"`
	CloudProvider            string             `yaml:"cloudProvider,omitempty" json:"cloudProvider,omitempty"`
	ManifestArtifactId       string             `yaml:"-" json:"manifestArtifactId"`
	NamespaceOverride        string             `yaml:"namespaceOverride,omitempty" json:"namespaceOverride,omitempty"`
	SkipExpressionEvaluation bool               `yaml:"skipExpressionEvaluation,omitempty" json:"skipExpressionEvaluation,omitempty"`
	Source                   string             `yaml:"source,omitempty" json:"source,omitempty"`
	TrafficManagement        *TrafficManagement `yaml:"trafficManagement,omitempty" json:"trafficManagement,omitempty"`
	Kinds                    []string           `yaml:"kinds,omitempty" json:"kinds,omitempty"`
	LabelSelectors           *LabelSelectors    `yaml:"labelSelectors,omitempty" json:"labelSelectors,omitempty"`
	Mode                     string             `yaml:"mode,omitempty" json:"mode,omitempty"`
	Options                  *Options           `yaml:"options,omitempty" json:"options,omitempty"`

}

// TrafficManagement is part of Stages
type TrafficManagement struct {
	Enabled bool `yaml:"enabled" json:"enabled"`
	Options struct {
		EnableTraffic bool          `yaml:"enableTraffic" json:"enableTraffic"`
		Services      []interface{} `yaml:"services" json:"services"`
	} `yaml:"options" json:"options"`
}

// LabelSelectors is part of Stages
type LabelSelectors struct {
	Selectors []struct {
		Key    string   `yaml:"key" json:"key"`
		Kind   string   `yaml:"kind" json:"kind"`
		Values []string `yaml:"values" json:"values"`
	} `yaml:"selectors" json:"selectors"`
}

// Options is part of Stages
type Options struct {
	Cascading bool `yaml:"cascading" json:"cascading"`
}

type DeleteManifest struct {
	TypeManifest `yaml:",inline" json:"-"`
	App          string `yaml:"-" json:"app,omitempty"`
	Location     string `yaml:"-" json:"location,omitempty"`
}

//type ManifestStage struct {
//	ExpectedArtifacts        []ExpectedArtifacts `yaml:"expectedArtifacts,omitempty" json:"expectedArtifacts,omitempty"`
//	InputArtifacts           []InputArtifacts    `yaml:"inputArtifacts,omitempty" json:"inputArtifacts,omitempty"`
//	Namespace                string              `yaml:"namespace,omitempty" json:"namespace,omitempty"`
//	OutputName               string              `yaml:"outputName,omitempty" json:"outputName,omitempty"`
//	Overrides                struct{}            `yaml:"overrides,omitempty" json:"overrides,omitempty"`
//	RefId                    string              `yaml:"refId,omitempty" json:"refId,omitempty"`
//	RequisiteStageRefIds     []string            `yaml:"requisiteStageRefIds" json:"requisiteStageRefIds"`
//	TemplateRenderer         string              `yaml:"templateRenderer,omitempty" json:"templateRenderer,omitempty"`
//	Account                  string              `yaml:"account,omitempty" json:"account,omitempty"`
//	CloudProvider            string              `yaml:"cloudProvider,omitempty" json:"cloudProvider,omitempty"`
//	ManifestArtifactId       string              `yaml:"-" json:"manifestArtifactId"`
//	Moniker                  *Moniker            `yaml:"moniker,omitempty" json:"moniker,omitempty"`
//	NamespaceOverride        string              `yaml:"namespaceOverride,omitempty" json:"namespaceOverride,omitempty"`
//	SkipExpressionEvaluation bool                `yaml:"skipExpressionEvaluation,omitempty" json:"skipExpressionEvaluation,omitempty"`
//	Source                   string              `yaml:"source,omitempty" json:"source,omitempty"`
//	TrafficManagement        *TrafficManagement  `yaml:"trafficManagement,omitempty" json:"trafficManagement,omitempty"`
//	App                      string              `yaml:"-" json:"app,omitempty"`
//	Kinds                    []string            `yaml:"kinds,omitempty" json:"kinds,omitempty"`
//	LabelSelectors           *LabelSelectors     `yaml:"labelSelectors,omitempty" json:"labelSelectors,omitempty"`
//	Location                 string              `yaml:"-" json:"location,omitempty"`
//	Mode                     string              `yaml:"mode,omitempty" json:"mode,omitempty"`
//	Options                  *Options            `yaml:"options,omitempty" json:"options,omitempty"`
//	// swinch extra fields
//	BakeStageRefIds string `yaml:"bakeStageRefIds,omitempty" json:"-"`
//}

//type ManualJudgment struct {
//	FailPipeline   bool          `yaml:"failPipeline" json:"failPipeline"`
//	IsNew          bool          `yaml:"isNew" json:"isNew"`
//	JudgmentInputs []interface{} `yaml:"judgmentInputs" json:"judgmentInputs"`
//	Notifications  []interface{} `yaml:"notifications" json:"notifications"`
//}
