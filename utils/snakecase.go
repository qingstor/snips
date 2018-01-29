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

import "strings"

// SnakeCase converts a string to snake case string.
func SnakeCase(original string) string {
	if strings.ToLower(original) != original {
		original = CamelCaseToSnakeCase(original)
	}
	return SnakeCaseToSnakeCase(original, true)
}

// SnakeCaseToSnakeCase converts SnakeCase to SnakeCase.
func SnakeCaseToSnakeCase(original string, clean ...bool) string {
	converted := original

	if len(clean) == 1 && clean[0] {
		converted = strings.Replace(converted, " ", "_", -1)
		converted = strings.Replace(converted, "-", "_", -1)
		converted = strings.ToLower(converted)
	}

	convertedParts := strings.Split(converted, "_")

	for index, part := range convertedParts {
		if part != "" {
			if word := lowercaseToLowercaseWordMap[part]; word != "" {
				part = word
			}
		}
		convertedParts[index] = part
	}

	return strings.Join(convertedParts, "_")
}

// SnakeCaseToCamelCase converts SnakeCase to CamelCase.
func SnakeCaseToCamelCase(original string) string {
	converted := SnakeCaseToSnakeCase(original, true)
	convertedParts := strings.Split(converted, "_")

	for index, part := range convertedParts {
		if part != "" {
			if word := lowercaseToCapitalizedWordMap[part]; word != "" {
				part = word
			} else {
				part = strings.ToUpper(string(part[0])) + part[1:]
			}
		}
		convertedParts[index] = part
	}

	return strings.Join(convertedParts, "")
}

// SnakeCaseToDashConnected converts SnakeCase to DashConnected.
func SnakeCaseToDashConnected(original string) string {
	return strings.Replace(SnakeCaseToSnakeCase(original, true), "_", "-", -1)
}
