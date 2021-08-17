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
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"swinch/domain/datastore"
	"swinch/domain/stage"
)

type DeleteManifest struct {
	Name                 string   `yaml:"name" json:"name"`
	Type                 string   `yaml:"type,omitempty" json:"type,omitempty"`
	RefId                string   `yaml:"refId,omitempty" json:"refId"`
	RequisiteStageRefIds []string `yaml:"requisiteStageRefIds" json:"requisiteStageRefIds"`

	Account  string `yaml:"account,omitempty" json:"account,omitempty"`
	App      string `yaml:"-" json:"app,omitempty"`
	Location string `yaml:"-" json:"location"`
	// Namespace not in spinnaker json struct
	Namespace          string          `yaml:"namespace,omitempty" json:"-"`
	Kinds              []string        `yaml:"kinds,omitempty" json:"kinds,omitempty"`
	LabelSelectors     *LabelSelectors `yaml:"labelSelectors,omitempty" json:"labelSelectors,omitempty"`
	Options            *Options        `yaml:"options,omitempty" json:"options,omitempty"`
	Mode               string          `yaml:"mode,omitempty" json:"mode,omitempty"`
	CloudProvider      string          `yaml:"cloudProvider,omitempty" json:"cloudProvider,omitempty"`
	ManifestArtifactId *string         `json:"manifestArtifactId,omitempty"`
}

type LabelSelectors struct {
	Selectors []struct {
		Key    string   `yaml:"key" json:"key"`
		Kind   string   `yaml:"kind" json:"kind"`
		Values []string `yaml:"values" json:"values"`
	} `yaml:"selectors" json:"selectors"`
}

type Options struct {
	Cascading bool `yaml:"cascading" json:"cascading"`
}

func (delm *DeleteManifest) ProcessDeleteManifest(p *Pipeline, stage *map[string]interface{}, metadata *stage.Stage) {
	delm.decode(stage)
	delm.expand(p, metadata)
	delm.updateStage(stage)
}

func (delm *DeleteManifest) decode(stage *map[string]interface{}) {
	decoderConfig := mapstructure.DecoderConfig{WeaklyTypedInput: true, Result: &delm}
	decoder, err := mapstructure.NewDecoder(&decoderConfig)
	if err != nil {
		log.Fatalf("err: %v", err)
	}

	err = decoder.Decode(stage)
	if err != nil {
		log.Fatalf("err: %v", err)
	}
}

func (delm DeleteManifest) expand(p *Pipeline, metadata *stage.Stage) {
	delm.App = p.Manifest.Metadata.Application
	if delm.Location != "" {
		delm.Namespace = delm.Location
	} else if delm.Namespace != "" {
		delm.Location = delm.Namespace
	}

	// RefId is either specified by the user or generated based on the stage index
	delm.RefId = metadata.RefId
}

func (delm *DeleteManifest) updateStage(stage *map[string]interface{}) {
	d := datastore.Datastore{}
	buffer := d.MarshalJSON(delm)
	stageMap := new(map[string]interface{})
	err := json.Unmarshal(buffer, stageMap)
	if err != nil {
		log.Fatalf("Failed to unmarshal JSON:  %v", err)
	}

	*stage = *stageMap
}
