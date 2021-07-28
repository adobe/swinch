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
	"fmt"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
)

type Pipeline struct {
	PipelineManifest
	PipelineSpec
	Bake
	Deploy
}

const (
	BakeManifest   = "bakeManifest"
	deployManifest = "deployManifest"
	deleteManifest = "deleteManifest"
	manualJudgment = "manualJudgment"
)

func (p *Pipeline) ExpandSpec() {
	for i := 0; i < len(p.Spec.Stages); i++ {
		//stage := &p.Spec.Stages[i]
		//stage.RefId = strconv.Itoa(i + 1)
		//log.Debugf("Running stage: %v, RefId: %v", i, stage.RefId)
		stage := new(Stage)
		err := mapstructure.Decode(&p.Spec.Stages[i], stage)
		if err != nil {
			log.Fatalf("err: %v", err)
		}
		//fmt.Println(stage)
		//fmt.Println(stage.Type)
		err = mapstructure.Decode(&p.Spec.Stages[i], stage)
		switch stage.Type {
		case BakeManifest:
			bake := new(Bake)
			err = mapstructure.Decode(&p.Spec.Stages[i], bake)
			if err != nil {
				log.Fatalf("err: %v", err)
			}
			fmt.Println("---------")
			fmt.Println(bake.Type)
			fmt.Println("---------")
			bake.bakeManifest()
			fmt.Println(bake.RequisiteStageRefIds)
			fmt.Println("---------")
			//stageList := make([]StageList, 0)
			//stageList := append(bake)
		case deployManifest:
			deploy := new(Deploy)
			fmt.Println("deploy spec")
			fmt.Println(&p.Spec.Stages[i])
			fmt.Println("---------")
			err = mapstructure.Decode(&p.Spec.Stages[i], deploy)
			if err != nil {
				log.Fatalf("err: %v", err)
			}
			fmt.Println("after decode")
			fmt.Println(deploy.Type)
			fmt.Println(deploy.RequisiteStageRefIds)
			fmt.Println("---------")
			deploy.DeployManifest(p)
		case deleteManifest:
			//stage.deleteManifest(p)
		default:
		}
	}
	//os.Exit(0)
}

type StageList struct {
}

type Graph struct {
	stages []*Stage
	edges map[*Stage][]*Stage
}

func (g *Graph) AddStage(s *Stage) {
	g.stages = append(g.stages, s)
}

func (g *Graph) AddEdge(s1, s2 *Stage){
	if g.edges == nil {
		g.edges = make(map[*Stage][]*Stage)
	}
	g.edges[s1] = append(g.edges[s1], s2)
	g.edges[s2] = append(g.edges[s2], s1)

}
