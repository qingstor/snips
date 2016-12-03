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
	"github.com/yunify/snips/utils"
	"path/filepath"
)

// A Template holds the information of template file.
type Template struct {
	Filename      string
	FileDirectory string
	FilePath      string
	FileContent   string

	Format         string
	IsNeedGenerate bool

	ID string // service, sub_service or types

	OutputFilename   string
	OutputFilePath   string
	OutputFileNaming *fileNamingConfiguration
}

// UpdateOutputFilename updates the output filename according to namingStyle.
func (t *Template) UpdateOutputFilename(filename string, namingStyle string) {
	switch namingStyle {
	case "snake_case":
		t.OutputFilename = utils.SnakeCase(filename)
	case "camel_case":
		t.OutputFilename = utils.CamelCase(filename)
	}
}

// UpdateOutputFilePath updates the absolute path of output file.
func (t *Template) UpdateOutputFilePath(parentPath string) error {
	relativePath := "" +
		t.FileDirectory + "/" +
		t.OutputFileNaming.Prefix +
		t.OutputFilename +
		t.OutputFileNaming.Suffix +
		t.OutputFileNaming.Extension

	targetFilePath, err := filepath.Abs(parentPath + "/" + relativePath)
	if err != nil {
		return err
	}

	t.OutputFilePath = targetFilePath
	return nil
}
