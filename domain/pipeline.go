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
)

type Pipeline struct {
	PipelineManifest
	PipelineSpec
	StageMetadata
	BakeManifest
	DeployManifest
	DeleteManifest
}

const (
	bakeManifest   = "bakeManifest"
	deployManifest = "deployManifest"
	deleteManifest = "deleteManifest"
	manualJudgment = "manualJudgment"
)


func (p *Pipeline) ProcessStages() {
	for i := 0; i < len(p.Spec.Stages); i++ {
		metadata := p.getStageMetadata(p, i)
		stage := &p.Spec.Stages[i]
		switch metadata.Type {
		case bakeManifest:
			p.ProcessBakeManifest(stage, &metadata)
		case deployManifest:
			p.ProcessDeployManifest(p, stage, &metadata)
		case deleteManifest:
			p.ProcessDeleteManifest(p, stage, &metadata)
		default:
			log.Fatalf("Failed to detect stage type: %v", metadata.Type)
		}
	}
}
