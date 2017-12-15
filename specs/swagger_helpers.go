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
	"strconv"
	"strings"

	"github.com/go-openapi/spec"
	"github.com/imdario/mergo"

	"github.com/yunify/snips/capsules"
)

func (s *Swagger) intermediateType(typeName, formatName string) string {
	intermediateTypesMap := map[string]string{
		"boolean":                  "boolean",
		"date":                     "timestamp",
		"integer":                  "integer",
		"integer-int32":            "integer",
		"integer-uint32":           "unsigned-integer",
		"integer-int64":            "long",
		"integer-uint64":           "unsigned-long",
		"number":                   "float",
		"number-float":             "float",
		"number-double":            "double",
		"string":                   "string",
		"string-byte":              "base64",
		"string-binary":            "binary",
		"string-boolean":           "boolean",
		"string-email":             "email",
		"string-password":          "string",
		"string-date-time":         "timestamp",
		"string-date-time-rfc822":  "timestamp",
		"string-date-time-iso8601": "timestamp",

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

func (s *Swagger) parseSchema(name string, schema *spec.Schema) *capsules.Property {
	targetSchema := schema
	targetType := ""
	targetExtraType := ""
	targetFormat := ""

	if refTokens := targetSchema.Ref.GetPointer().DecodedTokens(); len(refTokens) > 0 {
		targetType = strings.Join(schema.Type, "")
		switch targetType {
		case "array":
			targetExtraType = s.intermediateTypeOfSchema(targetSchema.Items.Schema)
		default:
			targetType = "object"
			targetExtraType = s.intermediateTypeOfSchema(targetSchema)
		}
	} else if targetSchema.AdditionalProperties != nil &&
		targetSchema.AdditionalProperties.Schema != nil &&
		targetSchema.AdditionalProperties.Allows {
		targetSchema = targetSchema.AdditionalProperties.Schema
		targetType = "map"
		targetExtraType = s.intermediateTypeOfSchema(targetSchema)
	} else {
		targetType = s.intermediateTypeOfSchema(targetSchema)
		switch targetType {
		case "array":
			targetExtraType = s.intermediateTypeOfSchema(targetSchema.Items.Schema)
		case "timestamp":
			targetFormat = s.intermediateTypeOfTime(targetSchema.Format)
		}
	}

	defaultValue := ""
	if schema.Default != nil {
		defaultValue = fmt.Sprintf("%v", schema.Default)
	}

	properties := map[string]*capsules.Property{}
	for name, schema := range schema.SchemaProps.Properties {
		property := s.parseSchema(name, &schema)
		property.ID = name
		property.Name = name
		properties[name] = property
	}

	return &capsules.Property{
		Description: targetSchema.Description,
		Type:        targetType,
		ExtraType:   targetExtraType,
		Format:      targetFormat,
		Default:     defaultValue,
		Properties:  properties,
		CommonValidations: capsules.CommonValidations{
			Enum:             s.parseEnum(targetSchema.Enum),
			Maximum:          targetSchema.Maximum,
			Minimum:          targetSchema.Minimum,
			MaxLength:        targetSchema.MaxLength,
			MinLength:        targetSchema.MinLength,
			ExclusiveMaximum: targetSchema.ExclusiveMaximum,
			ExclusiveMinimum: targetSchema.ExclusiveMinimum,
			Pattern:          targetSchema.Pattern,
			MaxItems:         targetSchema.MaxItems,
			MinItems:         targetSchema.MinItems,
			UniqueItems:      targetSchema.UniqueItems,
			MultipleOf:       targetSchema.MultipleOf,
		},
	}
}

func (s *Swagger) parseParameter(
	parameter *spec.Parameter,
	parameters *map[string]spec.Parameter) *capsules.Property {

	targetParameter := parameter
	targetType := ""
	targetFormat := ""
	targetExtraType := ""
	targetCollectionFormat := ""

	if refTokens := parameter.Ref.GetPointer().DecodedTokens(); len(refTokens) > 0 {
		current := (*parameters)[refTokens[len(refTokens)-1]]
		targetParameter = &current
	}

	targetType = s.intermediateType(targetParameter.Type, targetParameter.Format)
	if targetType == "timestamp" {
		targetFormat = s.intermediateTypeOfTime(targetParameter.Format)
	}
	if targetType == "array" {
		targetExtraType = s.intermediateType(
			targetParameter.Items.Type, targetParameter.Items.Format,
		)
		targetCollectionFormat = targetParameter.CollectionFormat
	}

	defaultValue := ""
	if targetParameter.Default != nil {
		defaultValue = fmt.Sprintf("%v", targetParameter.Default)
	}

	return &capsules.Property{
		ID:               targetParameter.Name,
		Name:             targetParameter.Name,
		Description:      targetParameter.Description,
		Type:             targetType,
		ExtraType:        targetExtraType,
		Format:           targetFormat,
		CollectionFormat: targetCollectionFormat,
		Default:          defaultValue,
		IsRequired:       targetParameter.Required,
		CommonValidations: capsules.CommonValidations{
			Enum:             s.parseEnum(targetParameter.Enum),
			Maximum:          targetParameter.Maximum,
			Minimum:          targetParameter.Minimum,
			MaxLength:        targetParameter.MaxLength,
			MinLength:        targetParameter.MinLength,
			ExclusiveMaximum: targetParameter.ExclusiveMaximum,
			ExclusiveMinimum: targetParameter.ExclusiveMinimum,
			Pattern:          targetParameter.Pattern,
			MaxItems:         targetParameter.MaxItems,
			MinItems:         targetParameter.MinItems,
			UniqueItems:      targetParameter.UniqueItems,
			MultipleOf:       targetParameter.MultipleOf,
		},
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
		Default:     defaultValue,
		CommonValidations: capsules.CommonValidations{
			Enum: s.parseEnum(targetHeader.Enum),
		},
	}
}

func (s *Swagger) parseOperation(
	uri string, method string, property *capsules.Property,
	specOperation *spec.Operation, swagger *spec.Swagger) *capsules.Operation {

	parsedURI := strings.Replace(uri, "?upload_id", "", -1)

	// If basePath is set, we should add it to Request URI.
	// OpenAPI Specification formulates that basePath must start with "/",
	// so we use len to judge should we need to add it.
	// ref: https://github.com/OAI/OpenAPI-Specification/blob/master/versions/2.0.md#fixed-fields
	if len(swagger.BasePath) > 1 {
		parsedURI = swagger.BasePath + parsedURI
	}

	operation := &capsules.Operation{
		ID:          specOperation.ID,
		Name:        specOperation.Summary,
		Consumes:    specOperation.Consumes,
		Description: specOperation.Description,
		Request: &capsules.Request{
			Method: method,
			Path:   parsedURI,
			Properties: &capsules.Property{
				ID:         specOperation.ID + "Input",
				Name:       specOperation.Summary + " Input",
				Properties: map[string]*capsules.Property{},
			},
			Query: &capsules.Property{
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
			FormData: &capsules.Property{
				ID:         specOperation.ID + "Input",
				Name:       specOperation.Summary + " Input",
				Properties: map[string]*capsules.Property{},
			},
			Body: &capsules.Property{},
		},
		Responses: make(map[int]*capsules.Response),
	}

	if specOperation.ExternalDocs != nil {
		operation.DocumentationURL = specOperation.ExternalDocs.URL
	}

	// Fill path params into request params
	mergo.Merge(operation.Request.Properties, property)

	for _, param := range specOperation.Parameters {
		switch param.In {
		case "path":
			property := s.parseParameter(&param, &swagger.Parameters)
			operation.Request.Properties.Properties[param.Name] = property
		case "query":
			property := s.parseParameter(&param, &swagger.Parameters)
			operation.Request.Query.Properties[param.Name] = property
		case "header":
			property := s.parseParameter(&param, &swagger.Parameters)
			operation.Request.Headers.Properties[param.Name] = property
		case "formData":
			property := s.parseParameter(&param, &swagger.Parameters)
			operation.Request.FormData.Properties[param.Name] = property
		case "body":
			operation.Request.Body = s.parseSchema(param.Name, param.Schema)
			if operation.Request.Body.Description == "" && param.Description != "" {
				operation.Request.Body.Description = param.Description
			}

			for name, schema := range param.Schema.Properties {
				property := s.parseSchema(name, &schema)
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

	for code, specResponse := range specOperation.Responses.ResponsesProps.StatusCodeResponses {
		operation.Responses[code] = &capsules.Response{
			StatusCode: &capsules.StatusCode{
				Code:        code,
				Description: specResponse.Description,
			},
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
		}

		for name, header := range specResponse.Headers {
			operation.Responses[code].Headers.Properties[name] = s.parseHeader(&header)
			operation.Responses[code].Headers.Properties[name].Name = name
			operation.Responses[code].Headers.Properties[name].ID = name
		}

		if specResponse.Schema != nil {
			operation.Responses[code].Body = s.parseSchema(strconv.Itoa(code), specResponse.Schema)

			for name, schema := range specResponse.Schema.Properties {
				property := s.parseSchema(name, &schema)
				property.Name = name
				property.ID = name
				operation.Responses[code].Elements.Properties[name] = property
			}

			for _, schemaKey := range specResponse.Schema.Required {
				if operation.Request.Elements.Properties[schemaKey] != nil {
					operation.Request.Elements.Properties[schemaKey].IsRequired = true
				}
			}
		}
	}

	return operation
}
