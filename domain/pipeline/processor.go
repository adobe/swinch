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
	"swinch/domain/stages"
)

type Processor struct {
	Manifest
	stages.Stages
}

func (ps *Processor) processManifest(manifest *Manifest) {
	ps.Stages.GetTypes()
	ps.Manifest = *manifest
	for i := 0; i < len(ps.Manifest.Spec.Stages); i++ {
		ps.Stage = ps.Decode(&ps.Manifest.Spec.Stages[i])
		ps.InitStage = &ps.Manifest.Spec.Stages[i]
		ps.AllStages = &ps.Manifest.Spec.Stages

		// Set some stage metadata
		ps.Stage.Metadata.RefId = strconv.Itoa(i + 1)
		// Propagate the manifest metadata to the stage
		ps.Stage.ManifestMetadata.Name = ps.Manifest.Metadata.Name
		ps.Stage.ManifestMetadata.Application = ps.Manifest.Metadata.Application

		stageType := stages.StageType(ps.Stage.Type)
		_, ok := ps.Types[stageType]
		if !ok {
			log.Fatalf("Failed to detect stage type: %v", ps.Stage.Type)
		}

		// "If stage fails" execution option has 4 scenarios as seen in the WebUI; to set one of them a bool combination of the below parameters is needed
		// to avoid complexity, the user will use ONLY the ifStageFails parameter (which exists only in the yaml)
		ps.Stage.ContinuePipeline = new(bool)
		ps.Stage.FailPipeline = new(bool)
		ps.Stage.CompleteOtherBranchesThenFail = new(bool)
		switch ps.Stage.IfStageFails {
		case "halt the entire pipeline":
			*ps.Stage.ContinuePipeline = false
			*ps.Stage.FailPipeline = true
			*ps.Stage.CompleteOtherBranchesThenFail = false
		case "halt this branch of the pipeline":
			*ps.Stage.ContinuePipeline = false
			*ps.Stage.FailPipeline = false
			*ps.Stage.CompleteOtherBranchesThenFail = false
		case "halt this branch and fail the pipeline once other branches complete":
			*ps.Stage.ContinuePipeline = false
			*ps.Stage.FailPipeline = false
			*ps.Stage.CompleteOtherBranchesThenFail = true
		case "ignore the failure":
			*ps.Stage.ContinuePipeline = true
			*ps.Stage.FailPipeline = false
			*ps.Stage.CompleteOtherBranchesThenFail = false
		// without these defaults, if the ifStageFails parameters is not set inside the chart,
		// the default option will be "halt this branch of the pipeline" instead of "halt the entire pipeline"
		default:
			*ps.Stage.ContinuePipeline = false
			*ps.Stage.FailPipeline = true
			*ps.Stage.CompleteOtherBranchesThenFail = false
		}

		//Overwrite the initial stage map with he newly generated stage spec
		*ps.InitStage = *ps.Types[stageType].MakeStage(&ps.Stage)
	}
}
