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
	"os"
	"swinch/domain"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs a swinch chart",
	Long:  `Installs a swinch chart.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		SetLogLevel(logLevel)
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Template call
		datastore := domain.Datastore{}
		outputPath = datastore.CreateTmpFolder()
		defer os.RemoveAll(outputPath)
		templateCmd.Run(cmd, []string{})

		// Apply call
		filePath = outputPath
		applyCmd.Run(cmd, []string{})
	},
}

func init() {
	installCmd.Flags().StringVarP(&chartPath, "chartPath", "c", "", "Dir path for chart")
	installCmd.Flags().StringVarP(&valuesFilePath, "valuesFilePath", "f", "", "Overwrite chart values file")
	installCmd.Flags().BoolVarP(&plan, "plan", "p", true, "Display plan while installing, no user input.")
	installCmd.MarkFlagRequired("chartPath")
	rootCmd.AddCommand(installCmd)
}
