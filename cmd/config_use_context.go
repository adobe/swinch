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
	"github.com/manifoldco/promptui"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"swinch/cmd/config"
)

// useContextCmd represents the use-context command
var useContextCmd = &cobra.Command{
	Use:   "use-context",
	Short: "Switches between Spinnaker contexts",
	Long:  `Switches between Spinnaker contexts; works in both interactive or non-interactive way`,
	Example: `Interactive vs non-interactive:
	swinch config use-context (displays the prompt to select a context)
	swinch config use-context spinnaker-dev (non-interactive, switches to the new context directly if it exists)`,
	PreRun: func(cmd *cobra.Command, args []string) {
		SetLogLevel(logLevel)
		ValidateConfigFile()
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := viper.ReadInConfig(); err == nil {
			context := useContextPromptUI()
			if context != "" {
				changeCurrentContext(context)
			}
		}
	},
}

func init() {
	configCmd.AddCommand(useContextCmd)
}

func useContextPromptUI() string {
	cd := config.ContextDefinition{}
	_, ctxList := cd.GetContexts()

	if len(ctxList) == 0 {
		log.Fatalf("The config file does not have any valid contexts")
	}

	_, args, _ := rootCmd.Find(os.Args)

	context := new(string)

	// Allow 'swinch config use-context context-name' subcommand to run without promptui
	if len(args) == 4 && args[1] == "config" && args[2] == "use-context" {
		for _, existingContext := range ctxList {
			if args[3] == existingContext {
				*context = args[3]
				break
			}
		}

		if *context == "" {
			log.Fatalf("The specified context '%s' does not exist in the contexts list", args[3])
		}
	} else {
		prompt := promptui.Select{
			Label: "Set a new Spinnaker Context",
			Items: ctxList,
		}
		_, ctx, err := prompt.Run()
		*context = ctx
		if err != nil {
			log.Fatalf("Exiting %v\n", err)
		}
	}
	return *context
}

func changeCurrentContext(newContext string) {
	viper.Set("current-context.name", newContext)

	err := viper.WriteConfig()
	if err != nil {
		log.Errorf("Error: %s", err)
	} else {
		log.Infof("Current context is now set to: '%s'", newContext)
	}

	scf := config.SpinConfigFile{}
	scf.GenerateSpinConfigFile()
}
