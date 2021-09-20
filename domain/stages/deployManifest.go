/*
Copyright 2021 Adobe. All rights reservedm.
This file is licensed to you under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License. You may obtain a copy
of the License at http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under
the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR REPRESENTATIONS
OF ANY KIND, either express or impliedm. See the License for the specific language
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

const deployManifest = "deployManifest"

type DeployManifest struct {
	Metadata                 `mapstructure:",squash"`
	Account                  string   `yaml:"account,omitempty" json:"account,omitempty"`
	CloudProvider            string   `json:"cloudProvider"`
	ManifestArtifactId       string   `json:"manifestArtifactId"`
	Moniker                  *Moniker `yaml:"moniker,omitempty" json:"moniker"`
	NamespaceOverride        string   `json:"namespaceOverride"`
	Overrides                struct{} `yaml:"overrides,omitempty" json:"overrides,omitempty"`
	Source                   string   `json:"source"`
	SkipExpressionEvaluation bool     `yaml:"skipExpressionEvaluation,omitempty" json:"skipExpressionEvaluation,omitempty"`

	BakeStageRefIds *int `yaml:"bakeStageRefIds,omitempty" json:"-"`
}

// Moniker is part of Stages
type Moniker struct {
	App string `yaml:"app" json:"app"`
}

func (dm DeployManifest) GetStageType() string {
	return deployManifest
}

func (dm DeployManifest) Process(stage *Stage) {
	dm.decode(stage)
	dm.expand(stage)
	dm.update(stage)
}

func (dm *DeployManifest) decode(stage *Stage) {
	decoderConfig := mapstructure.DecoderConfig{WeaklyTypedInput: true, Result: &dm}
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

func (dm *DeployManifest) expand(stage *Stage) {
	dm.Moniker = new(Moniker)
	dm.Moniker.App = stage.ManifestMetadata.Application

	bakeIndex := dm.getBakeIndex()
	bake := new(BakeManifest)
	err := mapstructure.Decode((*stage.Stages)[bakeIndex], bake)
	if err != nil {
		log.Fatalf("err: %v", err)
	}

	dm.ManifestArtifactId = bake.ExpectedArtifacts[0].Id
}

func (dm *DeployManifest) getBakeIndex() int {
	bakeStageIndex := new(int)
	// Bind deploy stage to a specific bake
	if dm.BakeStageRefIds == nil {
		// Presume a deploy stage has the bake stage as the first element in RequisiteStageRefIds
		*bakeStageIndex, _ = strconv.Atoi(dm.RequisiteStageRefIds[0])
	} else {
		*bakeStageIndex = *dm.BakeStageRefIds
	}
	// Convert from Spinnaker human readable indexing
	*bakeStageIndex -= 1

	return *bakeStageIndex
}

func (dm *DeployManifest) update(stage *Stage) {
	d := datastore.Datastore{}
	tmpStage := new(map[string]interface{})
	err := json.Unmarshal(d.MarshalJSON(dm), tmpStage)
	if err != nil {
		log.Fatalf("Failed to unmarshal JSON:  %v", err)
	}
	*stage.RawStage = *tmpStage
}
