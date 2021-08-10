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
	"swinch/cmd/config"
)

// deleteContextCmd represents the delete-context command
var deleteContextCmd = &cobra.Command{
	Use:   "delete-context",
	Short: "Deletes a Spinnaker context from the config file",
	Long:  `Deletes a Spinnaker context from the config file`,
	PreRun: func(cmd *cobra.Command, args []string) {
		SetLogLevel(logLevel)
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := viper.ReadInConfig(); err == nil {
			context := deleteContextPromptUI()
			if context != "" {
				deleteContext(context)
			}
		}
	},
}

func init() {
	configCmd.AddCommand(deleteContextCmd)
}

func deleteContextPromptUI() string {
	cd := config.ContextDefinition{}
	_, ctxList := cd.GetContexts()

	if len(ctxList) == 0 {
		log.Fatalf("The config file does not have any valid contexts")
	}

	prompt := promptui.Select{
		Label: "Delete Spinnaker Context",
		Items: ctxList,
	}
	_, context, err := prompt.Run()
	if err != nil {
		log.Fatalf("Exiting %v\n", err)
	}
	return context
}

func deleteContext(context string) {
	cd := config.ContextDefinition{}
	ctx, _ := cd.GetContexts()

	cc := config.CurrentContext{}
	currentCtx := cc.GetCurrentContext()

	var updatedCtx []config.ContextDefinition

	// Removing the context selected for deletion from the contexts list in the new updatedCtx slice
	for _, v := range ctx {
		if v.Name != context {
			updatedCtx = append(updatedCtx, v)
		}
	}

	viper.Set("contexts", updatedCtx)

	err := viper.WriteConfig()
	if err != nil {
		log.Errorf("Error: %s", err)
	} else {
		if context == currentCtx {
			log.Warnf("Context '%s' selected for deletion is set as 'current-context'; run 'swinch config use-context' to select another current-context", context)
		} else {
			log.Infof("Context '%s' was deleted from the config file", context)
		}
	}
}
