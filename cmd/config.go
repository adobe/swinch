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
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Tweak swinch config",
	Long:  `Tweak swinch config - setup is similar with kubernetes kubeconfig; one can easily define multiple contexts and switch between them`,
	Example: `Steps to initialize a custom config:
	swinch config generate (generates a mock config file with some example entries)
	swinch config add-context (add a new context to the list)
	swinch config get-contexts (prints all contexts)
	swinch config use-context (select a new current-context)
	swinch config delete-context (delete the example entries)`,
	PreRun: func(cmd *cobra.Command, args []string) {
		SetLogLevel(logLevel)
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
