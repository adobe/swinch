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
	"bytes"
	"github.com/Masterminds/sprig"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"swinch/domain/datastore"
	"swinch/domain/manifest"
	"text/template"
)

type Template struct {
	Values
	datastore.Datastore
}

func (t *Template) TemplateChart(chartPath, valuesFile, outputPath string, fullRender , excludeDefaultValues bool) {
	values := t.loadValuesFile(chartPath, valuesFile, excludeDefaultValues)
	for _, chartTemplate := range t.discoverTemplates(chartPath) {
		log.Debugf("Found chart template: %v", chartTemplate)

		buffer := t.templateFile(chartPath, chartTemplate.Name(), values)

		if fullRender != false {
			buffer = t.fullRender(buffer)
		}

		t.writeTemplateFile(outputPath, chartTemplate.Name(), buffer)
	}
}

func (t Template) discoverTemplates(chartPath string) []os.DirEntry {
	chartTemplates, err := os.ReadDir(path.Join(chartPath, TemplatesFolder))
	if err != nil {
		log.Fatalf("Error dicovering Chart templates: %v", err)
	}

	return chartTemplates
}

func (t Template) templateFile(chartPath, chartTemplate string, values Values) *bytes.Buffer {
	// Create a named template for each file
	templatePath := path.Join(chartPath, TemplatesFolder, chartTemplate)
	tpl := template.New(chartTemplate).Funcs(template.FuncMap(sprig.FuncMap()))
	tpl, err := tpl.ParseFiles(templatePath)
	if err != nil {
		log.Fatalf("Error in parsing: %v", err)
	}

	buffer := new(bytes.Buffer)
	err = tpl.Execute(buffer, values)
	if err != nil {
		log.Fatalf("Error templating: %v", err)
	}

	return buffer
}

func (t *Template) fullRender(buffer *bytes.Buffer) *bytes.Buffer {
	m := manifest.NewManifest{}
	manifests := m.Decode(buffer)
	buffer.Reset()
	for _, newManifest := range manifests {
		switch newManifest.Kind {
		case m.Application.GetKind():
			buffer.Write(t.MarshalYAML(m.Application.Load(newManifest).Manifest))
		case m.Pipeline.GetKind():
			buffer.Write(t.MarshalYAML(m.Pipeline.Load(newManifest).Manifest))
		}
	}
	return buffer
}

func (t *Template) writeTemplateFile(outputPath, chartTemplate string, buffer *bytes.Buffer) {
	t.Mkdir(path.Join(outputPath), FilePerm)
	t.WriteFile(path.Join(outputPath, chartTemplate), buffer.Bytes(), FilePerm)
}
