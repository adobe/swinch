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

package chart

import (
	"github.com/imdario/mergo"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"strings"
	"swinch/domain/datastore"
)

type Values struct {
	Values map[interface{}]interface{}
}

func (v *Values) loadValuesFile(chartPath, valuesFilePaths string, excludeDefaultValues bool) Values {
	paths := v.getPaths(chartPath, valuesFilePaths, excludeDefaultValues)
	d := datastore.Datastore{}

	for _, valuesFilePath := range paths {
		_, err := os.Stat(valuesFilePath)
		os.IsNotExist(err)
		values := d.UnmarshalYAMLValues(d.ReadFile(valuesFilePath))
		if err = mergo.Merge(&v.Values, values, mergo.WithOverride); err != nil {
			log.Fatalf(err.Error())
		}
	}

	return *v
}

func (v Values) getPaths(chartPath, cliPaths string, excludeDefaultValues bool) []string {
	paths := make([]string, 0)
	switch {
	case excludeDefaultValues == false:
		paths = append(paths, path.Join(chartPath, "/values.yaml"))
		fallthrough
	case cliPaths != "":
		paths = append(paths, strings.Split(cliPaths, ",")...)
	case len(paths) == 0:
		fallthrough
	default:
		log.Fatalf("Failed to find values file paths")

	}
	return paths
}
