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

package config

import (
	"encoding/base64"
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"swinch/domain/datastore"
)

const (
	CfgFolderName = "/.swinch/"
	CfgFolderPerm = 0700

	CfgFileName = "config.yaml"
	CfgFilePerm = 0600

	CfgSpinFileName = "context-spin-config.yaml"
	CfgSpinFilePerm = 0600
)

// SpinConfigFile struct used to populate ~/.swinch/context-spin-config.yaml; this file is served to the spin-cli calls
type SpinConfigFile struct {
	Gate struct {
		Endpoint string `yaml:"endpoint"`
	} `yaml:"gate"`

	Auth struct {
		Enabled bool `yaml:"enabled"`
		Ldap    struct {
			Username string `yaml:"username"`
			Password string `yaml:"password"`
		} `yaml:"ldap,omitempty"`
		Basic struct {
			Username string `yaml:"username"`
			Password string `yaml:"password"`
		} `yaml:"basic,omitempty"`
	} `yaml:"auth"`
}

// ContextDefinition struct used to populate ~/.swinch/config.yaml contexts
type ContextDefinition struct {
	Name     string
	Endpoint string
	Auth     string
	Username string
	Password string
}

// CurrentContext struct used to populate ~/.swinch/config.yaml current-context
type CurrentContext struct {
	Name string
}

type CPrompt struct {
	PUI       promptui.Prompt
	FieldName string
}

// GenerateSpinConfigFile method parses the ~/.swinch/config.yaml file, validates the current-context against all the available contexts
// and creates/populates the ~/.swinch/context-spin-config.yaml file used by the spin-cli calls
func (scf SpinConfigFile) GenerateSpinConfigFile() {
	cd := ContextDefinition{}
	ctx, _ := cd.GetContexts()

	cc := CurrentContext{}
	currentCtx := cc.GetCurrentContext()

	for _, context := range ctx {
		if context.Name == currentCtx {
			scf.Gate.Endpoint = context.Endpoint
			scf.Auth.Enabled = true

			switch context.Auth {
			case "ldap":
				scf.Auth.Ldap.Username = context.Username
				scf.Auth.Ldap.Password = Base64Decode(context.Password)
			case "basic":
				scf.Auth.Basic.Username = context.Username
				scf.Auth.Basic.Password = Base64Decode(context.Password)
			}
		}
	}

	d := datastore.Datastore{}
	spinCfgFile := d.MarshalYAML(&scf)
	d.WriteFile(HomeFolder()+CfgFolderName+CfgSpinFileName, spinCfgFile, CfgSpinFilePerm)
}

// GetContexts method parses the ~/.swinch/config.yaml file and returns the the contexts
func (cd ContextDefinition) GetContexts() ([]ContextDefinition, []string) {
	var ctx []ContextDefinition
	var ctxList []string

	if err := viper.UnmarshalKey("contexts", &ctx); err != nil {
		log.Fatalf("Error reading contexts: %s", err)
	}

	// ctxList string slice needed by promptui; adding to list only contexts with all the fields set
	for _, v := range ctx {
		if v.Name != "" && v.Endpoint != "" && v.Auth != "" && v.Username != "" && v.Password != "" {
			ctxList = append(ctxList, v.Name)
		}
	}
	return ctx, ctxList
}

// GetCurrentContext method parses the ~/.swinch/config.yaml file and returns the current-context as string type
func (cc CurrentContext) GetCurrentContext() string {
	if err := viper.UnmarshalKey("current-context", &cc); err != nil {
		log.Fatalf("Error reading current context: %s", err)
	}
	return cc.Name
}

// ValidateCurrentContext function validates that 'current-context' exists in the contexts list and it is valid (all fields populated); returns bool type
func (cd ContextDefinition) ValidateCurrentContext() error {
	_, ctxList := cd.GetContexts()

	cc := CurrentContext{}
	currentCtx := cc.GetCurrentContext()

	for _, context := range ctxList {
		if currentCtx == context {
			return nil
		}
	}

	return fmt.Errorf("curent context '%s' not valid", currentCtx)
}

func Base64Encode(data string) string {
	encodedData := base64.StdEncoding.EncodeToString([]byte(data))
	return encodedData
}

func Base64Decode(data string) string {
	decodedData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		log.Fatalf("Error decoding data: %s", err)
	}
	return string(decodedData)
}

// HomeFolder function retrieves the user's home folder; returns the home folder path as string type
func HomeFolder() string {
	home, err := homedir.Dir()
	if err != nil {
		log.Fatalf("Error establishing the home directory: %s", err)
	}
	return home
}
