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
	"bytes"
	"github.com/danielcoman/diff"
	log "github.com/sirupsen/logrus"
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
	},
	Run: func(cmd *cobra.Command, args []string) {
		PlanCmd()
	},
}

func init() {
	planCmd.Flags().StringVarP(&filePath, "filePath", "f", "", "Manifest file or directory, non recursive")
	planCmd.MarkFlagRequired("filePath")
	rootCmd.AddCommand(planCmd)
}

func PlanCmd() {
	m := manifest.Manifest{}
	a := Application{}
	p := Pipeline{}
	a.manifests, p.manifests = m.GetManifests(filePath)

	if len(a.manifests) > 0 {
		a.manifestActions(planAction)
	}
	if len(p.manifests) > 0 {
		p.manifestActions(planAction)
	}
}

func Changes(oldData, newData []byte) bool {
	changes := bytes.Compare(oldData, newData)
	if changes == 0 {
		log.Infof("No changes detected")
		return false
	}

	return true
}

func DiffChanges(oldData, newData []byte) {
	log.Infof(diff.LineDiff(string(oldData), string(newData)))
}
