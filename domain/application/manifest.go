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
	ApiVersion string   `yaml:"apiVersion" json:"apiVersion"`
	Kind       string   `yaml:"kind" json:"kind"`
	Metadata   Metadata `yaml:"metadata" json:"metadata"`
	Spec       Spec     `yaml:"spec" json:"spec"`
}

type Metadata struct {
	Name string `yaml:"name" json:"name"`
}

func (a *Application) GetKind() string {
	return Kind
}

func (a *Application) MakeManifest(spec Spec) *Application {
	a.ApiVersion = API
	a.Kind = Kind
	a.Metadata.Name = spec.Name
	a.Spec = spec
	return a
}

func (a *Application) Load(manifest interface{}) *Application {
	a.decode(manifest)
	a.inferFromManifest()

	err := a.validate()
	if err != nil {
		log.Fatalf("Application manifest validation failed: %v", err)
	}
	return a
}

func (a *Application) decode(manifest interface{}) {
	d := datastore.Datastore{}
	err := yaml.Unmarshal(d.MarshalYAML(manifest), &a)
	if err != nil {
		log.Fatalf("Error Load: %v", err)
	}
}

func (a *Application) inferFromManifest() {
	// Spinnaker requires lower case application name
	a.Spec.Name = strings.ToLower(a.Metadata.Name)
}

func (a *Application) validate() error {
	if len(a.Spec.Name) < 3 {
		return AppNameLen
	}
	return nil
}
