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
	"fmt"
	"strings"

	"github.com/go-openapi/spec"

	"github.com/yunify/snips/capsules"
)

func (s *Swagger) intermediateType(typeName, formatName string) string {
	intermediateTypesMap := map[string]string{
		"integer":                  "integer",
		"integer-int32":            "integer",
		"integer-int64":            "integer",
		"number":                   "float",
		"number-float":             "float",
		"number-double":            "double",
		"string":                   "string",
		"string-byte":              "base64",
		"string-binary":            "binary",
		"boolean":                  "boolean",
		"date":                     "timestamp",
		"string-date-time":         "timestamp",
		"string-date-time-rfc822":  "timestamp",
		"string-date-time-iso8601": "timestamp",
		"password":                 "string",

		"array":  "array",
		"object": "object",
		"map":    "map",
		"any":    "any",
	}

	compose := typeName
	if formatName != "" {
		compose += "-" + formatName
	}

	if intermediateTypesMap[compose] == "" {
		return compose
	}
	return intermediateTypesMap[compose]
}

func (s *Swagger) intermediateTypeOfTime(formatName string) string {
	intermediateTimeFormatMap := map[string]string{
		"date-time":        "ISO 8601",
		"date-time-rfc822": "RFC 822",
	}

	return intermediateTimeFormatMap[formatName]
}

func (s *Swagger) intermediateTypeOfSchema(schema *spec.Schema) string {
	if refTokens := schema.Ref.GetPointer().DecodedTokens(); len(refTokens) > 0 {
		return refTokens[len(refTokens)-1]
	}

	if strings.Join(schema.Type, "") != "" {
		return s.intermediateType(strings.Join(schema.Type, ""), schema.Format)
	}

	if schema.AdditionalProperties.Allows && schema.AdditionalProperties.Schema != nil {
		return s.intermediateTypeOfSchema(schema.AdditionalProperties.Schema)
	}

	return "any"
}

func (s *Swagger) parseEnum(enum []interface{}) []string {
	enums := []string{}
	for _, enumInterface := range enum {
		enums = append(enums, fmt.Sprintf("%v", enumInterface))
	}
	return enums
}

func (s *Swagger) parseSchema(schema *spec.Schema) *capsules.Property {
	targetSchema := schema
	targetType := ""
	targetExtraType := ""
	targetFormat := ""

	if refTokens := targetSchema.Ref.GetPointer().DecodedTokens(); len(refTokens) > 0 {
		targetType = "object"
		targetExtraType = s.intermediateTypeOfSchema(targetSchema)

	} else if targetSchema.AdditionalProperties != nil &&
		targetSchema.AdditionalProperties.Schema != nil &&
		targetSchema.AdditionalProperties.Allows {
		targetSchema = targetSchema.AdditionalProperties.Schema
		targetType = "map"
		targetExtraType = s.intermediateTypeOfSchema(targetSchema)

	} else {
		targetType = s.intermediateTypeOfSchema(targetSchema)
		if targetType == "array" {
			targetExtraType = s.intermediateTypeOfSchema(targetSchema.Items.Schema)
		}
		if targetType == "timestamp" {
			targetFormat = s.intermediateTypeOfTime(targetSchema.Format)
		}
	}

	defaultValue := ""
	if schema.Default != nil {
		defaultValue = fmt.Sprintf("%v", schema.Default)
	}

	properties := map[string]*capsules.Property{}
	for name, schema := range schema.SchemaProps.Properties {
		property := s.parseSchema(&schema)
		property.ID = name
		property.Name = name
		properties[name] = property
	}

	return &capsules.Property{
		Description: targetSchema.Description,
		Type:        targetType,
		ExtraType:   targetExtraType,
		Format:      targetFormat,
		Enum:        s.parseEnum(targetSchema.Enum),
		Default:     defaultValue,
		Properties:  properties,
	}
}

func (s *Swagger) parseParameter(
	parameter *spec.Parameter,
	parameters *map[string]spec.Parameter) *capsules.Property {

	targetParameter := parameter
	targetType := ""
	targetFormat := ""

	if refTokens := parameter.Ref.GetPointer().DecodedTokens(); len(refTokens) > 0 {
		current := (*parameters)[refTokens[len(refTokens)-1]]
		targetParameter = &current
	}

	targetType = s.intermediateType(targetParameter.Type, targetParameter.Format)
	if targetType == "timestamp" {
		targetFormat = s.intermediateTypeOfTime(targetParameter.Format)
	}

	defaultValue := ""
	if targetParameter.Default != nil {
		defaultValue = fmt.Sprintf("%v", targetParameter.Default)
	}

	return &capsules.Property{
		ID:          targetParameter.Name,
		Name:        targetParameter.Name,
		Description: targetParameter.Description,
		Type:        targetType,
		Format:      targetFormat,
		Enum:        s.parseEnum(targetParameter.Enum),
		Default:     defaultValue,
		IsRequired:  targetParameter.Required,
	}
}

