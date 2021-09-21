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
	"swinch/domain/application"
)

// applicationCmd represents the application command
var applicationCmd = &cobra.Command{
	Use:   "application",
	Short: "Run operations on a Spinnaker application",
	Long:  `Run operations on a Spinnaker application`,
	PreRun: func(cmd *cobra.Command, args []string) {
		SetLogLevel(logLevel)
		ValidateConfigFile()
		ValidateConfig()
	},
	Run: func(cmd *cobra.Command, args []string) {
		action := cmd.Parent().Use
		cmdAppAction(action)
	},
}

var PlanAppCmd = *applicationCmd
var ImportAppCmd = *applicationCmd
var DeleteAppCmd = *applicationCmd

func init() {
	// import flags
	ImportAppCmd.Flags().StringVarP(&applicationName, "application", "a", "", "Application name")
	ImportAppCmd.Flags().StringVarP(&filePath, "file", "f", "", "JSON file input")
	ImportAppCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Generated chart output path")
	ImportAppCmd.Flags().StringVarP(&chartName, "chart", "n", "", "Specify chart name for imported pipeline")
	ImportAppCmd.Flags().BoolVarP(&protectedImport, "protected-import", "", false, "Protect already created chart from overwriting")
	ImportAppCmd.MarkFlagRequired("application")
	ImportAppCmd.MarkFlagRequired("output")
	importCmd.AddCommand(&ImportAppCmd)

	// delete flags
	DeleteAppCmd.Flags().StringVarP(&applicationName, "application", "a", "", "Application name")
	DeleteAppCmd.MarkFlagRequired("application")
	deleteCmd.AddCommand(&DeleteAppCmd)

	// plan flags
	planCmd.AddCommand(&PlanAppCmd)
}

func cmdAppAction(action string) {
	a := application.Application{}
	switch action {
	case deleteAction:
		a.Delete(applicationName)
	case importAction:
		fmt.Println("Import TBA")
	default:
		log.Fatalf("Bad application action")
	}
}
