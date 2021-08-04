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
	"fmt"
	log "github.com/sirupsen/logrus"
	"swinch/cmd/config"
)

type ApplicationAPI struct {
	App string
	SpinCLI
}

var baseArgs = []string{
	"--no-color=false",
	"--config", config.HomeFolder()+config.CfgFolderName+config.CfgSpinFileName,
}

func (a *ApplicationAPI) NotFound() error {
	return fmt.Errorf("Application '%v' not found\n", a.App)
}

func (a *ApplicationAPI) deleteNotFound() error {
	return fmt.Errorf("attempting to delete application '%v' which does not exist, exiting", a.App)
}

func (a *ApplicationAPI) Get() []byte {
	args := []string{"application", "get", a.App}
	buffer, err := a.executeAppCmd(append(baseArgs, args...))
	a.status(err)
	return buffer.Bytes()
}

func (a ApplicationAPI) Save(filePath string) {
	args := []string{"application", "save", "--file", filePath}
	_, err := a.executeAppCmd(append(baseArgs, args...))
	a.status(err)
	if err == nil {
		log.Infof("Application '%v' updated successfuly", a.App)
	}
	defer a.rmTmp(filePath)
}

func (a ApplicationAPI) Delete() {
	args := []string{"application", "delete", a.App}
	_, err := a.executeAppCmd(append(baseArgs, args...))
	if err != nil {
		a.status(err)
	} else {
		log.Infof("Delete application '%v' success", a.App)
	}
}

func (a *ApplicationAPI) status(err error) {
	if err != nil {
		switch err.Error() {
		case a.NotFound().Error():
			log.Info(a.NotFound())
		case a.deleteNotFound().Error():
			log.Info(a.deleteNotFound())
		default:
			log.Fatalf("Failed to check application status: %v", err)
		}
	}
}
