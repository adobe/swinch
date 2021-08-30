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

package pipeline

import (
	log "github.com/sirupsen/logrus"
	"strconv"
	"swinch/domain/stage"
)

type Pipeline struct {
	Manifest
	Spec
	BakeManifest
	DeployManifest
	DeleteManifest
	ManualJudgment
	Wait
	Jenkins
	RunJobManifest
	stage.Stage
}

const (
	bakeManifest   = "bakeManifest"
	deployManifest = "deployManifest"
	deleteManifest = "deleteManifest"
	manualJudgment = "manualJudgment"
	wait           = "wait"
	jenkins        = "jenkins"
	runJobManifest = "runJobManifest"
)

func (p *Pipeline) ProcessStages() {
	for i := 0; i < len(p.Manifest.Spec.Stages); i++ {
		stageMap := &p.Manifest.Spec.Stages[i]
		metadata := p.GetStageMetadata(stageMap)
		metadata.RefId = strconv.Itoa(i + 1)
		switch metadata.Type {
		case bakeManifest:
			p.ProcessBakeManifest(p, stageMap, &metadata)
		case deployManifest:
			p.ProcessDeployManifest(p, stageMap, &metadata)
		case deleteManifest:
			p.ProcessDeleteManifest(p, stageMap, &metadata)
		case manualJudgment:
			p.ProcessManualJudgment(stageMap, &metadata)
		case wait:
			p.ProcessWait(stageMap, &metadata)
		case jenkins:
			p.ProcessJenkins(stageMap, &metadata)
		case runJobManifest:
			p.ProcessRunJobManifest(p, stageMap, &metadata)
		default:
			log.Fatalf("Failed to detect stage type: %v", metadata.Type)
		}
	}
}
