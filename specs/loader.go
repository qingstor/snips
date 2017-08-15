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

package specs

import (
	"errors"
	"io/ioutil"
	"path"

	"github.com/yunify/snips/capsules"
	"github.com/yunify/snips/constants"
)

// A Spec holds the information of an API spec file.
type Spec struct {
	Filename    string
	FilePath    string
	FileContent string

	Format string

	Data *capsules.Data
}

// LoadSpec loads a specification.
func LoadSpec(filePath, format string) (s *Spec, err error) {
	switch format {
	case constants.SpecFormatSwagger, constants.SpecFormatSwaggerV2,
		constants.SpecFormatSwaggerOpenAPI, constants.SpecFormatSwaggerOpenAPIV2:
		format = constants.SpecFormatSwaggerOpenAPIV2
	default:
		err = errors.New("Spec format not supported: " + format)
		return
	}

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return
	}
	swagger := Swagger{
		FilePath: filePath,
		Data:     &capsules.Data{},
	}
	err = swagger.Parse("v2.0")
	if err != nil {
		return
	}

	s = &Spec{
		Filename:    path.Base(filePath),
		FilePath:    filePath,
		FileContent: string(content),
		Format:      format,
		Data:        swagger.Data,
	}
	return
}
