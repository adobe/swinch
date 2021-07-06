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
	log "github.com/sirupsen/logrus"
	"strconv"
)

type Pipeline struct {
	PipelineManifest
	PipelineSpec
}

const (
	bakeManifest   = "bakeManifest"
	deployManifest = "deployManifest"
	deleteManifest = "deleteManifest"
)

func (p *Pipeline) ExpandSpec() {
	for i := 0; i < len(p.Spec.Stages); i++ {
		stage := &p.Spec.Stages[i]
		stage.RefId = strconv.Itoa(i + 1)
		log.Debugf("Running stage: %v, RefId: %v", i, stage.RefId)
		switch stage.Type {
		case bakeManifest:
			stage.bakeManifest()
		case deployManifest:
			stage.deployManifest(p)
		case deleteManifest:
			stage.deleteManifest(p)
		}
	}
}
