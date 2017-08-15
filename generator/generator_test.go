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

package generator

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yunify/snips/capsules"
	"github.com/yunify/snips/constants"
	"github.com/yunify/snips/specs"
	"github.com/yunify/snips/templates"
)

func TestGenerator(t *testing.T) {
	templatePath, err := filepath.Abs("../templates/fixtures/template_0")
	assert.Nil(t, err)
	loadedTemplates, manifest, err := templates.LoadTemplates(templatePath)
	assert.Nil(t, err)
	assert.Equal(t, 4, len(loadedTemplates))
	assert.Equal(t, "service", loadedTemplates["service"].ID)
	assert.Equal(t, "Mustache", manifest.Template.Format)

	fixtures, err := filepath.Abs("../specs/fixtures")
	assert.Nil(t, err)
	specInfo, err := specs.LoadSpec(
		path.Join(fixtures, "qingstor_sample", "api.json"),
		constants.SpecFormatSwaggerOpenAPI,
	)
	assert.Nil(t, err)
	assert.NotNil(t, specInfo)

	capsule := &capsules.BaseCapsule{
		CapsulePowder: &capsules.CapsulePowder{
			Data: specInfo.Data,
		},
	}

	codeGenerator := New()
	codeGenerator.Set(capsule, loadedTemplates["service"])

	err = codeGenerator.Render()
	assert.NotNil(t, err)
	assert.Equal(t, "Template format not supported: \"Mustache\"", err.Error())

	loadedTemplates["service"].Format = "Go"
	err = codeGenerator.Render()
	assert.Nil(t, err)
	assert.Equal(t, "\nQingStor\n", codeGenerator.code)

	loadedTemplates["service"].UpdateOutputFilename("hello", "camel_case")
	err = loadedTemplates["service"].UpdateOutputFilePath("./test")
	assert.Nil(t, err)
	err = codeGenerator.Write()
	assert.Nil(t, err)
	fileContent, err := ioutil.ReadFile("./test/Hello")
	assert.Nil(t, err)
	assert.Equal(t, "\nQingStor\n", string(fileContent))

	err = os.RemoveAll("./test")
	assert.Nil(t, err)

	codeGenerator.Clear()
	assert.Nil(t, codeGenerator.capsule)
	assert.Nil(t, codeGenerator.template)
	assert.Equal(t, "", codeGenerator.code)
}
