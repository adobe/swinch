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

package manifest

import (
	"bytes"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io"
	"swinch/domain/application"
	"swinch/domain/datastore"
	"swinch/domain/pipeline"
)

var Kinds = map[string]string{
	application.Kind: application.API,
	pipeline.Kind:    pipeline.API,
}

type Manifest struct {
	ApiVersion string      `yaml:"apiVersion" json:"apiVersion"`
	Kind       string      `yaml:"kind" json:"kind"`
	Metadata   interface{} `yaml:"metadata" json:"metadata"`
	Spec       interface{} `yaml:"spec" json:"spec"`
}


func (m *Manifest) GetManifests(filePath string) ([]application.Manifest, []pipeline.Manifest) {
	d := datastore.Datastore{}
	discoveredYAMLDocs := d.DiscoverYAMLFiles(filePath)
	return m.ReadManifest(discoveredYAMLDocs)
}

// ReadManifest used to load all m types
func (m *Manifest) ReadManifest(yamlFilesBuffer *bytes.Buffer) ([]application.Manifest, []pipeline.Manifest) {
	decoder := yaml.NewDecoder(yamlFilesBuffer)

	applications := make([]application.Manifest, 0)
	pipelines := make([]pipeline.Manifest, 0)

	for {
		manifest := new(Manifest)
		errDecode := decoder.Decode(&manifest)

		if manifest == nil {
			continue
		}
		if errors.Is(errDecode, io.EOF) {
			break
		}
		if errDecode != nil {
			log.Fatalf("Error reading YAML: %v", errDecode)
		}

		err := manifest.Validate()
		if err != nil {
			log.Fatal(err)
		}

		switch manifest.Kind {
		case application.Kind:
			app := application.Application{}
			app.LoadManifest(manifest)
			applications = append(applications, app.Manifest)
		case pipeline.Kind:
			pipe := pipeline.Pipeline{}
			pipe.LoadManifest(manifest)
			pipe.ProcessStages()
			pipelines = append(pipelines, pipe.Manifest)
		default:
			log.Fatalf("Error detecting manifest Kind")
		}
	}
	return applications, pipelines
}

func (m *Manifest) Validate() error {
	_, ok := Kinds[m.Kind]
	if ok {
		kindApiVersion := Kinds[m.Kind]
		if m.ApiVersion != kindApiVersion {
			return fmt.Errorf("bad api version, expected: %v, got: %v", kindApiVersion, m.ApiVersion)
		} else {
			return nil
		}
	} else {
		return fmt.Errorf("unknown manifest Kind: %v", m.Kind)
	}
}
