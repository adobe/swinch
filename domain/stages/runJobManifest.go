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
	"strconv"
	"swinch/domain/datastore"
)

const runJobManifest = "runJobManifest"

type RunJobManifest struct {
	Stage `mapstructure:",squash"`

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
	JobBakeStageRefIds            *int `yaml:"jobBakeStageRefIds,omitempty" json:"-"`
}

func (rjm RunJobManifest) GetStageType() string {
	return runJobManifest
}

func (rjm RunJobManifest) Process(stage *Stage) {
	rjm.decode(stage)
	rjm.expand(stage)
	rjm.update(stage)
}

func (rjm *RunJobManifest) decode(stage *Stage) {
	decoderConfig := mapstructure.DecoderConfig{WeaklyTypedInput: true, Result: &rjm}
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

func (rjm *RunJobManifest) expand(stage *Stage) {
	bakeIndex := rjm.getBakeIndex()
	//TODO get the bake stage without decoding
	bake := new(BakeManifest)
	err := mapstructure.Decode((*stage.Stages)[bakeIndex], bake)
	if err != nil {
		log.Fatalf("err: %v", err)
	}

	rjm.ManifestArtifactId = bake.ExpectedArtifacts[0].Id
}

func (rjm *RunJobManifest) getBakeIndex() int {
	bakeStageIndex := new(int)
	// Bind deploy stage to a specific bake
	if rjm.JobBakeStageRefIds == nil {
		// Presume a deploy stage has the bake stage as the first element in RequisiteStageRefIds
		*bakeStageIndex, _ = strconv.Atoi(rjm.RequisiteStageRefIds[0])
	} else {
		*bakeStageIndex = *rjm.JobBakeStageRefIds
	}
	// Convert from Spinnaker human readable indexing
	*bakeStageIndex -= 1

	return *bakeStageIndex
}

func (rjm *RunJobManifest) update(stage *Stage) {
	d := datastore.Datastore{}
	tmpStage := new(map[string]interface{})
	err := json.Unmarshal(d.MarshalJSON(rjm), tmpStage)
	if err != nil {
		log.Fatalf("Failed to unmarshal JSON:  %v", err)
	}
	*stage.RawStage = *tmpStage
}
