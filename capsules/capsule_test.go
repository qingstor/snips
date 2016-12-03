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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCapsulePowder_SetData(t *testing.T) {
	powder := &CapsulePowder{}
	powder.SetMode("sub_service", "bucket")
	assert.Equal(t, powder.CurrentSubServiceID, "bucket")
}