func (s *Swagger) parseHeader(header *spec.Header) *capsules.Property {
	targetHeader := header
	targetType := ""
	targetFormat := ""

	targetType = s.intermediateType(targetHeader.Type, targetHeader.Format)
	if targetType == "timestamp" {
		targetFormat = s.intermediateTypeOfTime(targetHeader.Format)
	}

	defaultValue := ""
	if header.Default != nil {
		defaultValue = fmt.Sprintf("%v", header.Default)
	}

	return &capsules.Property{
		Description: targetHeader.Description,
		Type:        targetType,
		Format:      targetFormat,
		Enum:        s.parseEnum(targetHeader.Enum),
		Default:     defaultValue,
	}
}

func (s *Swagger) parseOperation(
	uri string, method string,
	specOperation *spec.Operation, swagger *spec.Swagger) *capsules.Operation {

	parsedURI := strings.Replace(uri, "?upload_id", "", -1)

	operation := &capsules.Operation{
		ID:          specOperation.ID,
		Name:        specOperation.Summary,
		Description: specOperation.Description,
		Request: &capsules.Request{
			Method: method,
			URI:    parsedURI,
			Params: &capsules.Property{
				ID:         specOperation.ID + "Input",
				Name:       specOperation.Summary + " Input",
				Properties: map[string]*capsules.Property{},
			},
			Headers: &capsules.Property{
				ID:         specOperation.ID + "Input",
				Name:       specOperation.Summary + " Input",
				Properties: map[string]*capsules.Property{},
			},
			Elements: &capsules.Property{
				ID:         specOperation.ID + "Input",
				Name:       specOperation.Summary + " Input",
				Properties: map[string]*capsules.Property{},
			},
			Body: &capsules.Property{},
		},
		Response: &capsules.Response{
			StatusCodes: map[int]*capsules.StatusCode{},
			Headers: &capsules.Property{
				ID:         specOperation.ID + "Output",
				Name:       specOperation.Summary + " Output",
				Properties: map[string]*capsules.Property{},
			},
			Elements: &capsules.Property{
				ID:         specOperation.ID + "Output",
				Name:       specOperation.Summary + " Output",
				Properties: map[string]*capsules.Property{},
			},
			Body: &capsules.Property{},
		},
	}

	if specOperation.ExternalDocs != nil {
		operation.DocumentationURL = specOperation.ExternalDocs.URL
	}

	for _, param := range specOperation.Parameters {
		switch param.In {
		case "query":
			property := s.parseParameter(&param, &swagger.Parameters)
			operation.Request.Params.Properties[param.Name] = property
		case "header":
			property := s.parseParameter(&param, &swagger.Parameters)
			operation.Request.Headers.Properties[param.Name] = property
		case "body":
			operation.Request.Body = s.parseSchema(param.Schema)
			if operation.Request.Body.Description == "" && param.Description != "" {
				operation.Request.Body.Description = param.Description
			}

			for name, schema := range param.Schema.Properties {
				property := s.parseSchema(&schema)
				property.Name = name
				property.ID = name
				operation.Request.Elements.Properties[name] = property
			}

			for _, schemaKey := range param.Schema.Required {
				if operation.Request.Elements.Properties[schemaKey] != nil {
					operation.Request.Elements.Properties[schemaKey].IsRequired = true
				}
			}
		}
	}

	for code, statusCode := range specOperation.Responses.ResponsesProps.StatusCodeResponses {
		operation.Response.StatusCodes[code] = &capsules.StatusCode{
			Description: statusCode.Description,
		}
	}

	successResponse := specOperation.Responses.ResponsesProps.StatusCodeResponses[200]

	for name, header := range successResponse.Headers {
		operation.Response.Headers.Properties[name] = s.parseHeader(&header)
		operation.Response.Headers.Properties[name].Name = name
		operation.Response.Headers.Properties[name].ID = name
	}

	if successResponse.Schema != nil {
		operation.Response.Body = s.parseSchema(successResponse.Schema)

		for name, schema := range successResponse.Schema.Properties {
			property := s.parseSchema(&schema)
			property.Name = name
			property.ID = name
			operation.Response.Elements.Properties[name] = property
		}

		for _, schemaKey := range successResponse.Schema.Required {
			if operation.Request.Elements.Properties[schemaKey] != nil {
				operation.Request.Elements.Properties[schemaKey].IsRequired = true
			}
		}
	}

	return operation
}
