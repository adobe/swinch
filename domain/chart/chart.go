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

package chart

import (
	log "github.com/sirupsen/logrus"
	"path"
	"swinch/domain/datastore"
)

const (
	ChartTemplatesFolder = "templates"
	ChartValuesFile      = "values.yaml"
	ChartMetadataFile    = "Chart.yaml"
)

type Chart struct {
	OutputPath      string
	Kind            string
	ProtectedImport bool
	ChartMetadata
	ChartValues
	datastore.Datastore
}

// Import

func (c Chart) GenerateChart(manifest interface{}) {
	if c.FileExists(path.Join(c.OutputPath, c.ChartMetadata.Name, ChartValuesFile)) && c.ProtectedImport {
		log.Fatalf("Cannot import over an existing chart, values file present in path '%s'", path.Join(c.OutputPath, c.ChartMetadata.Name, ChartValuesFile))
	} else {
		c.Mkdir(path.Join(c.OutputPath, c.ChartMetadata.Name, "/", ChartTemplatesFolder), FilePerm)
		c.WriteChartMetadata()
		c.WriteChartValues()
		c.WriteManifest(manifest)
	}
}

// WriteChartMetadata default Chart metadata for imported pipelines
func (c Chart) WriteChartMetadata() {
	c.ChartMetadata = c.loadMetadata([]byte(DefaultChartMetadata))
	c.WriteYAML(c.ChartMetadata, path.Join(c.OutputPath, c.ChartMetadata.Name, ChartMetadataFile))
}

// WriteChartValues default Chart values for imported pipelines
func (c Chart) WriteChartValues() {
	c.WriteYAML(c.ChartValues.Values, path.Join(c.OutputPath, c.ChartMetadata.Name, ChartValuesFile))
}

func (c Chart) WriteManifest(manifest interface{}) {
	c.WriteYAML(manifest, path.Join(c.OutputPath, c.ChartMetadata.Name, "/", ChartTemplatesFolder, c.Kind+".yaml"))
}
