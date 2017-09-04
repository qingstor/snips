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

import (
	"go/format"
	"strings"
)

// BaseCapsule provides necessary data to fill the template.
type BaseCapsule struct {
	*CapsulePowder
	Language string // go golang etc
}

// FormatCode formats the generated code.
func (c *BaseCapsule) FormatCode(code string) (string, error) {
	switch strings.ToLower(c.Language) {
	case "go", "golang":
		goodCode, err := format.Source([]byte(code))
		return string(goodCode), err
	}
	return code, nil
}
