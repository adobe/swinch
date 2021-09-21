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

package application

import (
	log "github.com/sirupsen/logrus"
	"swinch/domain/datastore"
	"swinch/domain/util"
	"swinch/spincli"
)

type Application struct {
	Manifest
	spincli.ApplicationAPI
	util.Util
	datastore.Datastore
}

func (a *Application) Plan() {
	a.Apply(true, true)
}

func (a *Application) Apply(dryRun, plan bool) {
	existingApp := a.Get(a.Metadata.Name)
	changes := false
	newApp := false
	if len(existingApp) == 0 {
		newApp = true
	} else {
		changes = a.Changes(a.MarshalJSON(a.loadSpec(existingApp)), a.MarshalJSON(a.Spec))
	}

	if changes && plan {
		log.Infof("Planing changes for application '%v'", a.Metadata.Name)
		a.DiffChanges(a.MarshalJSON(a.loadSpec(existingApp)), a.MarshalJSON(a.Spec))
	}

	if !dryRun && (changes || newApp) {
		log.Infof("Saving application '%v'", a.Metadata.Name)
		a.Save(a.Metadata.Name, a.WriteJSONTmp(a.Spec))
	}
}

func (a *Application) Destroy() {
	a.Delete(a.Metadata.Name)
}
