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

package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReplaceCurlyBracketWithSquare(t *testing.T) {
	target := "/<bucketName>?stats"
	assert.Equal(t, target, ReplaceCurlyBracketWithSquare("/{bucketName}?stats"))
	assert.Equal(t, target, ReplaceCurlyBracketWithSquare("/{bucketName }?stats"))
	assert.Equal(t, target, ReplaceCurlyBracketWithSquare("/{ bucketName}?stats"))

	anotherTarget := "/<<bucketName>>?stats"
	assert.Equal(t, anotherTarget, ReplaceCurlyBracketWithSquare("/{{bucketName}}?stats"))
	assert.Equal(t, anotherTarget, ReplaceCurlyBracketWithSquare("/{{ bucketName }}?stats"))
}
