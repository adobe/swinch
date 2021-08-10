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
	"github.com/spf13/viper"
	"path"
	"swinch/cmd/config"
	"swinch/domain"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generates a mock swinch config file in '~/.swinch/config.yaml'",
	Long:  `Generates a mock swinch config file in '~/.swinch/config.yaml'`,
	PreRun: func(cmd *cobra.Command, args []string) {
		SetLogLevel(logLevel)
	},
	Run: func(cmd *cobra.Command, args []string) {
		generateConfig()
	},
}

func init() {
	configCmd.AddCommand(generateCmd)
}

func generateConfig() {
	contexts := []config.ContextDefinition{
		{
			Name:     "spinnaker-dev",
			Endpoint: "https://spinnaker-dev-api.example.com",
			Auth:     "ldap",
			Username: "username",
			Password: config.Base64Encode("base64EncodedPassword"),
		},
		{
			Name:     "spinnaker-prod",
			Endpoint: "https://spinnaker-prod-api.example.com",
			Auth:     "basic",
			Username: "username",
			Password: config.Base64Encode("base64EncodedPassword"),
		},
	}

	currentContext := config.CurrentContext{Name: "spinnaker-dev"}

	// Populate the generated config file with some mock values and set secure file permissions
	viper.SetDefault("contexts", contexts)
	viper.SetDefault("current-context", currentContext)
	viper.SetConfigPermissions(config.CfgFilePerm)

	// Write the config file to the default location ~/.swinch/config.yaml
	ds := domain.Datastore{}
	ds.Mkdir(path.Join(config.HomeFolder(), config.CfgFolderName), config.CfgFolderPerm)

	err := viper.SafeWriteConfigAs(config.HomeFolder() + config.CfgFolderName + config.CfgFileName)
	if err != nil {
		log.Fatalf("Error: %s", err)
	} else {
		log.Infof("Writing a mock swinch config file in: %s",
			config.HomeFolder()+config.CfgFolderName+config.CfgFileName+
				" \nNow run 'swinch config add-context' to append a real context to the config file and then 'swinch config use-context' to set it as current context")
	}
}
