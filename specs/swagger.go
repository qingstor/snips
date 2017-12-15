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
	"strings"

	"github.com/imdario/mergo"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"

	"github.com/yunify/snips/capsules"
)

// Swagger holds the data that to parse swagger spec.
type Swagger struct {
	FilePath string
	Data     *capsules.Data
}

// Parse parses swagger spec to data.
func (s *Swagger) Parse(version string) error {
	switch version {
	case "v2.0":
		document, err := loads.Spec(s.FilePath)
		if err != nil {
			return err
		}

		document, err = document.Expanded()
		if err != nil {
			return err
		}

		allProperties := map[string]*capsules.Property{}
		allOperations := map[string]map[string]*capsules.Operation{}
		s.parseOperations(document.Spec(), allProperties, allOperations)

		s.loadData(document.Spec())
		s.loadService(document.Spec(), allProperties, allOperations)
		s.loadSubService(document.Spec(), allProperties, allOperations)
		s.loadCustomizedTypes(document.Spec())
	default:
		return errors.New("Swagger version not supported: " + version)
	}

	return nil
}

func (s *Swagger) parseOperations(
	swagger *spec.Swagger,
	allProperties map[string]*capsules.Property,
	allOperations map[string]map[string]*capsules.Operation) {

	parseOperation := func(uri string, method string,
		specOperation *spec.Operation, property *capsules.Property) {
		if specOperation.ID == "PostObject" {
			return
		}

		sections := []string{}

		if len(specOperation.Tags) > 0 {
			for _, subServiceName := range specOperation.Tags {
				sections = append(sections, subServiceName+"SubService")
			}
		} else {
			serviceName := swagger.Info.Title
			sections = append(sections, serviceName+"Service")
		}

		for _, sectionName := range sections {
			if allProperties[sectionName] == nil {
				allProperties[sectionName] = property
			} else {
				mergo.Merge(allProperties[sectionName], property)
			}

			operation := s.parseOperation(uri, method, property, specOperation, swagger)
			if allOperations[sectionName] == nil {
				allOperations[sectionName] = map[string]*capsules.Operation{}
			}
			allOperations[sectionName][specOperation.ID] = operation
		}
	}

	for requestURI, pathItem := range swagger.Paths.Paths {
		property := &capsules.Property{
			Properties: map[string]*capsules.Property{},
		}

		for _, param := range pathItem.Parameters {
			paramProperty := s.parseParameter(&param, &swagger.Parameters)
			property.Properties[paramProperty.Name] = paramProperty
		}

		if pathItem.Get != nil {
			parseOperation(requestURI, "GET", pathItem.Get, property)
		}
		if pathItem.Put != nil {
			parseOperation(requestURI, "PUT", pathItem.Put, property)
		}
		if pathItem.Post != nil {
			parseOperation(requestURI, "POST", pathItem.Post, property)
		}
		if pathItem.Delete != nil {
			parseOperation(requestURI, "DELETE", pathItem.Delete, property)
		}
		if pathItem.Options != nil {
			parseOperation(requestURI, "OPTIONS", pathItem.Options, property)
		}
		if pathItem.Head != nil {
			parseOperation(requestURI, "HEAD", pathItem.Head, property)
		}
		if pathItem.Patch != nil {
			parseOperation(requestURI, "PATCH", pathItem.Patch, property)
		}
	}
}

func (s *Swagger) loadData(swagger *spec.Swagger) {
	if s.Data == nil {
		s.Data = &capsules.Data{}
	}
	s.Data.Service = nil
	s.Data.SubServices = map[string]*capsules.SubService{}
	s.Data.CustomizedTypes = map[string]*capsules.Property{}
}

func (s *Swagger) loadService(
	swagger *spec.Swagger,
	allProperties map[string]*capsules.Property,
	allOperations map[string]map[string]*capsules.Operation) {
	serviceName := swagger.Info.Title

	property := &capsules.Property{}
	if allProperties[swagger.Info.Title] != nil {
		property = allProperties[serviceName]
	}

	s.Data.Service = &capsules.Service{
		APIVersion:  swagger.Info.Version,
		Name:        serviceName,
		BasePath:    swagger.BasePath,
		Description: swagger.Info.Description,
		Properties:  property,
		Operations:  allOperations[serviceName+"Service"],
	}

	// Be compatible with QingCloud IaaS Services
	if strings.Contains(s.Data.Service.Name, "QingCloud") {
		for _, o := range s.Data.Service.Operations {
			o.Request.Query, o.Request.Elements = o.Request.Elements, o.Request.Query
		}
	}
}

func (s *Swagger) loadSubService(
	swagger *spec.Swagger,
	allProperties map[string]*capsules.Property,
	allOperations map[string]map[string]*capsules.Operation) {

	for subService, operations := range allOperations {
		if strings.Contains(subService, "SubService") {
			subServiceName := strings.Replace(subService, "SubService", "", -1)
			s.Data.SubServices[subServiceName] = &capsules.SubService{
				ID:         subServiceName,
				Name:       subServiceName,
				Properties: allProperties[subService],
				Operations: operations,
			}

			// Be compatible with QingCloud IaaS Services
			if strings.Contains(s.Data.Service.Name, "QingCloud") {
				for _, o := range s.Data.SubServices[subServiceName].Operations {
					o.Request.Query, o.Request.Elements = o.Request.Elements, o.Request.Query
				}
			}
		}
	}
}

func (s *Swagger) loadCustomizedTypes(swagger *spec.Swagger) {
	for name, definition := range swagger.Definitions {
		s.Data.CustomizedTypes[name] = s.parseSchema(name, &definition)
		s.Data.CustomizedTypes[name].ID = name
		s.Data.CustomizedTypes[name].Name = name

		for _, schemaKey := range definition.SchemaProps.Required {
			s.Data.CustomizedTypes[name].Properties[schemaKey].IsRequired = true
		}
	}
}
