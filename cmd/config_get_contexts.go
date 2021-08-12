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
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	log "github.com/sirupsen/logrus"
	"os"
	"swinch/cmd/config"

	"github.com/spf13/cobra"
)

// getContextsCmd represents the get-contexts command
var getContextsCmd = &cobra.Command{
	Use:   "get-contexts",
	Short: "Lists the available contexts present in the config file ~/.swinch/config.yaml",
	Long:  `Lists the available contexts present in the config file ~/.swinch/config.yaml`,
	PreRun: func(cmd *cobra.Command, args []string) {
		SetLogLevel(logLevel)
		ValidateConfigFile()
		ValidateConfig()
	},
	Run: func(cmd *cobra.Command, args []string) {
		getContexts()
	},
}

func init() {
	configCmd.AddCommand(getContextsCmd)
}

func getContexts() {
	cd := config.ContextDefinition{}
	ctx, _ := cd.GetContexts()

	cc := config.CurrentContext{}
	currentCtx := cc.GetCurrentContext()

	if len(ctx) == 0 {
		log.Fatalf("The config file does not have any valid contexts")
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "NAME", "CURRENT-CONTEXT", "ENDPOINT", "AUTH-METHOD", "USERNAME", "PASSWORD"})
	t.AppendSeparator()

	for index, context := range ctx {
		if context.Name == currentCtx {
			t.AppendRow([]interface{}{index + 1, context.Name, "*", context.Endpoint, context.Auth, context.Username, "hidden"})
		} else {
			t.AppendRow([]interface{}{index + 1, context.Name, " ", context.Endpoint, context.Auth, context.Username, "hidden"})
		}
	}

	t.SetStyle(table.Style{
		Name: "swinch",
		Box: table.BoxStyle{
			PaddingLeft:  "",
			PaddingRight: "     ",
		},
	})

	t.SetColumnConfigs([]table.ColumnConfig{
		{
			Name:        "NAME",
			Align:       text.AlignLeft,
			AlignHeader: text.AlignCenter,
		},
		{
			Name:        "CURRENT-CONTEXT",
			Align:       text.AlignCenter,
			AlignHeader: text.AlignCenter,
		},
		{
			Name:        "ENDPOINT",
			Align:       text.AlignLeft,
			AlignHeader: text.AlignCenter,
		},
		{
			Name:        "AUTH-METHOD",
			Align:       text.AlignCenter,
			AlignHeader: text.AlignCenter,
		},
		{
			Name:        "USERNAME",
			Align:       text.AlignCenter,
			AlignHeader: text.AlignCenter,
		},
		{
			Name:        "PASSWORD",
			Align:       text.AlignCenter,
			AlignHeader: text.AlignCenter,
		},
	})

	t.Render()
}
