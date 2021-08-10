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
	"github.com/Masterminds/sprig"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"swinch/domain/datastore"
	"text/template"
)

const (
	FilePerm = 0775
)

type Template struct {
	chartPath string
	Values
}

func (t Template) RenderChart(chartPath, valuesFile, outputPath string) {
	t.chartPath = chartPath
	if valuesFile == "" {
		valuesFile = path.Join(t.chartPath, "/values.yaml")
	}
	values := t.loadValuesFile(valuesFile)

	for _, chartFile := range t.discoverTemplates() {
		log.Debugf("Found chart template: %v", chartFile.Name())

		templatePath := path.Join(t.chartPath, TemplatesFolder, chartFile.Name())
		tpl := template.New(chartFile.Name()).Funcs(template.FuncMap(sprig.FuncMap()))
		tpl, err := tpl.ParseFiles(templatePath)
		if err != nil {
			log.Fatalf("Error in parsing: %v", err)
		}

		d := datastore.Datastore{}
		d.Mkdir(path.Join(outputPath), FilePerm)
		outFile, err := os.Create(path.Join(outputPath, chartFile.Name()))
		if err != nil {
			log.Fatalf("Error in output files: %v", err)
		}

		err = tpl.Execute(outFile, values)
		if err != nil {
			log.Fatalf("Error templating: %v", err)
		}

		err = outFile.Close()
		if err != nil {
			log.Fatalf("File error: %v", err)
		}
	}
}

func (t Template) discoverTemplates() []os.DirEntry {
	chartTemplates, err := os.ReadDir(path.Join(t.chartPath, TemplatesFolder))

	if err != nil {
		log.Fatalf("Error dicovering Chart templates: %v", err)
	}

	return chartTemplates
}
