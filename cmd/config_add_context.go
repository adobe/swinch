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
	"errors"
	"github.com/manifoldco/promptui"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os/user"
	"swinch/cmd/config"
)

// addContextCmd represents the add-context command
var addContextCmd = &cobra.Command{
	Use:   "add-context",
	Short: "Adds a new Spinnaker context to the config file",
	Long:  `Adds a new Spinnaker context to the config file`,
	PreRun: func(cmd *cobra.Command, args []string) {
		SetLogLevel(logLevel)
		ValidateConfigFile()
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := viper.ReadInConfig(); err == nil {
			fields := addContextPromptUI()
			if fields["confirmation"] == "y" {
				addNewContext(fields)
			}
		}
	},
}

func init() {
	configCmd.AddCommand(addContextCmd)
}

func addContextPromptUI() map[string]string {
	// Retrieve the current user for autocompletion purposes
	username := ""
	u, err := user.Current()
	if err == nil {
		username = u.Username
	}

	// Prompt validations
	validateContextDuplication := func(input string) error {
		cd := config.ContextDefinition{}
		ctx, _ := cd.GetContexts()

		if len(input) < 3 {
			return errors.New("must use at least 3 characters")
		} else {
			for _, v := range ctx {
				if input == v.Name {
					return errors.New("context already exists, choose a different name OR delete the entry and retry")
				}
			}
		}
		return nil
	}

	validateAuth := func(input string) error {
		if input != "ldap" && input != "basic" {
			return errors.New("must use ether 'ldap' OR 'basic' auth methods")
		}
		return nil
	}

	validateFieldLength := func(input string) error {
		if len(input) < 3 {
			return errors.New("must use at least 3 characters")
		}
		return nil
	}

	fields := map[string]string{}

	prompts := []config.CPrompt{
		{
			PUI: promptui.Prompt{
				Label:    "Spinnaker context name",
				Default:  "spinnaker-prod",
				Validate: validateContextDuplication,
			},
			FieldName: "name",
		},
		{
			PUI: promptui.Prompt{
				Label:    "Spinnaker API endpoint",
				Default:  "https://spinnaker-prod-api.example.com",
				Validate: validateFieldLength,
			},
			FieldName: "endpoint",
		},
		{
			PUI: promptui.Prompt{
				Label:    "Auth type (ldap OR basic)",
				Default:  "ldap",
				Validate: validateAuth,
			},
			FieldName: "auth",
		},
		{
			PUI: promptui.Prompt{
				Label:    "Username",
				Default:  username,
				Validate: validateFieldLength,
			},
			FieldName: "username",
		},
		{
			PUI: promptui.Prompt{
				Label:    "Password (hidden-base64encoded)",
				Mask:     '*',
				Validate: validateFieldLength,
			},
			FieldName: "password",
		},
		{
			PUI: promptui.Prompt{
				Label:     "Is the information above correct",
				IsConfirm: true,
			},
			FieldName: "confirmation",
		},
	}

	for _, prompt := range prompts {
		fields[prompt.FieldName], err = prompt.PUI.Run()
		if err != nil {
			log.Errorf("Exiting... %v\n", err)
			break
		}
	}
	return fields
}

func addNewContext(fields map[string]string) {
	cd := config.ContextDefinition{}
	ctx, _ := cd.GetContexts()

	newContext := config.ContextDefinition{
		Name:     fields["name"],
		Endpoint: fields["endpoint"],
		Auth:     fields["auth"],
		Username: fields["username"],
		Password: config.Base64Encode(fields["password"]),
	}

	ctx = append(ctx, newContext)

	viper.Set("contexts", ctx)

	err := viper.WriteConfig()
	if err != nil {
		log.Fatalf("Error: %s", err)
	} else {
		log.Infof("New context '%s' was added to the config file; you can now run `swinch config use-context` to set it as current-context", fields["name"])
	}
}
