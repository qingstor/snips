// +-------------------------------------------------------------------------
// | Copyright (C) 2016 Yunify, Inc.
// +-------------------------------------------------------------------------
// | Licensed under the Apache License, Version 2.0 (the "License");
// | you may not use this work except in compliance with the License.
// | You may obtain a copy of the License in the LICENSE file, or at:
// |
// | http://www.apache.org/licenses/LICENSE-2.0
// |
// | Unless required by applicable law or agreed to in writing, software
// | distributed under the License is distributed on an "AS IS" BASIS,
// | WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// | See the License for the specific language governing permissions and
// | limitations under the License.
// +-------------------------------------------------------------------------

package templates

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// ManifestConfigurations holds the information of manifest file.
type ManifestConfigurations struct {
	Template        *templateConfigurations                `json:"template" yaml:"template"`
	Output          *outputConfigurations                  `json:"output" yaml:"output"`
	TemplateFiles   map[string]*templateFileConfigurations `json:"template_files" yaml:"template_files"`
	SupportingFiles []string                               `json:"supporting_files" yaml:"supporting_files"`
	MetaData        map[string]interface{}                 `json:"metadata" yaml:"metadata"`
	WordMaps        *wordMaps                              `json:"word_maps" yaml:"word_maps"`
}

type templateConfigurations struct {
	Format string `json:"format" yaml:"format"`
}

type outputConfigurations struct {
	FileNaming *fileNamingConfiguration `json:"file_naming" yaml:"file_naming"`
}

type templateFileConfigurations struct {
	FilePath         string                   `json:"file_path" yaml:"file_path"`
	Format           string                   `json:"format" yaml:"format"`
	OutputFileNaming *fileNamingConfiguration `json:"output_file_naming" yaml:"output_file_naming"`
}

type fileNamingConfiguration struct {
	Style     string `json:"style" yaml:"style"`
	Extension string `json:"extension" yaml:"extension"`
	Prefix    string `json:"prefix" yaml:"prefix"`
	Suffix    string `json:"suffix" yaml:"suffix"`
}

type wordMaps struct {
	CapitalizedToCapitalized map[string]string `json:"capitalized_to_capitalized" yaml:"capitalized_to_capitalized"`
	LowercaseToLowercase     map[string]string `json:"lowercase_to_lowercase" yaml:"lowercase_to_lowercase"`
	LowercaseToCapitalized   map[string]string `json:"lowercase_to_capitalized" yaml:"lowercase_to_capitalized"`
	Abbreviate               []string          `json:"abbreviate" yaml:"abbreviate"`
}

func loadManifestFromTemplateDirectory(path string) (*ManifestConfigurations, error) {
	manifest := newDefaultManifest()

	if content, err := ioutil.ReadFile(path + "/manifest.json"); err == nil {
		err = json.Unmarshal(content, manifest)
		if err != nil {
			return nil, err
		}
		return manifest, nil
	}

	if content, err := ioutil.ReadFile(path + "/manifest.yaml"); err == nil {
		err = yaml.Unmarshal(content, manifest)
		if err != nil {
			return nil, err
		}
		return manifest, nil
	}

	return nil, errors.New("Template manifest file not found in " + path + ".")
}

func newDefaultManifest() *ManifestConfigurations {
	return &ManifestConfigurations{
		Template: &templateConfigurations{
			Format: "Go",
		},
		Output: &outputConfigurations{
			FileNaming: &fileNamingConfiguration{
				Style: "snake_case",
			},
		},
		TemplateFiles: map[string]*templateFileConfigurations{
			"shared": {
				FilePath:         "shared.tmpl",
				OutputFileNaming: &fileNamingConfiguration{},
			},
			"service": {
				FilePath:         "service.tmpl",
				OutputFileNaming: &fileNamingConfiguration{},
			},
			"sub_service": {
				FilePath:         "sub_service.tmpl",
				OutputFileNaming: &fileNamingConfiguration{},
			},
			"types": {
				FilePath:         "types.tmpl",
				OutputFileNaming: &fileNamingConfiguration{},
			},
		},
	}
}
