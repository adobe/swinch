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
	"errors"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"swinch/domain/datastore"
)

const (
	Kind = "Pipeline"
	API  = "spinnaker.adobe.com/alpha1"
)

var (
	PipeNameLen = errors.New("invalid name length, 4 char min")
)

type Manifest struct {
	ApiVersion string   `yaml:"apiVersion" json:"apiVersion"`
	Kind       string   `yaml:"kind" json:"kind"`
	Metadata   Metadata `yaml:"metadata" json:"metadata"`
	Spec       Spec     `yaml:"spec" json:"spec"`
}

type Metadata struct {
	Name        string `yaml:"name" json:"name"`
	Application string `yaml:"application" json:"application"`
}

func (p *Pipeline) GetKind() string {
	return Kind
}

func (p *Pipeline) MakeManifest(spec Spec) *Pipeline {
	p.ApiVersion = API
	p.Kind = Kind
	p.Metadata.Name = spec.Name
	p.Metadata.Application = spec.Application
	p.Spec = spec
	return p
}

func (p *Pipeline) Load(manifest interface{}) *Pipeline {
	p.decode(manifest)
	p.inferFromMetadata()
	p.process(&p.Manifest)

	err := p.Validate()
	if err != nil {
		log.Fatalf("Pipeline manifest validation failed: %v", err)
	}
	return p
}

func (p *Pipeline) decode(manifest interface{}) {
	d := datastore.Datastore{}
	err := yaml.Unmarshal(d.MarshalYAML(manifest), &p)
	if err != nil {
		log.Fatalf("Error Load: %v", err)
	}
}

func (p *Pipeline) inferFromMetadata() {
	p.Spec.Name = p.Metadata.Name
	p.Spec.Application = p.Metadata.Application
}

func (p *Pipeline) Validate() error {
	if len(p.Spec.Name) < 3 {
		return PipeNameLen
	}
	return nil
}
