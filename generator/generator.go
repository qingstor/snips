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
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/yunify/snips/capsules"
	"github.com/yunify/snips/templates"
)

// Generator will load a API spec file and covert it to corresponding programing code.
type Generator struct {
	capsule  capsules.Capsule
	template *templates.Template

	code string
}

// New creates a generator.
func New() *Generator {
	return &Generator{}
}

// Set set properties of a generator.
func (g *Generator) Set(capsule capsules.Capsule, template *templates.Template) {
	g.capsule = capsule
	g.template = template
	g.code = ""
}

// Clear clear the generator's properties.
func (g *Generator) Clear() {
	g.capsule = nil
	g.template = nil
	g.code = ""
}

// Run will render and write code.
func (g *Generator) Run() error {
	err := g.Render()
	if err != nil {
		return err
	}

	err = g.Write()
	if err != nil {
		return err
	}

	return nil
}

// Render coverts API spec data content to programing code.
func (g *Generator) Render() error {
	if g.template.IsNeedGenerate {
		switch g.template.Format {
		case "Go":
			buffer := bytes.Buffer{}
			target := template.Must(template.New("Template").Funcs(funcMap).Parse(g.template.FileContent))

			var err error
			err = target.Execute(&buffer, g.capsule)
			if err != nil {
				return err
			}

			g.code = buffer.String()
		default:
			return fmt.Errorf("Template format not supported: \"%s\"", g.template.Format)
		}
	} else {
		g.code = string(g.template.FileContent)
	}

	return nil
}

// Write writes the generated code into the target file.
func (g *Generator) Write() error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	fmt.Println("Generating to: ." + strings.Replace(g.template.OutputFilePath, pwd, "", -1))

	err = os.MkdirAll(filepath.Dir(g.template.OutputFilePath), 0755)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(g.template.OutputFilePath, []byte(g.code), 0644)
	if err != nil {
		return err
	}

	return nil
}
