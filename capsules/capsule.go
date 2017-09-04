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

// Capsule defines the common functions to implement for generator.
type Capsule interface {
	SetData(specData *Data)
	SetMode(templateID string, subServiceID ...string)
}

// CapsulePowder provides necessary data to fill the template.
type CapsulePowder struct {
	Data *Data

	CurrentTemplateID   string
	CurrentSubServiceID string
}

// SetData sets the API spec data.
func (c *CapsulePowder) SetData(specData *Data) {
	c.Data = specData
}

// SetMode sets the template ID and sub service ID.
func (c *CapsulePowder) SetMode(templateID string, subServiceID ...string) {
	c.CurrentTemplateID = templateID
	if len(subServiceID) == 1 {
		c.CurrentSubServiceID = subServiceID[0]
	}
}
