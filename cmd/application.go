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
	"swinch/domain/application"
	"swinch/domain/chart"
	"swinch/spincli"
)

// applicationCmd represents the application command
var applicationCmd = &cobra.Command{
	Use:   "application",
	Short: "Run operations on a Spinnaker application",
	Long:  `Run operations on a Spinnaker application`,
	PreRun: func(cmd *cobra.Command, args []string) {
		SetLogLevel(logLevel)
	},
	Run: func(cmd *cobra.Command, args []string) {
		action := cmd.Parent().Use
		Application{}.cmdActions(applicationName, action)
	},
}

var PlanAppCmd = *applicationCmd
var ImportAppCmd = *applicationCmd
var DeleteAppCmd = *applicationCmd

func init() {
	// import flags
	ImportAppCmd.Flags().StringVarP(&applicationName, "application", "a", "", "Application name")
	ImportAppCmd.Flags().StringVarP(&filePath, "filePath", "f", "", "JSON file input")
	ImportAppCmd.Flags().StringVarP(&outputPath, "outputPath", "o", "", "Generated chart output path")
	ImportAppCmd.Flags().StringVarP(&chartName, "chartName", "n", "", "Specify chart name for imported pipeline")
	ImportAppCmd.Flags().BoolVarP(&protectedImport, "protectedImport", "", false, "Protect already created chart from overwriting")
	ImportAppCmd.MarkFlagRequired("application")
	ImportAppCmd.MarkFlagRequired("outputPath")
	importCmd.AddCommand(&ImportAppCmd)

	// delete flags
	DeleteAppCmd.Flags().StringVarP(&applicationName, "application", "a", "", "Application name")
	DeleteAppCmd.MarkFlagRequired("application")
	deleteCmd.AddCommand(&DeleteAppCmd)

	// plan flags
	planCmd.AddCommand(&PlanAppCmd)
}

type Application struct {
	manifests []application.Manifest
	application.Application
	spincli.ApplicationAPI
	chart.Chart
}

func (a Application) cmdActions(app, action string) {
	a.App = app
	switch action {
	case deleteAction:
		a.Delete()
	case importAction:
		a.importChart()
	default:
		log.Fatalf("Bad application action")
	}
}

func (a Application) manifestActions(action string) {
	for i := 0; i < len(a.manifests); i++ {
		manifest := &a.manifests[i]
		a.App = a.manifests[i].Metadata.Name
		switch action {
		case applyAction:
			dryRun := false
			a.save(manifest.Spec, dryRun)
		case deleteAction:
			a.Delete()
		case planAction:
			dryRun := true
			a.save(manifest.Spec, dryRun)
		default:
			log.Fatalf("Bad application action")
		}
	}
}

func (a *Application) save(spec application.Spec, dryRun bool) {
	app := a.Get()
	changes := false
	newApp := false
	if len(app) == 0 {
		newApp = true
	} else {
		changes = Changes(a.MarshalJSON(a.LoadSpec(app)), a.MarshalJSON(spec))
	}

	if changes && plan {
		log.Infof("Planing changes for application '%v'", a.App)
		DiffChanges(a.MarshalJSON(a.LoadSpec(app)), a.MarshalJSON(spec))
	}

	if !dryRun && (changes || newApp) {
		log.Infof("Saving application '%v'", a.App)
		a.Save(a.WriteJSONTmp(spec))
	}
}

func (a *Application) importChart() {
	a.OutputPath = outputPath
	a.ProtectedImport = protectedImport
	a.Kind = "application"

	data := new([]byte)
	if filePath != "" {
		*data = a.ReadFile(filePath)
	} else {
		*data = a.Get()
	}

	manifest := a.MakeManifest(a.LoadSpec(*data))
	a.ChartMetadata.Name = chartName
	if a.ChartMetadata.Name == "" {
		a.ChartMetadata.Name = manifest.Metadata.Name
	}

	a.ChartValues.Values = map[interface{}]interface{}{a.Kind: map[string]string{"name": manifest.Metadata.Name}}

	a.GenerateChart(manifest)
}
