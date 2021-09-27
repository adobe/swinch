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

package stages

import (
	"encoding/json"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"swinch/domain/datastore"
)

const ethosNamespaceCreate StageType = "EthosNamespaceCreate"

type EthosNamespaceCreate struct {
	Metadata `mapstructure:",squash"`
	Common   `mapstructure:",squash"`

	Alias      string                         `yaml:"alias" json:"alias"`
	IsNew      bool                           `yaml:"isNew,omitempty" json:"isNew,omitempty"`
	Parameters EthosNamespaceCreateParameters `yaml:"parameters" json:"parameters"`
}

type EthosNamespaceCreateParameters struct {
	AdirCogs           string `yaml:"AdirCogs" json:"AdirCogs"`
	AdusCogs           string `yaml:"AdusCogs" json:"AdusCogs"`
	AdusOpex           string `yaml:"AdusOpex" json:"AdusOpex"`
	Cluster            string `yaml:"Cluster" json:"Cluster"`
	LdapEditGroup      string `yaml:"LdapEditGroup" json:"LdapEditGroup"`
	LdapViewGroup      string `yaml:"LdapViewGroup" json:"LdapViewGroup"`
	Namespace          string `yaml:"Namespace" json:"Namespace"`
	Project            string `yaml:"Project" json:"Project"`
	SKMSservice        string `yaml:"SKMSservice" json:"SKMSservice"`
	SpinnakerOnlyGroup string `yaml:"SpinnakerOnlyGroup" json:"SpinnakerOnlyGroup"`
}

func (enc EthosNamespaceCreate) MakeStage(stage *Stage) *map[string]interface{} {
	enc.decode(stage)
	return enc.encode()
}

func (enc *EthosNamespaceCreate) decode(stage *Stage) {
	decoderConfig := mapstructure.DecoderConfig{WeaklyTypedInput: true, Result: &enc}
	decoder, err := mapstructure.NewDecoder(&decoderConfig)
	if err != nil {
		log.Fatalf("err: %v", err)
	}

	err = decoder.Decode(stage.Metadata)
	if err != nil {
		log.Fatalf("error decoding stage metadata: %v", err)
	}
	err = decoder.Decode(stage.Spec)
	if err != nil {
		log.Fatalf("error decoding stage spec: %v", err)
	}
}

func (enc *EthosNamespaceCreate) encode() *map[string]interface{} {
	d := datastore.Datastore{}
	stage := new(map[string]interface{})
	err := json.Unmarshal(d.MarshalJSON(enc), stage)
	if err != nil {
		log.Fatalf("Failed to unmarshal JSON:  %v", err)
	}
	return stage
}
