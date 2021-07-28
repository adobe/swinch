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
	"fmt"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"strconv"
)

type Deploy struct {
	Name                 string `yaml:"name" json:"name"`
	Type                 string `yaml:"type,omitempty" json:"type,omitempty"`
	RefId                string `yaml:"refId,omitempty" json:"refId,omitempty"`
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
	BakeStageRefIds string `yaml:"bakeStageRefIds,omitempty" json:"-"`
}

// SwinchDeployExtra add extra fields not present on Spinnaker struct
type SwinchDeployExtra struct {
	BakeStageRefIds string `yaml:"bakeStageRefIds,omitempty" json:"-"`
}

func (d *Deploy) DeployManifest(p *Pipeline) {
	d.Moniker = new(Moniker)
	d.Moniker.App = p.Metadata.Application
	bakeStageIndex := new(int)
	fmt.Println("deploy")
	fmt.Println(d.RequisiteStageRefIds)
	if d.BakeStageRefIds == "" {
		// Presume a deploy stage has the bake stage as the first element in RequisiteStageRefIds
		*bakeStageIndex = d.RequisiteStageRefIds[0]
	} else {
		*bakeStageIndex, _ = strconv.Atoi(d.BakeStageRefIds)
	}

	// Convert from Spinnaker human readable indexing
	*bakeStageIndex -= 1

	bake := new(Bake)
	err := mapstructure.Decode(p.Spec.Stages[*bakeStageIndex], bake)
	if err != nil {
		log.Fatalf("err: %v", err)
	}

	d.ManifestArtifactId = bake.ExpectedArtifacts[0].Id
}

// Moniker is part of Stages
type Moniker struct {
	App string `yaml:"app" json:"app"`
}
