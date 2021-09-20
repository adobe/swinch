package pipeline

import (
	log "github.com/sirupsen/logrus"
	"strconv"
	"swinch/domain/stages"
)

type Processor struct {
	Manifest
	stages.Stage
	stages.BakeManifest
	stages.DeployManifest
	stages.DeleteManifest
	stages.ManualJudgment
	stages.Wait
	stages.Jenkins
	stages.RunJobManifest
}

type stage interface {
	Process(*stages.Stage)
	GetStageType() string
}

func (ps Processor) process(manifest *Manifest) {
	ps.Manifest = *manifest
	ps.Stages = &ps.Manifest.Spec.Stages
	for i := 0; i < len(ps.Manifest.Spec.Stages); i++ {
		ps.RawStage = &ps.Manifest.Spec.Stages[i]
		// Decode the raws stage from the manifest into an internal stage struct
		ps.Stage = ps.GetStage(ps.RawStage)

		// Set some stage metadata
		ps.Stage.Metadata.RefId = strconv.Itoa(i + 1)

		// Propagate the manifest metadata to the stage
		ps.Stage.ManifestMetadata.Name = ps.Manifest.Metadata.Name
		ps.Stage.ManifestMetadata.Application = ps.Manifest.Metadata.Application

		switch ps.Stage.Type {
		case ps.BakeManifest.GetStageType():
			var s stage = ps.BakeManifest
			s.Process(&ps.Stage)
		case ps.DeleteManifest.GetStageType():
			var s stage = ps.DeleteManifest
			s.Process(&ps.Stage)
		case ps.DeployManifest.GetStageType():
			var s stage = ps.DeployManifest
			s.Process(&ps.Stage)
		case ps.Jenkins.GetStageType():
			var s stage = ps.Jenkins
			s.Process(&ps.Stage)
		case ps.ManualJudgment.GetStageType():
			var s stage = ps.ManualJudgment
			s.Process(&ps.Stage)
		case ps.RunJobManifest.GetStageType():
			var s stage = ps.RunJobManifest
			s.Process(&ps.Stage)
		case ps.Wait.GetStageType():
			var s stage = ps.RunJobManifest
			s.Process(&ps.Stage)
		default:
			log.Fatalf("Failed to detect stage type: %v", ps.Stage.Metadata.Type)
		}
	}
}
