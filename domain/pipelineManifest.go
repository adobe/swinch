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
	"errors"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

const (
	PipelineKind = "Pipeline"
	PipelineAPI  = "spinnaker.adobe.com/alpha1"
)

var (
	PipeNameLen = errors.New("invalid name length, 4 char min")
)

type PipelineManifest struct {
	ApiVersion string           `yaml:"apiVersion" json:"apiVersion"`
	Kind       string           `yaml:"kind" json:"kind"`
	Metadata   PipelineMetadata `yaml:"metadata" json:"metadata"`
	Spec       PipelineSpec     `yaml:"spec" json:"spec"`
}

type PipelineMetadata struct {
	Name        string `yaml:"name" json:"name"`
	Application string `yaml:"application" json:"application"`
}

func (pm *PipelineManifest) LoadManifest(manifest interface{}) {
	d := Datastore{}
	err := yaml.Unmarshal(d.marshalYAML(manifest), &pm)
	if err != nil {
		log.Fatalf("Error LoadManifest: %v", err)
	}

	pm.inferFromMetadata()

	err = pm.Validate()
	if err != nil {
		log.Fatalf("Pipeline manifest validation failed: %v", err)
	}
}

func (pm *PipelineManifest) inferFromMetadata() {
	pm.Spec.Name = pm.Metadata.Name
	pm.Spec.Application = pm.Metadata.Application
}

func (pm PipelineManifest) Validate() error {
	if len(pm.Spec.Name) < 3 {
		return PipeNameLen
	}
	return nil
}
