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
	"github.com/spf13/cobra"
	"swinch/domain/application"
	"swinch/domain/manifest"
	"swinch/domain/pipeline"
	"swinch/domain/stages"
)

// planCmd represents the plan command
var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Plan",
	Long:  `Plan`,
	PreRun: func(cmd *cobra.Command, args []string) {
		SetLogLevel(logLevel)
		ValidateConfigFile()
		ValidateConfig()
	},
	Run: func(cmd *cobra.Command, args []string) {
		Plan()
	},
}

func init() {
	planCmd.Flags().StringVarP(&filePath, "file", "f", "", "Manifest file or directory, non recursive")
	planCmd.MarkFlagRequired("file")
	rootCmd.AddCommand(planCmd)
}

func Plan() {
	m := manifest.Manifest{}
	a := application.Application{}
	p := pipeline.Pipeline{}

	manifests := m.GetManifests(filePath)
	for _, manifest := range manifests {
		switch manifest.Kind {
		case a.Manifest.Kind:
			a.LoadManifest(manifest)
			a.Plan()
		case p.GetKind():
			p.LoadManifest(manifest)
			s := stages.Processor{}
			s.Process(&p.Manifest)
			p.Plan()
		}
	}
}
