package pipeline

import "strconv"

func (p *Pipeline) ProcessStages() {
	for i := 0; i < len(p.Manifest.Spec.Stages); i++ {
		stageMap := &p.Manifest.Spec.Stages[i]
		metadata := p.GetStageMetadata(stageMap)
		metadata.RefId = strconv.Itoa(i + 1)
		switch metadata.Type {
		case bakeManifest:
			p.ProcessBakeManifest(stageMap, &metadata)
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
