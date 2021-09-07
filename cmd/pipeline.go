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
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"swinch/domain/pipeline"
)

// pipelineCmd represents the pipeline command
var pipelineCmd = &cobra.Command{
	Use:   "pipeline",
	Short: "Import Spinnaker pipelines",
	Long:  `Import Spinnaker pipelines`,
	PreRun: func(cmd *cobra.Command, args []string) {
		SetLogLevel(logLevel)
		ValidateConfigFile()
		ValidateConfig()
	},
	Run: func(cmd *cobra.Command, args []string) {
		action := cmd.Parent().Use
		cmdPipeAction(applicationName, pipelineName, action)
	},
}

var PlanPipeCmd = *pipelineCmd
var ImportPipeCmd = *pipelineCmd
var DeletePipeCmd = *pipelineCmd

func init() {
	// import flags
	ImportPipeCmd.Flags().StringVarP(&applicationName, "application", "a", "", "Application name")
	ImportPipeCmd.Flags().StringVarP(&pipelineName, "pipeline", "p", "", "Pipeline name")
	ImportPipeCmd.Flags().StringVarP(&filePath, "file", "f", "", "JSON file input")
	ImportPipeCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Generated chart output path")
	ImportPipeCmd.Flags().StringVarP(&chartName, "chart", "n", "", "Specify chart name for imported pipeline")
	ImportPipeCmd.Flags().BoolVarP(&protectedImport, "protected-import", "", false, "Protect already created chart from overwriting")
	ImportPipeCmd.MarkFlagRequired("application")
	ImportPipeCmd.MarkFlagRequired("pipeline")
	ImportPipeCmd.MarkFlagRequired("output")
	importCmd.AddCommand(&ImportPipeCmd)

	// delete flags
	DeletePipeCmd.Flags().StringVarP(&applicationName, "application", "a", "", "Application name")
	DeletePipeCmd.Flags().StringVarP(&pipelineName, "pipeline", "p", "", "Pipeline name")
	DeletePipeCmd.MarkFlagRequired("application")
	DeletePipeCmd.MarkFlagRequired("pipeline")
	deleteCmd.AddCommand(&DeletePipeCmd)

	// plan flags
	planCmd.AddCommand(&PlanPipeCmd)
}

func cmdPipeAction(app, pipe, action string) {
	p := pipeline.Pipeline{}
	switch action {
	case deleteAction:
		p.Delete(app, pipe)
	case importAction:
		fmt.Println("Import TBA")
	default:
		log.Fatalf("Bad application action")
	}
}
