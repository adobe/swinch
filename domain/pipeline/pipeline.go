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
	log "github.com/sirupsen/logrus"
	"swinch/domain/datastore"
	"swinch/domain/util"
	"swinch/spincli"
)

type Pipeline struct {
	Manifest
	Processor
	Spec
	util.Util
	spincli.PipelineAPI
	datastore.Datastore
}

func (p *Pipeline) Plan() {
	p.Apply(true, true)
}

func (p *Pipeline) Apply(dryRun, plan bool) {
	existingPipe := p.Get(p.Metadata.Application, p.Metadata.Name)
	changes := false
	newPipe := false
	if len(existingPipe) == 0 {
		newPipe = true
	} else {
		changes = p.Changes(p.MarshalJSON(p.LoadSpec(existingPipe)), p.MarshalJSON(p.Spec))
	}

	if changes && plan {
		log.Infof("Planing changes for pipeline '%v' in application '%v'", p.Metadata.Name, p.Metadata.Application)
		p.DiffChanges(p.MarshalJSON(p.LoadSpec(existingPipe)), p.MarshalJSON(p.Spec))
	}

	if !dryRun && (changes || newPipe) {
		log.Infof("Saving pipeline '%v' in application '%v'", p.Metadata.Name, p.Metadata.Application)
		p.Save(p.Metadata.Application, p.Metadata.Name, p.WriteJSONTmp(p.Spec))
	}
}

func (p *Pipeline) Destroy() {
	p.Delete(p.Metadata.Name, p.Metadata.Application)
}

// Import TBA
//func (p *Pipeline) importChart() {
//	p.OutputPath = outputPath
//	p.ProtectedImport = protectedImport
//	p.Kind = "pipeline"
//
//	data := new([]byte)
//	if filePath != "" {
//		*data = p.ReadFile(filePath)
//	} else {
//		*data = p.Get()
//	}
//
//	manifest := p.MakeManifest(p.LoadSpec(*data))
//	p.Chart.Metadata.Name = chartName
//	if p.Chart.Metadata.Name == "" {
//		p.Chart.Metadata.Name = manifest.Metadata.Name
//	}
//
//	p.Values.Values = map[interface{}]interface{}{p.Kind: map[string]string{"name": manifest.Metadata.Name}}
//
//	p.GenerateChart(manifest)
//}
