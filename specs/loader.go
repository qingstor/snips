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
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/yunify/snips/capsules"
	"github.com/yunify/snips/utils"
)

// A Service holds the information of API service.
type Service struct {
	Filename         string
	FilePath         string
	LatestAPIVersion *APIVersion
	APIVersions      map[string]*APIVersion
}

// A APIVersion holds the information of an API service's version.
type APIVersion struct {
	Filename string
	FilePath string
	Spec     *Spec
}

// A Spec holds the information of an API spec file.
type Spec struct {
	Filename    string
	FilePath    string
	FileContent string

	Format string

	Data *capsules.Data
}

// LoadServices walks through the specs directory and load API spec information.
func LoadServices(specDirectory, specFormat string, serviceModule string) (*Service, error) {
	if serviceModule != strings.ToLower(serviceModule) {
		serviceModule = utils.CamelCaseToSnakeCase(serviceModule)
	}
	serviceModule = utils.SnakeCaseToSnakeCase(serviceModule, true)

	if _, err := os.Stat(specDirectory + "/" + serviceModule); err != nil {
		return nil, fmt.Errorf("spec of service \"%s\" not found", serviceModule)
	}

	service := &Service{
		Filename: serviceModule,
		FilePath: specDirectory + "/" + serviceModule,
	}

	apiVersions, err := LoadAPIVersions(service, specFormat)

	if err != nil {
		return nil, err
	}
	service.APIVersions = apiVersions

	var latestAPIVersion *APIVersion
	for _, apiVersion := range apiVersions {
		if latestAPIVersion != nil {
			if apiVersion.Filename > latestAPIVersion.Filename {
				latestAPIVersion = apiVersion
			}
		} else {
			latestAPIVersion = apiVersion
		}
	}
	service.LatestAPIVersion = latestAPIVersion
	service.APIVersions["latest"] = latestAPIVersion

	return service, nil
}

// LoadAPIVersions loads all API version files information.
func LoadAPIVersions(service *Service, specFormat string) (map[string]*APIVersion, error) {
	apiVersions := map[string]*APIVersion{}

	files, err := ioutil.ReadDir(service.FilePath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !file.IsDir() {
			continue
		}
		//skip hidden dir.
		if strings.HasPrefix(file.Name(), "."){
			continue
		}

		apiVersion := &APIVersion{
			Filename: file.Name(),
			FilePath: service.FilePath + "/" + file.Name(),
		}

		var format, filename string
		switch specFormat {
		case "Swagger", "Swagger-v2.0", "OpenAPI", "OpenAPI-v2.0":
			format = "swagger"
			filename = "api_v2.0.json"
		default:
			return apiVersions, errors.New("Spec format not supported: " + specFormat)
		}

		specFilePath := apiVersion.FilePath + "/" + format + "/" + filename
		specFileContent, err := ioutil.ReadFile(specFilePath)
		if err != nil {
			return apiVersions, err
		}

		apiVersion.Spec = &Spec{
			Filename:    filename,
			FilePath:    specFilePath,
			FileContent: string(specFileContent),
			Format:      specFormat,
			Data:        &capsules.Data{},
		}

		swagger := Swagger{
			FilePath: specFilePath,
			Data:     apiVersion.Spec.Data,
		}
		err = swagger.Parse("v2.0")
		if err != nil {
			return apiVersions, err
		}

		apiVersions[apiVersion.Filename] = apiVersion
	}

	return apiVersions, nil
}
