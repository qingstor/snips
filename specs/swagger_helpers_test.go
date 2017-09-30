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

	"github.com/yunify/snips/capsules"
)

func TestSwagger_getIntermediateType(t *testing.T) {
	swagger := Swagger{}

	assert.Equal(t, "integer", swagger.intermediateType("integer", "int32"))
	assert.Equal(t, "long", swagger.intermediateType("integer", "int64"))
	assert.Equal(t, "binary", swagger.intermediateType("string", "binary"))
	assert.Equal(t, "timestamp", swagger.intermediateType("string", "date-time-rfc822"))
	assert.Equal(t, "string", swagger.intermediateType("string", "password"))
}

func TestSwagger_getIntermediateTypeOfTime(t *testing.T) {
	swagger := Swagger{}

	assert.Equal(t, "ISO 8601", swagger.intermediateTypeOfTime("date-time"))
	assert.Equal(t, "RFC 822", swagger.intermediateTypeOfTime("date-time-rfc822"))
}

func TestSwagger_getIntermediateTypeOfSchema(t *testing.T) {
	filePath, err := filepath.Abs("fixtures/qingstor_sample/api.json")
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
	filePath, err := filepath.Abs("fixtures/qingstor_sample/api.json")
	assert.Nil(t, err)

	document, err := loads.Spec(filePath)
	assert.Nil(t, err)

	document, err = document.Expanded()
	assert.Nil(t, err)

	swagger := Swagger{}
	status := document.Spec().Paths.Paths["/{bucketName}?stats&csv={csvArrayTest}"].Get.Responses.StatusCodeResponses[200].Schema.Properties["status"]
	enum := swagger.parseEnum(status.Enum)
	assert.Equal(t, "active", enum[0])
	assert.Equal(t, "suspended", enum[1])
}

func TestSwagger_parseSchema(t *testing.T) {
	filePath, err := filepath.Abs("fixtures/qingstor_sample/api.json")
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
	filePath, err := filepath.Abs("fixtures/qingstor_sample/api.json")
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

	csv := document.Spec().Paths.Paths["/{bucketName}?stats&csv={csvArrayTest}"].Parameters[2]
	property = swagger.parseParameter(&csv, &document.Spec().Parameters)
	assert.Equal(t, "array", property.Type)
	assert.Equal(t, "csv", property.CollectionFormat)
	assert.Equal(t, "string", property.ExtraType)
	assert.Equal(t, false, property.IsRequired)

	numberValidation := document.Spec().Paths.Paths["/{bucketName}?validations&number={validationsNumberTest}&string={validationsStringTest}"].Parameters[2]
	property = swagger.parseParameter(&numberValidation, &document.Spec().Parameters)
	assert.Equal(t, 10, int(*property.Maximum))
	assert.Equal(t, 1, int(*property.Minimum))

	stringValidation := document.Spec().Paths.Paths["/{bucketName}?validations&number={validationsNumberTest}&string={validationsStringTest}"].Parameters[3]
	property = swagger.parseParameter(&stringValidation, &document.Spec().Parameters)
	assert.Equal(t, 10, int(*property.MaxLength))
	assert.Equal(t, 1, int(*property.MinLength))

	formDataTest := document.Spec().Paths.Paths["/{bucketName}?validations&number={validationsNumberTest}&string={validationsStringTest}"].Post.Parameters[0]
	property = swagger.parseParameter(&formDataTest, &document.Spec().Parameters)
	assert.Equal(t, "string", property.Type)
	assert.Equal(t, "", property.ExtraType)
	assert.Equal(t, false, property.IsRequired)
}

func TestSwagger_parseHeader(t *testing.T) {
	filePath, err := filepath.Abs("fixtures/qingstor_sample/api.json")
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
	filePath, err := filepath.Abs("fixtures/qingstor_sample/api.json")
	assert.Nil(t, err)

	document, err := loads.Spec(filePath)
	assert.Nil(t, err)

	document, err = document.Expanded()
	assert.Nil(t, err)

	swagger := Swagger{}
	operation := document.Spec().Paths.Paths["/{bucketName}"].Get
	property := &capsules.Property{
		Properties: map[string]*capsules.Property{
			"bucketName": {},
		},
	}
	parsedOperation := swagger.parseOperation("/{bucketName}", "GET", property, operation, document.Spec())
	assert.Equal(t, "ListObjects", parsedOperation.ID)
	assert.Equal(t, "GET Bucket (List Objects)", parsedOperation.Name)
	assert.Equal(t, 1, len(parsedOperation.Request.Properties.Properties))
	assert.Equal(t, 4, len(parsedOperation.Request.Query.Properties))
	assert.Equal(t, 9, len(parsedOperation.Responses[200].Elements.Properties))
}
