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

package datastore

import (
	"bytes"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
	"path"
	"path/filepath"
)

const (
	FilePerm = 0775
)

type Datastore struct {
}

// LoadYAMLFiles receives a folder path, reads all yaml files, merges them in a buffer and returns it
func (d Datastore) LoadYAMLFiles(path string) *bytes.Buffer {
	yamlFilesBuffer := new(bytes.Buffer)

	switch location, err := os.Stat(path); {
	case err != nil:
		log.Fatalf("error: %v", err)
	case location.IsDir() == true:
		files, err := os.ReadDir(path)
		if err != nil {
			log.Fatalf("error: %v", err)
		}

		for _, file := range files {
			if file.IsDir() {
				continue
			} else {
				if filepath.Ext(file.Name()) == ".yaml" {
					yamlFilesBuffer.Write(d.ReadFile(path + "/" + file.Name()))
				} else {
					continue
				}
			}
		}
	case location.IsDir() == false:
		if filepath.Ext(path) == ".yaml" {
			yamlFilesBuffer.Write(d.ReadFile(path))
		} else {
			log.Errorf("Please provide an yaml file")
		}
	}

	return yamlFilesBuffer
}

func (d Datastore) WriteJSON(data interface{}, outputPath string) {
	byteData := d.MarshalJSON(data)
	d.WriteFile(outputPath, byteData, FilePerm)
}

func (d Datastore) WriteJSONTmp(data interface{}) (filePath string) {
	byteData := d.MarshalJSON(data)
	return d.writeTmpFile(byteData)
}

func (d Datastore) WriteYAML(data interface{}, outputPath string) {
	byteData := d.MarshalYAML(data)
	d.WriteFile(outputPath, byteData, FilePerm)
}

func (d *Datastore) MarshalJSON(data interface{}) []byte {
	byteData, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Fatalf("Failed to marshal JSON:  %v", err)
	}
	return byteData
}

func (d *Datastore) MarshalYAML(data interface{}) []byte {
	byteData := new(bytes.Buffer)
	yamlEncoder := yaml.NewEncoder(byteData)
	yamlEncoder.SetIndent(2)
	err := yamlEncoder.Encode(&data)
	if err != nil {
		log.Fatalf("Error marshal YAML:  %v", err)
	}
	return byteData.Bytes()
}

// Utils

func (d *Datastore) CreateTmpFolder() string {
	w, err := os.MkdirTemp("", "tempfolder")
	if err != nil {
		log.Fatalf("Failed to create the temp folder: %v", err)
	}
	return w
}

func (d *Datastore) writeTmpFile(byteData []byte) string {
	w, err := os.CreateTemp("" /* /tmp dir. */, "tempfile")
	if err != nil {
		log.Fatalf("Failed to create the temp file: %v", err)
	}
	bytesWrite, errWrite := w.Write(byteData)
	if errWrite != nil || bytesWrite == 0 {
		log.Fatalf("Failed to write tmp file: %v", errWrite)
	}
	return w.Name()
}

func (d Datastore) ReadFile(filePath string) []byte {
	byteData, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error reading file in path: %v, %v", filePath, err)
	}
	return byteData
}

func (d Datastore) WriteFile(outputPath string, byteData []byte, perm int) {
	filePath := path.Join(outputPath)
	err := os.WriteFile(filePath, byteData, os.FileMode(perm))
	if err != nil {
		log.Fatalf("Failed to write:  %v in path: %v", err, filePath)
		log.Debugf("Writing: %v", byteData)
	}
}

func (d Datastore) Mkdir(path string, perm int) {
	if _, errStat := os.Stat(path); os.IsNotExist(errStat) {
		err := os.MkdirAll(path, os.FileMode(perm))
		if err != nil {
			log.Fatalf("Error mkdir:  %v", err)
		}
	}
}

func (d Datastore) FileExists(path string) bool {
	if _, errStat := os.Stat(path); os.IsNotExist(errStat) {
		return false
	} else {
		return true
	}
}
