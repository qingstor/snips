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

package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSnakeCase(t *testing.T) {
	assert.Equal(t, "latest_version", SnakeCase("lastest_version"))
	assert.Equal(t, "x_qs_date", SnakeCase("x-qs-date"))
	assert.Equal(t, "x_qs_date", SnakeCase("x-qs-date"))
	assert.Equal(t, "hello_world", SnakeCase("Hello World"))
	assert.Equal(t, "qingstor", SnakeCase("QingStor"))
	assert.Equal(t, "qingstor", SnakeCase("Qing Stor"))
}

func TestSnakeCaseToSnakeCase(t *testing.T) {
	assert.Equal(t, "latest_version", SnakeCaseToSnakeCase("lastest_version"))
	assert.Equal(t, "x_qs_date", SnakeCaseToSnakeCase("X-QS-Date", true))
	assert.Equal(t, "hello_world!", SnakeCaseToSnakeCase("Hello World!", true))
}

func TestSnakeCaseToCamelCase(t *testing.T) {
	assert.Equal(t, "CamelCase", SnakeCaseToCamelCase("camel_case"))
	assert.Equal(t, "BucketACL", SnakeCaseToCamelCase("bucket_acl"))
	assert.Equal(t, "AllowedOrigin", SnakeCaseToCamelCase("allowed_origin"))
	assert.Equal(t, "PartNumberMarker", SnakeCaseToCamelCase("part_number_marker"))
	assert.Equal(t, "CamelCase", SnakeCaseToCamelCase("camel_case"))
	assert.Equal(t, "CamelCase", SnakeCaseToCamelCase("camel_case"))
	assert.Equal(t, "CamelCase", SnakeCaseToCamelCase("camel_case"))
	assert.Equal(t, "XQSDate", SnakeCaseToCamelCase("X-QS-Date"))
}

func TestSnakeCaseToDashConnected(t *testing.T) {
	assert.Equal(t, "camel-case", SnakeCaseToDashConnected("camel_case"))
	assert.Equal(t, "bucket-acl", SnakeCaseToDashConnected("bucket_acl"))
}
