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
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadTemplates_0(t *testing.T) {
	fullPath, err := filepath.Abs("fixtures/template_0")
	assert.Nil(t, err)

	templates, manifest, err := LoadTemplates(fullPath)
	assert.Nil(t, err)

	assert.Equal(t, 4, len(templates))
	assert.Equal(t, "{{$service := .Data.Service}}\n{{camelCase $service.Name}}\n", string(templates["service"].FileContent))
	assert.Equal(t, "sub service template\n", string(templates["sub_service"].FileContent))
	assert.Equal(t, "types template\n", string(templates["types"].FileContent))
	assert.Equal(t, "shared\n", string(templates["shared"].FileContent))
	assert.Equal(t, "Mustache", manifest.Template.Format)
}

func TestLoadTemplates_1(t *testing.T) {
	fullPath, err := filepath.Abs("fixtures/template_1")
	assert.Nil(t, err)

	templates, manifest, err := LoadTemplates(fullPath)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(templates))
	assert.Equal(t, "service test template\n", string(templates["service"].FileContent))
	assert.Equal(t, "Go", manifest.Template.Format)
}

func TestLoadTemplates_2(t *testing.T) {
	fullPath, err := filepath.Abs("fixtures/template_2")
	assert.Nil(t, err)

	templates, manifest, err := LoadTemplates(fullPath)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(templates))
	assert.Equal(t, "types template\n", string(templates["types"].FileContent))
	assert.Equal(t, "camel_case", string(templates["types"].OutputFileNaming.Style))
	assert.Equal(t, ".any", string(templates["types"].OutputFileNaming.Extension))
	assert.Equal(t, "qs_", string(templates["types"].OutputFileNaming.Prefix))
	assert.Equal(t, "types", string(templates["types"].ID))
	assert.Equal(t, "snake_case", manifest.Output.FileNaming.Style)
}

func TestLoadTemplates_3(t *testing.T) {
	fullPath, err := filepath.Abs("fixtures/template_3")
	assert.Nil(t, err)

	templates, manifest, err := LoadTemplates(fullPath)
	assert.Nil(t, err)

	assert.Equal(t, 2, len(templates))
	assert.Equal(t, "camel_case", manifest.Output.FileNaming.Style)
	assert.Equal(t, ".rb", manifest.Output.FileNaming.Extension)
}
