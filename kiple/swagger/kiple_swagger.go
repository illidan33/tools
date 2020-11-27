package swagger

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/illidan33/tools/common"
	"github.com/illidan33/tools/gen"
	"path/filepath"
)

type CmdKipleSwagger struct {
	ServideDir string
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
	cmdtp.ServideDir, err = filepath.Abs(cmdtp.ServideDir)
	if err != nil {
		return err
	}

	// for test
	if cmdtp.IsDebug {
		fmt.Printf("%#v\n", cmdtp.Environments)
		if cmdtp.Environments.PackageName == "main" {
			cmdtp.Environments.CmdDir = filepath.Join(common.GetGoPath(), "/src/github.com/illidan33/tools/example/service")
			cmdtp.Environments.CmdFileName = "user_service.go"
		}
	}

	return nil
}

func (cmdtp *CmdKipleSwagger) Parse() error {
	if !common.IsExists(cmdtp.Environments.CmdDir) {
		return errors.New("empty cmddir")
	}
	cmdtp.Template.ModelList = map[string]gen.TemplateModel{}
	cmdtp.Template.ImportList = map[string]string{}
	cmdtp.Template.LoadImportList = map[string]bool{}
	cmdtp.Template.Swagger.Paths = map[string]map[string]SwaggerPath{}
	cmdtp.Template.Swagger.Definitions = map[string]SwaggerDefinition{}

	err := cmdtp.Template.ParseServiceDir(cmdtp.ServideDir)
	if err != nil {
		return err
	}

	content, err := json.Marshal(cmdtp.Template.Swagger)
	if err != nil {
		return err
	}
	fmt.Println(content)

	return nil
}
