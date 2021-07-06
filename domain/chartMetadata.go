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
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"path"
)

var DefaultChartMetadata = `
apiVersion: v1
description: "This is a generated chart"
version: 0.0.1
`

type ChartMetadata struct {
	ApiVersion  string `yaml:"apiVersion" json:"apiVersion"`
	Description string `yaml:"description" json:"description"`
	Name        string `yaml:"name" json:"name"`
	Version     string `yaml:"version" json:"version"`
}

func (m ChartMetadata) loadMetadataFile(ChartPath string) ChartMetadata {
	d := Datastore{}
	metadataBuffer := d.ReadFile(path.Join(ChartPath, "/Chart.yaml"))
	return m.loadMetadata(metadataBuffer)
}

func (m ChartMetadata) loadMetadata(byteData []byte) ChartMetadata {
	err := yaml.Unmarshal(byteData, &m)
	if err != nil {
		log.Fatalf("Error loading Chart metadata: %v", err)
	}
	return m
}
