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

package domain

import (
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
)

type StageMetadata struct {
	Name                 string   `yaml:"name" json:"name"`
	Type                 string   `yaml:"type,omitempty" json:"type,omitempty"`
	RefId                string   `yaml:"refId,omitempty" json:"refId,omitempty"`
	RequisiteStageRefIds []string `yaml:"requisiteStageRefIds" json:"requisiteStageRefIds"`
}

func (sm *StageMetadata) getStageMetadata(stage *map[string]interface{}) StageMetadata {
	sm.decode(stage)
	return *sm
}

func (sm *StageMetadata) decode(stage *map[string]interface{}) {
	decoderConfig := mapstructure.DecoderConfig{WeaklyTypedInput: true, Result: &sm}
	decoder, err := mapstructure.NewDecoder(&decoderConfig)
	if err != nil {
		log.Fatalf("err: %v", err)
	}

	err = decoder.Decode(stage)
	if err != nil {
		log.Fatalf("err: %v", err)
	}
}
