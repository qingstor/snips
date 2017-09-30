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

func TestSwagger_Parse(t *testing.T) {
	filePath, err := filepath.Abs("fixtures/qingstor_sample/api.json")
	assert.Nil(t, err)
	swagger := &Swagger{
		FilePath: filePath,
	}

	err = swagger.Parse("v2.0")
	assert.Nil(t, err)

	assert.Equal(t, "QingStor", swagger.Data.Service.Name)
	assert.Equal(t, "Bucket", swagger.Data.SubServices["Bucket"].Name)
	assert.Equal(t, 7, len(swagger.Data.SubServices["Bucket"].Operations))

	owner := swagger.Data.SubServices["Bucket"].Operations["ListObjects"].Responses[200].Elements.Properties["owner"]
	assert.Equal(t, "object", owner.Type)
	assert.Equal(t, "owner", owner.ExtraType)

	listBuckets := swagger.Data.Service.Operations["ListBuckets"]
	location := listBuckets.Request.Headers.Properties["Location"]
	assert.Equal(t, "Location", location.Name)
	assert.Equal(t, "string", location.Type)

	bucket := swagger.Data.CustomizedTypes["bucket"]
	assert.Equal(t, "bucket", bucket.Name)
	assert.Equal(t, "Bucket", bucket.Description)
	assert.Equal(t, "object", bucket.Type)
}
