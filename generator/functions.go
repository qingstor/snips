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

package generator

import (
	"net/http"
	"sort"
	"strings"
	"text/template"

	"github.com/yunify/snips/capsules"
	"github.com/yunify/snips/utils"
)

var funcMap = template.FuncMap{
	"snakeCase": utils.SnakeCase,
	"camelCase": utils.CamelCase,

	"lower":          lower,
	"lowerFirst":     utils.LowerFirstCharacter,
	"lowerFirstWord": utils.LowerFirstWord,
	"upperFirst":     utils.UpperFirstCharacter,
	"normalized":     normalized,
	"dashConnected":  dashConnected,

	"commaConnected":          commaConnected,
	"commaConnectedWithQuote": commaConnectedWithQuote,

	"replace":     replace,
	"passThrough": passThrough,
	"exist":       exist,

	"firstPropertyIDInCustomizedType": firstPropertyIDInCustomizedType,

	"statusText": statusText,
}

func lower(original string) string {
	return strings.ToLower(original)
}

func normalized(original string) string {
	return utils.CamelCaseToCamelCase(utils.SnakeCaseToSnakeCase(original))
}

func dashConnected(original string) string {
	return utils.SnakeCaseToDashConnected(utils.SnakeCase(original))
}

func commaConnected(stringArray []string) string {
	return strings.Join(stringArray, ", ")
}

func commaConnectedWithQuote(stringArray []string) string {
	quoteStringArray := []string{}
	for _, value := range stringArray {
		quoteStringArray = append(quoteStringArray, `"`+value+`"`)
	}
	return strings.Join(quoteStringArray, ", ")
}

func replace(s, old, new string, n int) string {
	return strings.Replace(s, old, new, n)
}

func passThrough(data ...interface{}) []interface{} {
	return data
}

func firstPropertyIDInCustomizedType(customizedType *capsules.Property) string {
	keys := []string{}
	for key := range customizedType.Properties {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	if len(keys) > 0 {
		return keys[0]
	}

	return ""
}

// statusText translates the integer status code into string text in camelcase.
// For example:
//     200 -> "OK"
//     201 -> "Created"
//     418 -> "Teapot"
func statusText(statusCode int) (statusText string) {
	statusText = http.StatusText(statusCode)

	// Replace special HTTP status description.
	statusText = strings.Replace(statusText, "I'm a teapot", "Teapot", -1)

	// Remove dash and space.
	statusText = strings.Replace(statusText, "-", "", -1)
	statusText = strings.Replace(statusText, " ", "", -1)

	return
}

func exist(data interface{}, index string) bool {
	switch data := data.(type) {
	case []interface{}:
		for _, v := range data {
			if v == index {
				return true
			}
		}
		return false
	case map[string]interface{}:
		_, ok := data[index]
		return ok
	default:
		return false
	}
}
