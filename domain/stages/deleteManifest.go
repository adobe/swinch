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
)

const deleteManifest = "deleteManifest"

type DeleteManifest struct {
	Metadata `mapstructure:",squash"`

	Account  string `yaml:"account,omitempty" json:"account,omitempty"`
	App      string `yaml:"-" json:"app,omitempty"`
	Location string `yaml:"-" json:"location,omitempty"`
	// Namespace not in spinnaker json struct
	Namespace          string          `yaml:"namespace,omitempty" json:"-"`
	Kinds              []string        `yaml:"kinds,omitempty" json:"kinds,omitempty"`
	LabelSelectors     *LabelSelectors `yaml:"labelSelectors,omitempty" json:"labelSelectors,omitempty"`
	Options            *Options        `yaml:"options,omitempty" json:"options,omitempty"`
	Mode               string          `yaml:"mode,omitempty" json:"mode,omitempty"`
	CloudProvider      string          `yaml:"cloudProvider,omitempty" json:"cloudProvider,omitempty"`
	ManifestArtifactId *string         `json:"manifestArtifactId,omitempty"`

	ContinuePipeline              bool `yaml:"continuePipeline,omitempty" json:"continuePipeline,omitempty"`
	FailPipeline                  bool `yaml:"failPipeline,omitempty" json:"failPipeline,omitempty"`
	CompleteOtherBranchesThenFail bool `yaml:"completeOtherBranchesThenFail,omitempty" json:"completeOtherBranchesThenFail,omitempty"`
}

type LabelSelectors struct {
	Selectors []struct {
		Key    string   `yaml:"key" json:"key"`
		Kind   string   `yaml:"kind" json:"kind"`
		Values []string `yaml:"values" json:"values"`
	} `yaml:"selectors" json:"selectors"`
}

type Options struct {
	Cascading          bool `yaml:"cascading" json:"cascading"`
	GracePeriodSeconds int  `yaml:"gracePeriodSeconds" json:"gracePeriodSeconds"`
}

func (delm DeleteManifest) GetStageType() string {
	return deleteManifest
}

func (delm DeleteManifest) Process(stage *Stage) {
	delm.decode(stage)
	delm.expand(stage)
	delm.update(stage)
}

func (delm *DeleteManifest) decode(stage *Stage) {
	decoderConfig := mapstructure.DecoderConfig{WeaklyTypedInput: true, Result: &delm}
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

func (delm *DeleteManifest) expand(stage *Stage) {
	delm.App = stage.ManifestMetadata.Application
	if delm.Location != "" {
		delm.Namespace = delm.Location
	} else if delm.Namespace != "" {
		delm.Location = delm.Namespace
	}
}

func (delm *DeleteManifest) update(stage *Stage) {
	d := datastore.Datastore{}
	tmpStage := new(map[string]interface{})
	err := json.Unmarshal(d.MarshalJSON(delm), tmpStage)
	if err != nil {
		log.Fatalf("Failed to unmarshal JSON:  %v", err)
	}
	*stage.RawStage = *tmpStage
}
