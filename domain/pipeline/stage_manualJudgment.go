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

type ManualJudgment struct {
	FailPipeline   bool          `yaml:"failPipeline" json:"failPipeline"`
	IsNew          bool          `yaml:"isNew" json:"isNew"`
	JudgmentInputs []interface{} `yaml:"judgmentInputs" json:"judgmentInputs"`
	Notifications  []interface{} `yaml:"notifications" json:"notifications"`
}