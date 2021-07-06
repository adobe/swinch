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

package spincli

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	spincmd "github.com/spinnaker/spin/cmd"
	spinapplication "github.com/spinnaker/spin/cmd/application"
	spinpipeline "github.com/spinnaker/spin/cmd/pipeline"
	"os"
)

type SpinCLI struct {
}

func (s SpinCLI) getCmd() (*bytes.Buffer, *cobra.Command, *spincmd.RootOptions) {
	buffer := new(bytes.Buffer)
	cmd, options := spincmd.NewCmdRoot(buffer, buffer)
	return buffer, cmd, options
}

func (s SpinCLI) executeAppCmd(args []string) (bytes.Buffer, error) {
	buffer, cmd, options := s.getCmd()
	appCmd := spinapplication.NewApplicationCmd(options)
	cmd.AddCommand(appCmd)
	cmd.SetArgs(args)
	return *buffer, cmd.Execute()
}

func (s SpinCLI) executePipeCmd(args []string) (bytes.Buffer, error) {
	buffer, cmd, options := s.getCmd()
	pipeCmd, _ := spinpipeline.NewPipelineCmd(options)
	cmd.AddCommand(pipeCmd)
	cmd.SetArgs(args)
	return *buffer, cmd.Execute()
}

func (s SpinCLI) rmTmp(file string) {
	errRemove := os.Remove(file)
	if errRemove != nil {
		log.Warnf("error removing the temp file: %v", errRemove)
	}
}
