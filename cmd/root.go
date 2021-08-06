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
	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var (
	cfgFile         string
	logLevel        string
	plan            bool
	filePath        string
	outputPath      string
	valuesFilePath  string
	chartName       string
	protectedImport bool
	application     string
	pipeline        string
	chartPath       string
)

const (
	SwinchVersion = "0.0.6"
)

const (
	planAction   = "plan"
	importAction = "import"
	applyAction  = "apply"
	deleteAction = "delete"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "swinch",

	Short: "Generate Spinnaker applications and pipelines from a kubernetes like objects",
	Long:  "Generate Spinnaker applications and pipelines from a kubernetes like objects",
	Run:   func(cmd *cobra.Command, args []string) { fmt.Println("Try the -h flag for possible options") },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	SetLogLevel(logLevel)
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.swinch.yaml)")
	rootCmd.PersistentFlags().StringVarP(&logLevel, "verbosity", "v", log.InfoLevel.String(), "Log level (debug, info, warn, error, fatal, panic")

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// SetLogLevel set the log level
func SetLogLevel(logLevel string) {
	level, _ := log.ParseLevel(logLevel)
	log.SetLevel(level)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".swinch" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".swinch")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
