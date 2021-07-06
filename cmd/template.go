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
	"swinch/domain"
)

// templateCmd represents the template command
var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Generate manifests from chart domain",
	Long:  "Template command will generate the Spinnaker Application manifest from the chart domain",
	PreRun: func(cmd *cobra.Command, args []string) {
		SetLogLevel(logLevel)
	},
	Run: func(cmd *cobra.Command, args []string) {
		Template()
	},
}

func init() {
	templateCmd.Flags().StringVarP(&chartPath, "chartPath", "c", "", "Dir path for chart")
	templateCmd.Flags().StringVarP(&valuesFilePath, "valuesFilePath", "f", "", "Overwrite chart values file")
	templateCmd.Flags().StringVarP(&outputPath, "outputPath", "o", "", "Dir path for writing templated manifests")
	templateCmd.MarkFlagRequired("chartPath")
	templateCmd.MarkFlagRequired("outputPath")
	rootCmd.AddCommand(templateCmd)
}

func Template() {
	t := domain.Template{}
	t.RenderChart(chartPath, valuesFilePath, outputPath)
}
