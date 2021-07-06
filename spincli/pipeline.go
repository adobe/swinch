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
	"errors"
	log "github.com/sirupsen/logrus"
)

type PipelineAPI struct {
	App  string `json:"application"`
	Pipe string `json:"name"`
	SpinCLI
}

func (p *PipelineAPI) NotFound() error {
	return errors.New("404 ")
}

func (p *PipelineAPI) NotAllowed() error {
	return errors.New("403 ")
}

func (p *PipelineAPI) unhandledNotFound() error {
	return errors.New("unhandled response 404: ")
}

func (p *PipelineAPI) unhandledNotAllowed() error {
	return errors.New("unhandled response 403: ")
}

func (p *PipelineAPI) renameNotAllowed() error {
	return errors.New("400 ")
}

func (p *PipelineAPI) Get() []byte {
	args := []string{"pipeline", "get", "--application", p.App, "--name", p.Pipe}
	buffer, err := p.executePipeCmd(append(baseArgs, args...))
	log.Debugf("Spinnaker get response: %v", err)
	p.status(err)
	return buffer.Bytes()
}

func (p PipelineAPI) Save(filePath string) {
	args := []string{"pipeline", "save", "--file", filePath}
	_, err := p.executePipeCmd(append(baseArgs, args...))
	p.status(err)
	if err == nil {
		log.Infof("Pipeline '%v' in application '%v' updated successfuly", p.Pipe, p.App)
	}
	defer p.rmTmp(filePath)
}

func (p PipelineAPI) Delete() {
	args := []string{"pipeline", "delete", "--application", p.App, "--name", p.Pipe}
	_, err := p.executePipeCmd(append(baseArgs, args...))
	if err != nil {
		p.status(err)
	} else {
		log.Infof("Delete pipeline '%v' success", p.Pipe)
	}
}

func (p *PipelineAPI) status(err error) {
	if err != nil {
		switch err.Error() {
		case p.NotFound().Error():
			log.Infof("Pipeline '%v' not found", p.Pipe)
		case p.NotAllowed().Error():
			log.Fatalf("Attempting action on pipeline '%v' from application '%v' which does not exist", p.Pipe, p.App)
		case p.unhandledNotFound().Error():
			log.Fatalf("Request repeated too quickly")
		case p.unhandledNotAllowed().Error():
			log.Fatalf("Request repeated too quickly")
		case p.renameNotAllowed().Error():
			log.Fatalf("Renaming an existing pipeline is not supported")
		default:
			log.Fatalf("Failed to check pipeline status: %v", err)
		}
	}
}
