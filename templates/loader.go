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
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/imdario/mergo"
)

// LoadTemplates read the template manifest and load all templates info.
func LoadTemplates(templateDirectory string) (map[string]*Template, *ManifestConfigurations, error) {
	templates := map[string]*Template{}

	manifest, err := loadManifestFromTemplateDirectory(templateDirectory)
	if err != nil {
		return templates, nil, err
	}

	for templateID := range manifest.TemplateFiles {
		serviceTemplate, err := loadTemplate(templateDirectory, templateID, manifest)
		if err == nil {
			templates[templateID] = serviceTemplate
		}
	}

	for index, supportingFileRelativePath := range manifest.SupportingFiles {
		supportingFile, err := loadSupportingFile(templateDirectory, supportingFileRelativePath)
		if err != nil {
			return templates, manifest, err
		}
		templates[fmt.Sprintf("%s-%d", "supporting", index)] = supportingFile
	}

	return templates, manifest, nil
}

func loadTemplate(templateDirectory, templateID string, manifest *ManifestConfigurations) (*Template, error) {
	templateConfiguration := manifest.TemplateFiles[templateID]
	if templateConfiguration == nil {
		return nil, fmt.Errorf(`configuration of template "%s" not found`, templateID)
	}

	templateFilePath := templateDirectory + "/" + templateConfiguration.FilePath
	templateFileContent, err := ioutil.ReadFile(templateFilePath)
	if err != nil {
		return nil, err
	}

	format := templateConfiguration.Format
	if format == "" {
		format = manifest.Template.Format
	}

	template := &Template{
		Filename:         filepath.Base(templateFilePath),
		FileDirectory:    filepath.Dir(templateConfiguration.FilePath),
		FilePath:         templateFilePath,
		FileContent:      string(templateFileContent),
		Format:           format,
		IsNeedGenerate:   true,
		ID:               templateID,
		OutputFileNaming: templateConfiguration.OutputFileNaming,
	}
	if manifest.Output.FileNaming != nil {
		mergo.Merge(template.OutputFileNaming, *manifest.Output.FileNaming)
	}

	return template, nil
}

func loadSupportingFile(templateDirectory, relativeFilePath string) (*Template, error) {
	filePath := templateDirectory + "/" + relativeFilePath
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return &Template{
		Filename:       filepath.Base(filePath),
		FileDirectory:  filepath.Dir(relativeFilePath),
		FilePath:       filePath,
		FileContent:    string(fileContent),
		IsNeedGenerate: false,
	}, nil
}
