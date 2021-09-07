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
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete the Application or Pipeline form a manifest",
	Long:  "Delete the Application or Pipeline form a manifest",
	PreRun: func(cmd *cobra.Command, args []string) {
		SetLogLevel(logLevel)
		ValidateConfigFile()
		ValidateConfig()
	},
	Run: func(cmd *cobra.Command, args []string) {
		Delete()
	},
}

func init() {
	deleteCmd.Flags().StringVarP(&filePath, "file", "f", "", "Manifest file")
	deleteCmd.MarkFlagRequired("file")
	rootCmd.AddCommand(deleteCmd)
}

func Delete() {
	m := manifest.Manifest{}
	a := application.Application{}
	p := Pipeline{}

	manifests := m.GetManifests(filePath)
	for _, manifest := range manifests {

		switch manifest.Kind {
		case a.Manifest.Kind:
			a.LoadManifest(manifest)
			a.Delete()
		case p.Manifest.Kind:
			p.LoadManifest(manifest)
		}
	}
}
