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

func TestLoadManifestFromTemplateDirectory_0(t *testing.T) {
	fullPath, err := filepath.Abs("fixtures/template_0")
	assert.Nil(t, err)

	manifest, err := loadManifestFromTemplateDirectory(fullPath)
	assert.Nil(t, err)

	assert.Equal(t, "Mustache", manifest.Template.Format)
	assert.Equal(t, "service.tmpl", manifest.TemplateFiles["service"].FilePath)
	assert.Equal(t, "shared.tmpl", manifest.TemplateFiles["shared"].FilePath)
}

func TestLoadManifestFromTemplateDirectory_1(t *testing.T) {
	fullPath, err := filepath.Abs("fixtures/template_1")
	assert.Nil(t, err)

	manifest, err := loadManifestFromTemplateDirectory(fullPath)
	assert.Nil(t, err)

	assert.Equal(t, "service_test.tmpl", manifest.TemplateFiles["service"].FilePath)
	assert.Equal(t, "qs_", manifest.TemplateFiles["service"].OutputFileNaming.Prefix)
	assert.Equal(t, "_service", manifest.TemplateFiles["service"].OutputFileNaming.Suffix)
	assert.Equal(t, "Mustache", manifest.TemplateFiles["service"].Format)
	assert.Equal(t, "sub_service.tmpl", manifest.TemplateFiles["sub_service"].FilePath)
}

func TestLoadManifestFromTemplateDirectory_2(t *testing.T) {
	fullPath, err := filepath.Abs("fixtures/template_2")
	assert.Nil(t, err)

	manifest, err := loadManifestFromTemplateDirectory(fullPath)
	assert.Nil(t, err)

	assert.Equal(t, "types.tmpl", manifest.TemplateFiles["types"].FilePath)
	assert.Equal(t, "qs_", manifest.TemplateFiles["types"].OutputFileNaming.Prefix)
	assert.Equal(t, "", manifest.TemplateFiles["types"].OutputFileNaming.Suffix)
	assert.Equal(t, "ignore.tmpl", manifest.TemplateFiles["ignore"].FilePath)
}

func TestLoadManifestFromTemplateDirectory_3(t *testing.T) {
	fullPath, err := filepath.Abs("fixtures/template_3")
	assert.Nil(t, err)

	manifest, err := loadManifestFromTemplateDirectory(fullPath)
	assert.Nil(t, err)

	assert.Equal(t, "camel_case", manifest.Output.FileNaming.Style)
	assert.Equal(t, ".rb", manifest.Output.FileNaming.Extension)
	assert.Equal(t, 2, len(manifest.SupportingFiles))
}

func TestNewDefaultManifest(t *testing.T) {
	manifest := newDefaultManifest()

	assert.Equal(t, "Go", manifest.Template.Format)
	assert.Equal(t, "snake_case", manifest.Output.FileNaming.Style)
	assert.Equal(t, "service.tmpl", manifest.TemplateFiles["service"].FilePath)
	assert.Equal(t, "sub_service.tmpl", manifest.TemplateFiles["sub_service"].FilePath)
	assert.Equal(t, "types.tmpl", manifest.TemplateFiles["types"].FilePath)
	assert.Nil(t, manifest.TemplateFiles["other"])
}
