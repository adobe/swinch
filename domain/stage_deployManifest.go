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

package domain

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
)

type DeployManifest struct {
	Name                 string `yaml:"name" json:"name"`
	Type                 string `yaml:"type,omitempty" json:"type,omitempty"`
	RefId                string `yaml:"refId,omitempty" json:"refId"`
	RequisiteStageRefIds []int  `yaml:"requisiteStageRefIds" json:"requisiteStageRefIds"`

	Account            string              `json:"account"`
	ExpectedArtifacts  []ExpectedArtifacts `yaml:"expectedArtifacts,omitempty" json:"expectedArtifacts,omitempty"`
	CloudProvider      string              `json:"cloudProvider"`
	ManifestArtifactId string              `json:"manifestArtifactId"`
	Moniker            *Moniker            `yaml:"moniker,omitempty" json:"moniker,omitempty"`
	NamespaceOverride  string              `json:"namespaceOverride"`
	Overrides          struct{}            `json:"overrides"`
	Source             string              `json:"source"`

	//SwinchDeployExtra
	BakeStageRefIds *int `yaml:"bakeStageRefIds,omitempty" json:"-"`
}

// Moniker is part of Stages
type Moniker struct {
	App string `yaml:"app" json:"app"`
}

func (dm *DeployManifest) ProcessDeployManifest(p *Pipeline, stage *map[string]interface{}, metadata *StageMetadata) {
	dm.decode(stage)
	dm.expand(p, metadata)
	dm.updateStage(stage)
}

func (dm *DeployManifest) expand(p *Pipeline, metadata *StageMetadata) {
	fmt.Println(dm.RefId)
	dm.Moniker = new(Moniker)
	dm.Moniker.App = p.Metadata.Application
	bakeStageIndex := new(int)

	// Bind deploy stage to a specific bake
	if dm.BakeStageRefIds == nil {
		// Presume a deploy stage has the bake stage as the first element in RequisiteStageRefIds
		*bakeStageIndex = dm.RequisiteStageRefIds[0]
	} else {
		*bakeStageIndex = *dm.BakeStageRefIds
	}

	// Convert from Spinnaker human readable indexing
	*bakeStageIndex -= 1

	//TODO get the bake stage without decoding
	bake := new(DeployManifest)
	err := mapstructure.Decode(p.Spec.Stages[*bakeStageIndex], bake)
	if err != nil {
		log.Fatalf("err: %v", err)
	}

	dm.ManifestArtifactId = bake.ExpectedArtifacts[0].Id

	//RefId is either specified by the user or generated based on the stage index
	dm.RefId = metadata.RefId
}

func (dm *DeployManifest) updateStage(stage *map[string]interface{}) {
	d := Datastore{}
	buffer := d.MarshalJSON(dm)
	stageMap := new(map[string]interface{})
	err := json.Unmarshal(buffer, stageMap)
	if err != nil {
		log.Fatalf("Failed to unmarshal JSON:  %v", err)

	}

	*stage = *stageMap
}

func (dm *DeployManifest) decode(stage *map[string]interface{}) {
	err := mapstructure.Decode(stage, *dm)
	if err != nil {
		log.Fatalf("err: %v", err)
	}
}
