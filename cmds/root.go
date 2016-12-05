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
	"github.com/yunify/snips/generator"
	"github.com/yunify/snips/metadata"
	"github.com/yunify/snips/specs"
	"github.com/yunify/snips/templates"
	"github.com/yunify/snips/utils"
)

var (
	flagVersion           bool
	codeServiceModule     string
	codeServiceAPIVersion string
	codeSpecDirectory     string
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
		&codeServiceModule, "service", "m", "",
		"Choose the service to use.",
	)
	RootCMD.Flags().StringVarP(
		&codeServiceAPIVersion, "service-api-version", "n", "latest",
		"Choose the service API version to use.",
	)
	RootCMD.Flags().StringVarP(
		&codeSpecDirectory, "spec", "s", "",
		"Specify spec files directory.",
	)
	RootCMD.Flags().StringVarP(
		&codeSpecFormat, "spec-format", "", "Swagger-v2.0",
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
	Short: "A code generator for QingCloud & QingStor SDKs.",
	Long: `A code generator for QingCloud & QingStor SDKs.
It is used to generate code from our public APIs currently.

For example:
  $ snips -m QingStor -n latest \
          -s ./specs -t ./templates/qingstor/go \
          -o ./publish/qingstor-sdk-go/service
  $ ...
  $ snips --service=QingStor \
          --service-api-version=latest \
          --spec=./specs \
          --template=./templates/qingstor/ruby \
          --output=./publish/qingstor-sdk-ruby/lib/qingstor/sdk/service
  $ ...

Copyright (C) 2016 Yunify, Inc.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if flagVersion {
			return nil
		}

		if codeSpecDirectory == "" {
			return errors.New("please specify spec files directory")
		}
		if _, err := os.Stat(codeSpecDirectory); err != nil {
			return errors.New("please make sure the specs directory exists")
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
			fmt.Println("snips version " + metadata.Version)
			return
		}

		loadedTemplates, _, err := templates.LoadTemplates(codeTemplateDirectory)
		utils.CheckErrorForExit(err)
		fmt.Println("Loaded templates from " + codeTemplateDirectory)
		fmt.Println(len(loadedTemplates), "template(s) detected.")

		service, err := specs.LoadServices(codeSpecDirectory, codeSpecFormat, codeServiceModule)
		utils.CheckErrorForExit(err)
		apiVersion := service.APIVersions[codeServiceAPIVersion]
		if apiVersion == nil {
			utils.CheckErrorForExit(fmt.Errorf(
				"API version \"%s\" of service \"%s\" not found",
				codeServiceAPIVersion,
				utils.SnakeCaseToCamelCase(service.Filename),
			))
		}

		fmt.Printf(
			"Loaded service %s (%s) from %s\n\n",
			utils.SnakeCaseToCamelCase(service.Filename),
			service.LatestAPIVersion.Filename,
			codeSpecDirectory)

		if spec := apiVersion.Spec; spec != nil {
			codeCapsule := &capsules.BaseCapsule{CapsulePowder: &capsules.CapsulePowder{}}
			codeGenerator := generator.New()

			codeCapsule.SetData(spec.Data)
			codeCapsule.SetVersioning(codeServiceAPIVersion != "latest")

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
