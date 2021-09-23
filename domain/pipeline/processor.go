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

		//Overwrite the initial stage map with he newly generated stage spec
		*ps.InitStage = *ps.Types[stageType].MakeStage(&ps.Stage)
	}
}
