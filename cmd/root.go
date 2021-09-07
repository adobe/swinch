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
	"os"
	"swinch/cmd/config"
	"swinch/domain/datastore"
)

var (
	logLevel        string
	plan            bool
	filePath        string
	outputPath      string
	valuesFilePath  string
	chartName       string
	protectedImport bool
	applicationName string
	pipelineName    string
	chartPath       string
	fullRender      bool
)

const (
	importAction = "import"
	deleteAction = "delete"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "swinch",

	Short: "Generate Spinnaker applications and pipelines from a kubernetes like objects",
	Long:  "Generate Spinnaker applications and pipelines from a kubernetes like objects",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	SetLogLevel(logLevel)
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&logLevel, "verbosity", "v", log.InfoLevel.String(), "Log level (debug, info, warn, error, fatal, panic")
}

// SetLogLevel set the log level
func SetLogLevel(logLevel string) {
	level, _ := log.ParseLevel(logLevel)
	log.SetLevel(level)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Search config in "~/.swinch/config.yaml".
	viper.AddConfigPath(config.HomeFolder() + config.CfgFolderName)
	viper.SetConfigName("config.yaml")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv() // read in environment variables that match
}

func ValidateConfigFile() {
	d := datastore.Datastore{}
	if !d.FileExists(config.HomeFolder() + config.CfgFolderName + config.CfgFileName) {
		log.Errorf("Config file not found, please generate one:")
		configCmd.Usage()
		os.Exit(1)
	}
}

func ValidateConfig() {
	cd := config.ContextDefinition{}
	// If a config file is found, read it in and validate current-context against available contexts
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("A parsing error detected in '%s': '%s'", viper.ConfigFileUsed(), err)
	}
	err = cd.ValidateCurrentContext()
	if err != nil {
		log.Fatalf("Context validation error: %s", err)
	}
	log.Debugf("Using config file: '%s' with current-context as '%s'", viper.ConfigFileUsed(), viper.Get("current-context.name"))
}
