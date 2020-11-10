package client

import (
	"errors"
	"fmt"
	"github.com/illidan33/tools/common"
	"path/filepath"
)

type CmdGenClient struct {
	DocUrl      string
	ServiceName string
	IsDebug     bool

	Environments common.CmdFilePath
	Template     TemplateGenClient
}

func (cmdtp *CmdGenClient) String() string {
	return cmdtp.ServiceName
}

func (cmdtp *CmdGenClient) Init() error {
	if cmdtp.DocUrl == "" {
		return errors.New("DocUrl required")
	}
	if cmdtp.ServiceName == "" {
		return errors.New("ServiceName required")
	}

	var err error
	cmdtp.Environments, err = common.GetGenEnvironmentValues()
	if err != nil {
		return err
	}
	cmdtp.Template.PackageName = "client_" + common.ToLowerSnakeCase(cmdtp.ServiceName)
	cmdtp.Template.ClientModel.ModelName = common.ToUpperCamelCase(cmdtp.Template.PackageName)

	if cmdtp.IsDebug {
		fmt.Printf("%#v\n", cmdtp.Environments)
	}

	// init func map
	cmdtp.Template.RegisteFuncMap()

	return nil
}

func (cmdtp *CmdGenClient) Parse() error {
	err := cmdtp.Template.ParseSwagger(cmdtp.DocUrl)
	if err != nil {
		return err
	}

	bfModel, err := cmdtp.Template.ParseTemplate(templateModelTxt, "templateModelTxt", cmdtp.Template)
	if err != nil {
		return err
	}
	if cmdtp.IsDebug {
		fmt.Println(bfModel.String())
	}

	bfClient, err := cmdtp.Template.ParseTemplate(templateClientTxt, "templateClientTxt", cmdtp.Template)
	if err != nil {
		return err
	}
	if cmdtp.IsDebug {
		fmt.Println(bfClient.String())
	}

	folderPath := filepath.Join(cmdtp.Environments.CmdDir, cmdtp.Template.PackageName)
	err = cmdtp.Template.FormatCodeToFile(filepath.Join(folderPath, "types_generate.go"), bfModel)
	if err != nil {
		return err
	}

	err = cmdtp.Template.FormatCodeToFile(filepath.Join(folderPath, "client_generate.go"), bfClient)
	if err != nil {
		return err
	}

	err = cmdtp.Template.ParseTemplateAndFormatToFile(cmdtp.Environments.CmdDir)
	if err != nil {
		return err
	}
	return nil
}
