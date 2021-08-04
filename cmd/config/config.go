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
	"github.com/manifoldco/promptui"
	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"swinch/domain"
)

const (
	CfgFolderName = "/.swinch/"
	CfgFolderPerm = 0700

	CfgFileName = "config.yaml"
	CfgFilePerm = 0600

	CfgSpinFileName = "context-spin-config.yaml"
	CfgSpinFilePerm = 0600
)

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

type ContextDefinition struct {
	Name     string
	Endpoint string
	Auth     string
	Username string
	Password string
}

type CurrentContext struct {
	Name string
}

type CPrompt struct {
	PUI       promptui.Prompt
	FieldName string
}

func (scf SpinConfigFile) GenerateSpinConfigFile() {
	cd := ContextDefinition{}
	ctx, _ := cd.GetContexts()

	cc := CurrentContext{}
	currentCtx := cc.GetCurrentContext()

	var cfg SpinConfigFile

	for _, context := range ctx {
		if context.Name == currentCtx {
			cfg.Gate.Endpoint = context.Endpoint
			cfg.Auth.Enabled = true

			switch context.Auth {
			case "ldap":
				cfg.Auth.Ldap.Username = context.Username
				cfg.Auth.Ldap.Password = Base64Actions("decode", context.Password)
			case "basic":
				cfg.Auth.Basic.Username = context.Username
				cfg.Auth.Basic.Password = Base64Actions("decode", context.Password)
			}
		}
	}

	ds := domain.Datastore{}
	spinCfgFile := ds.MarshalYAML(&cfg)
	ds.WriteFile(HomeFolder()+CfgFolderName+CfgSpinFileName, spinCfgFile, CfgSpinFilePerm)
}

func (cd ContextDefinition) GetContexts() ([]ContextDefinition, []string) {
	var ctx []ContextDefinition
	var ctxList []string

	if err := viper.UnmarshalKey("contexts", &ctx); err != nil {
		log.Fatalf("Error reading contexts: %s", err)
		return nil, nil
	} else {
		// ctxList string slice needed by promptui; adding to list only contexts with all the fields set
		for _, v := range ctx {
			if v.Name != "" && v.Endpoint != "" && v.Auth != "" && v.Username != "" && v.Password != "" {
				ctxList = append(ctxList, v.Name)
			}
		}
		return ctx, ctxList
	}
}

func (cc CurrentContext) GetCurrentContext() string {
	var currentCtx CurrentContext
	if err := viper.UnmarshalKey("current-context", &currentCtx); err != nil {
		log.Fatalf("Error reading current context: %s", err)
		return ""
	} else {
		return currentCtx.Name
	}
}

func Base64Actions(action, data string) string {
	switch action {
	case "encode":
		encodedData := base64.StdEncoding.EncodeToString([]byte(data))
		return encodedData
	case "decode":
		decodedData, err := base64.StdEncoding.DecodeString(data)
		if err != nil {
			log.Fatalf("Error decoding data: %s", err)
			return ""
		} else {
			return string(decodedData)
		}
	default:
		return ""
	}
}

// ValidateCurrentContext function validates that 'current-context' exists in the contexts list and it is valid (all fields populated); returns bool type
func ValidateCurrentContext() bool {
	cd := ContextDefinition{}
	_, ctxList := cd.GetContexts()

	cc := CurrentContext{}
	currentCtx := cc.GetCurrentContext()

	contextExists := false

	for _, context := range ctxList {
		if currentCtx == context {
			contextExists = true
		}
	}

	return contextExists
}

func HomeFolder() string {
	home, err := homedir.Dir()
	if err != nil {
		log.Fatalf("Error establishing the home directory: %s", err)
		return ""
	} else {
		return home
	}
}
