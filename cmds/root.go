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

package cmds

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yunify/snips/capsules"
	"github.com/yunify/snips/constants"
	"github.com/yunify/snips/generator"
	"github.com/yunify/snips/specs"
	"github.com/yunify/snips/templates"
	"github.com/yunify/snips/utils"
)

var (
	flagVersion           bool
	codeSpecFile          string
	codeSpecFormat        string
	codeTemplateDirectory string
	codeOutputDirectory   string
)

func init() {
	RootCMD.Flags().BoolVarP(
		&flagVersion, "version", "v", false,
		"Show version.",
	)
	RootCMD.Flags().StringVarP(
		&codeSpecFile, "file", "f", "",
		"Specify the spec file.",
	)
	RootCMD.Flags().StringVarP(
		&codeSpecFormat, "format", "", "Swagger-v2.0",
		"Specify the format of spec file.",
	)
	RootCMD.Flags().StringVarP(
		&codeTemplateDirectory, "template", "t", "",
		"Specify template files directory.",
	)
	RootCMD.Flags().StringVarP(
		&codeOutputDirectory, "output", "o", "",
		"Specify the output directory.",
	)
}

// RootCMD represents the base command when called without any subcommands.
var RootCMD = &cobra.Command{
	Use:   "snips",
	Short: "A code generator for RESTful APIs.",
	Long: `A code generator for RESTful APIs.

For example:
  $ snips -f ./specs/qingstor/api.json
          -t ./templates/qingstor/go \
          -o ./publish/qingstor-sdk-go/service
  $ ...
  $ snips --file=./specs/qingstor/api.json \
          --template=./templates/qingstor/ruby \
          --output=./publish/qingstor-sdk-ruby/lib/qingstor/sdk/service
  $ ...

Copyright (C) 2016-2017 Yunify, Inc.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if flagVersion {
			return nil
		}

		if codeSpecFile == "" {
			return errors.New("please specify the spec file")
		}
		if _, err := os.Stat(codeSpecFile); err != nil {
			return errors.New("please make sure the spec file exists")
		}

		if codeTemplateDirectory == "" {
			return errors.New("please specify templates directory")
		}
		if _, err := os.Stat(codeTemplateDirectory); err != nil {
			return errors.New("please make sure the templates directory exists")
		}

		if codeOutputDirectory == "" {
			return errors.New("please specify output files directory")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if flagVersion {
			fmt.Println("snips version " + constants.Version)
			return
		}

		loadedTemplates, manifest, err := templates.LoadTemplates(codeTemplateDirectory)
		utils.CheckErrorForExit(err)
		fmt.Println("Loaded templates from " + codeTemplateDirectory)
		fmt.Println(len(loadedTemplates), "template(s) detected.")

		spec, err := specs.LoadSpec(codeSpecFile, codeSpecFormat)
		utils.CheckErrorForExit(err)
		fmt.Printf("Loaded specification file %s (%s)\n\n", codeSpecFile, codeSpecFormat)

		if spec != nil {
			codeCapsule := &capsules.BaseCapsule{CapsulePowder: &capsules.CapsulePowder{}}
			codeGenerator := generator.New()

			if manifest.MetaData != nil {
				spec.Data.MetaData = manifest.MetaData
			}
			codeCapsule.SetData(spec.Data)

			sharedTemplateContent := ""
			if template := loadedTemplates["shared"]; template != nil {
				sharedTemplateContent = template.FileContent
			}

			if template := loadedTemplates["service"]; template != nil {
				template.FileContent += sharedTemplateContent
				template.UpdateOutputFilename(spec.Data.Service.Name, template.OutputFileNaming.Style)
				template.UpdateOutputFilePath(codeOutputDirectory)
				codeCapsule.SetMode("", "")
				codeGenerator.Set(codeCapsule, template)
				err = codeGenerator.Run()
				utils.CheckErrorForExit(err)
			}

			if template := loadedTemplates["sub_service"]; template != nil {
				template.FileContent += sharedTemplateContent
				for _, subService := range spec.Data.SubServices {
					template.UpdateOutputFilename(subService.Name, template.OutputFileNaming.Style)
					template.UpdateOutputFilePath(codeOutputDirectory)
					codeCapsule.SetMode(template.ID, subService.ID)
					codeGenerator.Set(codeCapsule, template)
					err = codeGenerator.Run()
					utils.CheckErrorForExit(err)
				}
			}

			if template := loadedTemplates["types"]; template != nil {
				template.FileContent += sharedTemplateContent
				template.UpdateOutputFilename("types", template.OutputFileNaming.Style)
				template.UpdateOutputFilePath(codeOutputDirectory)
				codeCapsule.SetMode("", "")
				codeGenerator.Set(codeCapsule, template)
				err = codeGenerator.Run()
				utils.CheckErrorForExit(err)
			}
		}

		fmt.Println("\nEverything looks fine.")
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by entrance of the application.
// It only needs to happen once to the rootCMD.
func Execute() {
	if err := RootCMD.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
