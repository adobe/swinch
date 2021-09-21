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
	"swinch/domain/manifest"
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
		runPlan()
	},
}

func init() {
	planCmd.Flags().StringVarP(&filePath, "file", "f", "", "Manifest file or directory, non recursive")
	planCmd.MarkFlagRequired("file")
	rootCmd.AddCommand(planCmd)
}

func Plan(m manifest.M) {
	m.Plan()
}

func runPlan() {
	m := manifest.NewManifest{}
	manifests := m.GetManifests(filePath)
	for _, newManifest := range manifests {
		switch newManifest.Kind {
		case m.Application.GetKind():
			Plan(m.Application.Load(newManifest))
		case m.Pipeline.GetKind():
			Plan(m.Pipeline.Load(newManifest))
		}
	}
}
