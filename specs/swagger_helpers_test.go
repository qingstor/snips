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

	"github.com/go-openapi/loads"

	"github.com/stretchr/testify/assert"
)

func TestSwagger_getIntermediateType(t *testing.T) {
	swagger := Swagger{}

	assert.Equal(t, "integer", swagger.intermediateType("integer", "int32"))
	assert.Equal(t, "integer", swagger.intermediateType("integer", "int64"))
	assert.Equal(t, "binary", swagger.intermediateType("string", "binary"))
	assert.Equal(t, "timestamp", swagger.intermediateType("string", "date-time-rfc822"))
	assert.Equal(t, "string", swagger.intermediateType("password", ""))
}

func TestSwagger_getIntermediateTypeOfTime(t *testing.T) {
	swagger := Swagger{}

	assert.Equal(t, "ISO 8601", swagger.intermediateTypeOfTime("date-time"))
	assert.Equal(t, "RFC 822", swagger.intermediateTypeOfTime("date-time-rfc822"))
}

func TestSwagger_getIntermediateTypeOfSchema(t *testing.T) {
	filePath, err := filepath.Abs("fixtures/qingstor_sample/2016-01-06/swagger/api_v2.0.json")
	assert.Nil(t, err)

	document, err := loads.Spec(filePath)
	assert.Nil(t, err)

	document, err = document.Expanded()
	assert.Nil(t, err)

	swagger := Swagger{}
	keys := document.Spec().Paths.Paths["/{bucketName}"].Get.Responses.StatusCodeResponses[200].Schema.Properties["keys"].Items.Schema
	assert.Equal(t, "key", swagger.intermediateTypeOfSchema(keys))
}

func TestSwagger_parseEnum(t *testing.T) {
	filePath, err := filepath.Abs("fixtures/qingstor_sample/2016-01-06/swagger/api_v2.0.json")
	assert.Nil(t, err)

	document, err := loads.Spec(filePath)
	assert.Nil(t, err)

	document, err = document.Expanded()
	assert.Nil(t, err)

	swagger := Swagger{}
	status := document.Spec().Paths.Paths["/{bucketName}?stats"].Get.Responses.StatusCodeResponses[200].Schema.Properties["status"]
	enum := swagger.parseEnum(status.Enum)
	assert.Equal(t, "active", enum[0])
	assert.Equal(t, "suspended", enum[1])
}

func TestSwagger_parseSchema(t *testing.T) {
	filePath, err := filepath.Abs("fixtures/qingstor_sample/2016-01-06/swagger/api_v2.0.json")
	assert.Nil(t, err)

	document, err := loads.Spec(filePath)
	assert.Nil(t, err)

	document, err = document.Expanded()
	assert.Nil(t, err)

	swagger := Swagger{}

	keys := document.Spec().Paths.Paths["/{bucketName}"].Get.Responses.StatusCodeResponses[200].Schema.Properties["keys"].Items.Schema
	property := swagger.parseSchema(keys)
	assert.Equal(t, "object", property.Type)
	assert.Equal(t, "key", property.ExtraType)
	assert.Equal(t, false, property.IsRequired)

	body := document.Spec().Paths.Paths["/{bucketName}"].Head.Responses.StatusCodeResponses[200].Schema
	bodyProperty := swagger.parseSchema(body)
	assert.Equal(t, "binary", bodyProperty.Type)
	assert.Equal(t, "", bodyProperty.ExtraType)
	assert.Equal(t, "This is response body", bodyProperty.Description)
	assert.Equal(t, false, property.IsRequired)
}

func TestSwagger_parseParameter(t *testing.T) {
	filePath, err := filepath.Abs("fixtures/qingstor_sample/2016-01-06/swagger/api_v2.0.json")
	assert.Nil(t, err)

	document, err := loads.Spec(filePath)
	assert.Nil(t, err)

	document, err = document.Expanded()
	assert.Nil(t, err)

	swagger := Swagger{}
	delimiter := document.Spec().Paths.Paths["/{bucketName}"].Get.Parameters[1]
	property := swagger.parseParameter(&delimiter, &document.Spec().Parameters)
	assert.Equal(t, "string", property.Type)
	assert.Equal(t, "", property.ExtraType)
	assert.Equal(t, false, property.IsRequired)
}

func TestSwagger_parseHeader(t *testing.T) {
	filePath, err := filepath.Abs("fixtures/qingstor_sample/2016-01-06/swagger/api_v2.0.json")
	assert.Nil(t, err)

	document, err := loads.Spec(filePath)
	assert.Nil(t, err)

	document, err = document.Expanded()
	assert.Nil(t, err)

	swagger := Swagger{}
	location := document.Spec().Paths.Paths["/"].Get.Responses.StatusCodeResponses[200].Headers["fake"]
	property := swagger.parseHeader(&location)
	assert.Equal(t, "string", property.Type)
	assert.Equal(t, "", property.ExtraType)
	assert.Equal(t, false, property.IsRequired)
}

func TestSwagger_parseOperation(t *testing.T) {
	filePath, err := filepath.Abs("fixtures/qingstor_sample/2016-01-06/swagger/api_v2.0.json")
	assert.Nil(t, err)

	document, err := loads.Spec(filePath)
	assert.Nil(t, err)

	document, err = document.Expanded()
	assert.Nil(t, err)

	swagger := Swagger{}
	operation := document.Spec().Paths.Paths["/{bucketName}"].Get
	parsedOperation := swagger.parseOperation("/{bucketName}", "GET", operation, document.Spec())
	assert.Equal(t, "ListObjects", parsedOperation.ID)
	assert.Equal(t, "GET Bucket (List Objects)", parsedOperation.Name)
	assert.Equal(t, 4, len(parsedOperation.Request.Params.Properties))
	assert.Equal(t, 9, len(parsedOperation.Response.Elements.Properties))
}
