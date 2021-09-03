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
	"strconv"
	"swinch/domain/datastore"
	"swinch/domain/stage"
)

type RunJobManifest struct {
	Name                 string   `yaml:"name" json:"name"`
	Type                 string   `yaml:"type,omitempty" json:"type,omitempty"`
	RefId                string   `yaml:"refId,omitempty" json:"refId"`
	RequisiteStageRefIds []string `yaml:"requisiteStageRefIds" json:"requisiteStageRefIds"`

	IsNew                 bool   `yaml:"isNew,omitempty" json:"isNew,omitempty"`
	Account               string `yaml:"account" json:"account"`
	Credentials           string `yaml:"credentials" json:"credentials"`
	Alias                 string `yaml:"alias" json:"alias"`
	Application           string `yaml:"application" json:"application"`
	CloudProvider         string `yaml:"cloudProvider" json:"cloudProvider"`
	Source                string `yaml:"source" json:"source"`
	ManifestArtifactId    string `json:"manifestArtifactId"`
	ConsumeArtifactSource string `yaml:"consumeArtifactSource" json:"consumeArtifactSource"`

	ContinuePipeline              bool `yaml:"continuePipeline,omitempty" json:"continuePipeline,omitempty"`
	FailPipeline                  bool `yaml:"failPipeline,omitempty" json:"failPipeline,omitempty"`
	CompleteOtherBranchesThenFail bool `yaml:"completeOtherBranchesThenFail,omitempty" json:"completeOtherBranchesThenFail,omitempty"`

	// Swinch only field
	JobBakeStageRefIds *int `yaml:"bakeStageRefIds,omitempty" json:"-"`
}

func (rjm RunJobManifest) ProcessRunJobManifest(p *Pipeline, stageMap *map[string]interface{}, metadata *stage.Stage) {
	rjm.decode(stageMap)
	rjm.expand(p, metadata)
	rjm.update(stageMap)
}

func (rjm *RunJobManifest) decode(stageMap *map[string]interface{}) {
	decoderConfig := mapstructure.DecoderConfig{WeaklyTypedInput: true, Result: &rjm}
	decoder, err := mapstructure.NewDecoder(&decoderConfig)
	if err != nil {
		log.Fatalf("err: %v", err)
	}

	err = decoder.Decode(stageMap)
	if err != nil {
		log.Fatalf("err: %v", err)
	}
}

func (rjm *RunJobManifest) expand(p *Pipeline, metadata *stage.Stage) {
	bakeStageIndex := new(int)

	// Bind job stage to a specific bake
	if rjm.JobBakeStageRefIds == nil {
		// Presume a job stage has the bake stage as the first element in RequisiteStageRefIds
		*bakeStageIndex, _ = strconv.Atoi(rjm.RequisiteStageRefIds[0])
	} else {
		*bakeStageIndex = *rjm.JobBakeStageRefIds
	}

	// Convert from Spinnaker human-readable indexing
	*bakeStageIndex -= 1

	//TODO get the bake stage without decoding
	bake := new(BakeManifest)
	err := mapstructure.Decode(p.Manifest.Spec.Stages[*bakeStageIndex], bake)
	if err != nil {
		log.Fatalf("err: %v", err)
	}

	rjm.ManifestArtifactId = bake.ExpectedArtifacts[0].Id

	// RefId is either specified by the user or generated based on the stage index
	rjm.RefId = metadata.RefId
}

func (rjm *RunJobManifest) update(stageMap *map[string]interface{}) {
	d := datastore.Datastore{}
	buffer := new(map[string]interface{})
	err := json.Unmarshal(d.MarshalJSON(rjm), buffer)
	if err != nil {
		log.Fatalf("Failed to unmarshal JSON:  %v", err)
	}
	*stageMap = *buffer
}
