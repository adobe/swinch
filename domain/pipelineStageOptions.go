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

// Moniker is part of Stages
type Moniker struct {
	App string `yaml:"app" json:"app"`
}

// TrafficManagement is part of Stages
type TrafficManagement struct {
	Enabled bool `yaml:"enabled" json:"enabled"`
	Options struct {
		EnableTraffic bool          `yaml:"enableTraffic" json:"enableTraffic"`
		Services      []interface{} `yaml:"services" json:"services"`
	} `yaml:"options" json:"options"`
}

// LabelSelectors is part of Stages
type LabelSelectors struct {
	Selectors []struct {
		Key    string   `yaml:"key" json:"key"`
		Kind   string   `yaml:"kind" json:"kind"`
		Values []string `yaml:"values" json:"values"`
	} `yaml:"selectors" json:"selectors"`
}

// Options is part of Stages
type Options struct {
	Cascading bool `yaml:"cascading" json:"cascading"`
}
