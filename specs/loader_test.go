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
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadServices(t *testing.T) {
	filePath, err := filepath.Abs("fixtures")
	assert.Nil(t, err)

	serviceInfo, err := LoadServices(filePath, "Swagger-v2.0", "qingstor_sample")
	assert.Nil(t, err)
	assert.NotNil(t, serviceInfo)

	serviceInfo, err = LoadServices(filePath, "Swagger-v2.0", "QingStorSample")
	assert.Nil(t, err)
	assert.NotNil(t, serviceInfo)

	assert.Equal(t, "2016-01-06", serviceInfo.LatestAPIVersion.Spec.Data.Service.APIVersion)
	assert.Equal(t, "2016-01-06", serviceInfo.APIVersions["latest"].Spec.Data.Service.APIVersion)
	assert.Equal(t, "QingStor", serviceInfo.LatestAPIVersion.Spec.Data.Service.Name)
	assert.Equal(t, "QingStor", serviceInfo.APIVersions["latest"].Spec.Data.Service.Name)
	assert.Equal(t, 3, len(serviceInfo.LatestAPIVersion.Spec.Data.CustomizedTypes))
	assert.Equal(t, 3, len(serviceInfo.APIVersions["latest"].Spec.Data.CustomizedTypes))
}
