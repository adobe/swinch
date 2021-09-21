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
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply or sync an Application or Pipeline from a manifest",
	Long:  "Apply or sync an Application or Pipeline from a manifest",
	PreRun: func(cmd *cobra.Command, args []string) {
		SetLogLevel(logLevel)
		ValidateConfigFile()
		ValidateConfig()
	},
	Run: func(cmd *cobra.Command, args []string) {
		Apply()
	},
}

func init() {
	applyCmd.Flags().StringVarP(&filePath, "file", "f", "", "Manifest file or directory, non recursive")
	applyCmd.Flags().BoolVarP(&plan, "plan", "p", true, "Display plan before apply, no user input.")
	applyCmd.MarkFlagRequired("file")
	rootCmd.AddCommand(applyCmd)
}

func Apply() {
	m := manifest.Manifest{}
	a := application.Application{}
	p := pipeline.Pipeline{}

	manifests := m.GetManifests(filePath)
	for _, manifest := range manifests {
		switch manifest.Kind {
		case a.Manifest.Kind:
			a.LoadManifest(manifest)
			a.Apply(true, plan)
		case p.Manifest.Kind:
			p.LoadManifest(manifest)
		}
	}
}
