package stages

type Stages struct {
	Types map[StageType]S
	StageType
	Stage
	BakeManifest
	DeleteManifest
	DeployManifest
	ManualJudgment
	RunJobManifest
	Jenkins
	Wait
}

type StageType string

type S interface {
	MakeStage(*Stage) *map[string]interface{}
}

func (ss *Stages) addStageDefinition(stageType StageType, stage S) {
	ss.Types[stageType] = stage
}

func (ss *Stages) GetTypes() {
	ss.Types = make(map[StageType]S)
	ss.addStageDefinition(bakeManifest, BakeManifest{})
	ss.addStageDefinition(deleteManifest, DeleteManifest{})
	ss.addStageDefinition(deployManifest, DeployManifest{})
	ss.addStageDefinition(jenkins, Jenkins{})
	ss.addStageDefinition(manualJudgment, ManualJudgment{})
	ss.addStageDefinition(pipeline, Pipeline{})
	ss.addStageDefinition(runJobManifest, RunJobManifest{})
	ss.addStageDefinition(wait, Wait{})
}
