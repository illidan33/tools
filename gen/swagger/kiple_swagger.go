package swagger

import (
	"errors"
	"fmt"
	"path/filepath"
	"tools/common"
)

type CmdKipleSwagger struct {
	Controller string
	Pojo       string
	IsInit     uint8
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
			cmdtp.Environments.CmdDir = filepath.Join(common.GetGoPath(), "/github.com/illidan33/tools/example/swag")
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
	//cmdtp.Template.ModelList = map[string]template.TemplatePkgModel{}
	cmdtp.Template.Swagger.Swagger = "2.0"
	cmdtp.Template.TemplateIris = TemplateIris{
		Parties:      map[string]*TemplateIrisParty{},
		Applications: map[string]*TemplateIrisApplication{},
		Controllers:  map[string]TemplateIrisController{},
		Funcs:        map[string]TemplateIrisFunc{},
	}
	cmdtp.Template.Swagger.Schemes = "{{ marshal .Schemes }}"
	cmdtp.Template.Swagger.Host = "{{.Host}}"
	cmdtp.Template.Swagger.BasePath = "{{.BasePath}}"
	cmdtp.Template.Swagger.Info = SwaggerInfo{
		Contact: map[string]string{},
	}
	cmdtp.Template.Swagger.Paths = map[string]map[string]SwaggerPath{}
	cmdtp.Template.Swagger.Definitions = map[string]SwaggerDefinition{}
	cmdtp.Template.IsInit = cmdtp.IsInit
	cmdtp.Template.RouterName = "RegisterGlobalModel"
	cmdtp.Template.ControllerBefore = "BeforeActivation"

	return nil
}

func (cmdtp *CmdKipleSwagger) Parse() error {
	if !common.IsExists(cmdtp.Environments.CmdDir) {
		return errors.New("empty cmddir")
	}

	var err error
	err = cmdtp.Template.ParseSwagTitle(filepath.Join(cmdtp.Environments.CmdDir, cmdtp.Environments.CmdFileName))
	if err != nil {
		return err
	}
	err = cmdtp.Template.ParsePojoDir(cmdtp.Pojo)
	if err != nil {
		return err
	}

	err = cmdtp.Template.ParseControllerRouter(cmdtp.Controller)
	if err != nil {
		return err
	}
	err = cmdtp.Template.ParseControllers(cmdtp.Controller)
	if err != nil {
		return err
	}
	err = cmdtp.Template.SetSwaggerPaths()
	if err != nil {
		return err
	}

	if cmdtp.Template.IsInit == 0 {
		err = cmdtp.Template.FormatToFiles(cmdtp.Environments.CmdDir)
		if err != nil {
			return err
		}
	} else {
		err = cmdtp.Template.OverWriteControllerDir(cmdtp.Controller)
		if err != nil {
			return err
		}
	}

	return nil
}
