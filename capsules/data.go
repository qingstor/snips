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

package capsules

// Data stores the data of an API service to render.
type Data struct {
	Service         *Service
	SubServices     map[string]*SubService
	CustomizedTypes map[string]*Property
}

// Service stores the data of a service.
type Service struct {
	APIVersion  string
	Name        string
	Description string
	Properties  *Property
	Operations  map[string]*Operation
}

// SubService stores the data of an sub service.
type SubService struct {
	ID         string
	Name       string
	Properties *Property
	Operations map[string]*Operation
}

// Operation stores the data of an operation.
type Operation struct {
	ID               string
	Name             string
	Description      string
	DocumentationURL string
	Request          *Request
	Response         *Response
}

// Request stores the data of request section.
type Request struct {
	Method   string
	URI      string
	Params   *Property
	Headers  *Property
	Elements *Property
	Body     *Property
}

// Response stores the data of response section.
type Response struct {
	StatusCodes map[int]*StatusCode
	Headers     *Property
	Elements    *Property
	Body        *Property
}

// StatusCode stores the data of status code.
type StatusCode struct {
	Description string
}

// Property describes info of a property.
type Property struct {
	ID          string
	Name        string
	Description string
	Type        string
	ExtraType   string
	Format      string
	Enum        []string
	Default     string
	IsRequired  bool
	Properties  map[string]*Property
}
