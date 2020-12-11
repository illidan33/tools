package swagger

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ghodss/yaml"
	"path/filepath"
	"tools/common"
	"tools/gen"
)

type CmdKipleSwagger struct {
	Controller string
	Pojo       string
	IsDebug    bool

	Environments common.CmdFilePath
	Template     TemplateKipleSwagger
}

func (cmdtp *CmdKipleSwagger) String() string {
	return cmdtp.Environments.PackageName
}

func (cmdtp *CmdKipleSwagger) Init() error {
	cmdtp.Template.InitTemplateFuncs()

	var err error
	cmdtp.Environments, err = common.GetGenEnvironmentValues()
	if err != nil {
		return err
	}

	// for test
	if cmdtp.IsDebug {
		fmt.Printf("%#v\n", cmdtp.Environments)
		if cmdtp.Environments.PackageName == "main" {
			cmdtp.Environments.CmdDir = filepath.Join(common.GetGoPath(), "/src/github.com/illidan33/tools/example/swag")
			cmdtp.Environments.CmdFileName = "main.go"
			cmdtp.Environments.CmdLine = "7"
		}
	}
	cmdtp.Controller, err = filepath.Abs(filepath.Join(cmdtp.Environments.CmdDir, cmdtp.Controller))
	if err != nil {
		return err
	}
	cmdtp.Pojo, err = filepath.Abs(filepath.Join(cmdtp.Environments.CmdDir, cmdtp.Pojo))
	if err != nil {
		return err
	}

	return nil
}

func (cmdtp *CmdKipleSwagger) Parse() error {
	if !common.IsExists(cmdtp.Environments.CmdDir) {
		return errors.New("empty cmddir")
	}
	cmdtp.Template.ModelList = map[string]gen.TemplateModel{}
	cmdtp.Template.Swagger.Swagger = "2.0"
	cmdtp.Template.ControllerUrls = map[string]string{}
	cmdtp.Template.TemplateSwaggerPaths = []TemplateSwaggerPath{}
	cmdtp.Template.Swagger.Schemes = "{{ marshal .Schemes }}"
	cmdtp.Template.Swagger.Host = "{{.Host}}"
	cmdtp.Template.Swagger.BasePath = "{{.BasePath}}"
	cmdtp.Template.Swagger.Info = SwaggerInfo{
		Contact: struct {
		}{},
	}
	cmdtp.Template.Swagger.Paths = map[string]map[string]SwaggerPath{}
	cmdtp.Template.Swagger.Definitions = map[string]SwaggerDefinition{}

	var err error
	err = cmdtp.Template.ParseSwagTitle(filepath.Join(cmdtp.Environments.CmdDir, cmdtp.Environments.CmdFileName))
	if err != nil {
		return err
	}

	err = cmdtp.Template.ParsePojoDir(cmdtp.Pojo)
	if err != nil {
		return err
	}

	err = cmdtp.Template.ParseControllerDir(cmdtp.Controller)
	if err != nil {
		return err
	}
	cmdtp.Template.SetSwaggerPaths()

	content, err := json.MarshalIndent(cmdtp.Template.Swagger.SwaggerRoot, "", "  ")
	if err != nil {
		return err
	}
	yamlContent, err := yaml.JSONToYAML(content)
	//yamlContent, err := yaml.Marshal(cmdtp.Template.Swagger.SwaggerRoot)
	if err != nil {
		return err
	}

	err = cmdtp.Template.WriteToFile(filepath.Join(cmdtp.Environments.CmdDir, "docs/swagger.json"), content)
	if err != nil {
		return err
	}
	err = cmdtp.Template.WriteToFile(filepath.Join(cmdtp.Environments.CmdDir, "docs/swagger.yaml"), yamlContent)
	if err != nil {
		return err
	}

	docContent, err := json.MarshalIndent(cmdtp.Template.Swagger, "", "  ")
	if err != nil {
		return err
	}
	bt, err := cmdtp.Template.ParseTemplate(templateSwagTxt, "templateSwagTxt", map[string]string{
		"Docs": fmt.Sprintf("`%s`", string(docContent)),
	})
	if err != nil {
		return err
	}
	err = cmdtp.Template.FormatCodeToFile(filepath.Join(cmdtp.Environments.CmdDir, "docs/docs.go"), bt)
	if err != nil {
		return err
	}

	return nil
}
