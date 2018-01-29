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
	"regexp"
	"strings"
	"unicode"
)

// CamelCase converts a string to camel case string.
func CamelCase(original string) string {
	if strings.ToLower(original) == original {
		original = SnakeCaseToCamelCase(original)
	}

	return CamelCaseToCamelCase(original, true)
}

// CamelCaseToCamelCase converts CamelCase to CamelCase.
func CamelCaseToCamelCase(original string, clean ...bool) string {
	converted := original
	for before, after := range capitalizedToCapitalizedWordMap {
		converted = strings.Replace(converted, before, after, -1)
	}

	if len(clean) == 1 && clean[0] {
		converted = strings.Replace(converted, " ", "", -1)
		converted = strings.Replace(converted, "-", "", -1)
		converted = strings.Replace(converted, "_", "", -1)
	}

	return converted
}

// CamelCaseToSnakeCase converts CamelCase to SnakeCase.
func CamelCaseToSnakeCase(original string) string {
	converted := CamelCaseToCamelCase(original, true)

	for _, word := range abbreviateWordMap {
		after := strings.ToLower(word)
		if after != "" {
			after = strings.ToUpper(string(after[0])) + after[1:]
		}

		converted = strings.Replace(converted, word, after, -1)
	}

	upper := regexp.MustCompile(`(^|\W)[A-Z]`).FindString(converted)
	converted = strings.Replace(converted, upper, strings.ToLower(upper), 1)

	oldRunes := []rune(converted)
	newRunes := []rune{}

	for index := 0; index < len(oldRunes); index++ {
		if unicode.IsUpper(oldRunes[index]) {
			newRunes = append(newRunes, '_')
			newRunes = append(newRunes, unicode.ToLower(oldRunes[index]))
		} else {
			newRunes = append(newRunes, oldRunes[index])
		}
	}

	return string(newRunes)
}

// CamelCaseToDashConnected converts CamelCase to DashConnected.
func CamelCaseToDashConnected(original string) string {
	return strings.Replace(CamelCaseToSnakeCase(original), "_", "-", -1)
}

// LowerFirstCharacter makes the first character in a string lowercase.
func LowerFirstCharacter(original string) string {
	if len(original) > 0 {
		original = strings.ToLower(original[:1]) + original[1:]
	}
	return original
}

// UpperFirstCharacter makes the first character in a string uppercase.
func UpperFirstCharacter(original string) string {
	if len(original) > 0 {
		original = strings.ToUpper(original[:1]) + original[1:]
	}
	return original
}

// LowerFirstWord makes the first word in a string lowercase.
func LowerFirstWord(original string) string {
	split := strings.Split(CamelCaseToSnakeCase(original), "_")
	if len(split) > 0 {
		return split[0] + original[len(split[0]):]
	}
	return original
}
