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

package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"swinch/domain"
	"swinch/spincli"
)

// pipelineCmd represents the pipeline command
var pipelineCmd = &cobra.Command{
	Use:   "pipeline",
	Short: "Import Spinnaker pipelines",
	Long:  `Import Spinnaker pipelines`,
	PreRun: func(cmd *cobra.Command, args []string) {
		SetLogLevel(logLevel)
	},
	Run: func(cmd *cobra.Command, args []string) {
		action := cmd.Parent().Use
		Pipeline{}.cmdActions(application, pipeline, action)
	},
}

var PlanPipeCmd = *pipelineCmd
var ImportPipeCmd = *pipelineCmd
var DeletePipeCmd = *pipelineCmd

func init() {
	// import flags
	ImportPipeCmd.Flags().StringVarP(&application, "application", "a", "", "Application name")
	ImportPipeCmd.Flags().StringVarP(&pipeline, "pipeline", "p", "", "Pipeline name")
	ImportPipeCmd.Flags().StringVarP(&filePath, "filePath", "f", "", "JSON file input")
	ImportPipeCmd.Flags().StringVarP(&outputPath, "outputPath", "o", "", "Generated chart output path")
	ImportPipeCmd.Flags().StringVarP(&chartName, "chartName", "n", "", "Specify chart name for imported pipeline")
	ImportAppCmd.Flags().BoolVarP(&protectedImport, "protectedImport", "n", false, "Protect already created chart from overwriting")
	ImportPipeCmd.MarkFlagRequired("application")
	ImportPipeCmd.MarkFlagRequired("pipeline")
	ImportPipeCmd.MarkFlagRequired("outputPath")
	importCmd.AddCommand(&ImportPipeCmd)

	// delete flags
	DeletePipeCmd.Flags().StringVarP(&application, "application", "a", "", "Application name")
	DeletePipeCmd.Flags().StringVarP(&pipeline, "pipeline", "p", "", "Pipeline name")
	DeletePipeCmd.MarkFlagRequired("application")
	DeletePipeCmd.MarkFlagRequired("pipeline")
	deleteCmd.AddCommand(&DeletePipeCmd)

	// plan flags
	planCmd.AddCommand(&PlanPipeCmd)
}

type Pipeline struct {
	manifests []domain.PipelineManifest
	domain.Pipeline
	spincli.PipelineAPI
	domain.Chart
}

func (p Pipeline) cmdActions(app, pipe, action string) {
	p.App = app
	p.Pipe = pipe
	switch action {
	case deleteAction:
		p.Delete()
	case importAction:
		p.importChart()
	default:
		log.Fatalf("Bad application action")
	}
}

func (p Pipeline) manifestActions(action string) {
	for i := 0; i < len(p.manifests); i++ {
		manifest := &p.manifests[i]
		p.App = manifest.Metadata.Application
		p.Pipe = manifest.Metadata.Name
		switch action {
		case applyAction:
			dryRun := false
			p.save(manifest.Spec, dryRun)
		case deleteAction:
			p.Delete()
		case planAction:
			dryRun := true
			p.save(manifest.Spec, dryRun)
		default:
			log.Fatalf("Bad application action")
		}
	}
}

func (p *Pipeline) save(spec domain.PipelineSpec, dryRun bool) {
	pipe := p.Get()
	changes := false
	newPipe := false
	if len(pipe) == 0 {
		newPipe = true
	} else {
		changes = Changes(p.MarshalJSON(p.LoadSpec(pipe)), p.MarshalJSON(spec))
	}

	if changes && plan {
		log.Infof("Planing changes for pipeline '%v' in application '%v'", p.Pipe, p.App)
		DiffChanges(p.MarshalJSON(p.LoadSpec(pipe)), p.MarshalJSON(spec))
	}

	if !dryRun && (changes || newPipe) {
		log.Infof("Saving pipeline '%v' in application '%v'", p.Pipe, p.App)
		p.Save(p.WriteJSONTmp(spec))
	}
}

func (p *Pipeline) importChart() {
	p.OutputPath = outputPath
	p.ProtectedImport = protectedImport
	p.Kind = "pipeline"

	data := new([]byte)
	if filePath != "" {
		*data = p.ReadFile(filePath)
	} else {
		*data = p.Get()
	}

	manifest := p.MakePipelineManifest(p.LoadSpec(*data))
	p.ChartMetadata.Name = chartName
	if p.ChartMetadata.Name == "" {
		p.ChartMetadata.Name = manifest.Metadata.Name
	}

	p.ChartValues.Values = map[interface{}]interface{}{p.Kind: map[string]string{"name": manifest.Metadata.Name}}

	p.GenerateChart(manifest)
}
