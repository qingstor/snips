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
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yunify/snips/constants"
)

func TestLoadSpec(t *testing.T) {
	fixtures, err := filepath.Abs("fixtures")
	assert.Nil(t, err)

	specInfo, err := LoadSpec(
		path.Join(fixtures, "qingstor_sample", "api.json"),
		constants.SpecFormatSwaggerOpenAPI,
	)
	assert.Nil(t, err)
	assert.NotNil(t, specInfo)

	assert.Equal(t, "2016-01-06", specInfo.Data.Service.APIVersion)
	assert.Equal(t, "2016-01-06", specInfo.Data.Service.APIVersion)
	assert.Equal(t, "QingStor", specInfo.Data.Service.Name)
	assert.Equal(t, "QingStor", specInfo.Data.Service.Name)
	assert.Equal(t, 3, len(specInfo.Data.CustomizedTypes))
	assert.Equal(t, 3, len(specInfo.Data.CustomizedTypes))
}
