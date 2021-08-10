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

package application

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"strings"
	"swinch/domain/datastore"
)

const (
	Kind = "Application"
	API  = "spinnaker.adobe.com/alpha1"
)

var (
	AppNameLen = errors.New("invalid name length, 4 char min")
)

type Manifest struct {
	ApiVersion string              `yaml:"apiVersion" json:"apiVersion"`
	Kind       string              `yaml:"kind" json:"kind"`
	Metadata   Metadata `yaml:"metadata" json:"metadata"`
	Spec       Spec     `yaml:"spec" json:"spec"`
}

type Metadata struct {
	Name string `yaml:"name" json:"name"`
}

func (am *Manifest) MakeManifest(spec Spec) *Manifest {
	am.ApiVersion = API
	am.Kind = Kind
	am.Metadata.Name = spec.Name
	am.Spec = spec
	return am
}

func (am *Manifest) LoadManifest(manifest interface{}) {
	d := datastore.Datastore{}
	err := yaml.Unmarshal(d.MarshalYAML(manifest), &am)
	if err != nil {
		log.Fatalf("Error LoadManifest: %v", err)
	}

	am.inferFromManifest()

	err = am.Validate()
	if err != nil {
		log.Fatalf("Application manifest validation failed: %v", err)
	}
}

func (am *Manifest) inferFromManifest() {
	// Spinnaker requires lower case application name
	am.Spec.Name = strings.ToLower(am.Metadata.Name)
}

func (am Manifest) Validate() error {
	if len(am.Spec.Name) < 3 {
		return AppNameLen
	}
	return nil
}
