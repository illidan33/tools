package client

import (
	"errors"
	"fmt"
	"github.com/illidan33/tools/common"
	"github.com/illidan33/tools/gen"
	"os"
	"path/filepath"
	"strings"
)

type CmdGenClient struct {
	DocUrl      string
	ServiceName string
	IsDebug     bool

	Environments common.CmdFilePath
	TemplateGenClient
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
	cmdtp.PackageName = "client_" + common.ToLowerSnakeCase(cmdtp.ServiceName)
	cmdtp.ClientModel.ModelName = common.ToUpperCamelCase(cmdtp.PackageName)

	if cmdtp.IsDebug {
		fmt.Printf("%#v\n", cmdtp.Environments)
	}
	return nil
}

func (cmdtp *CmdGenClient) Parse() error {
	err := cmdtp.ParseSwagger(cmdtp.DocUrl)
	if err != nil {
		return err
	}

	folderPath := filepath.Join(cmdtp.Environments.CmdDir, cmdtp.PackageName)
	if !common.IsExists(folderPath) {
		err := os.MkdirAll(folderPath, os.ModePerm)
		if err != nil {
			return err
		}
	}

	cmdtp.InitTemplateFuncs()
	cmdtp.RegisteTemplateFunc(map[string]interface{}{
		"isModel": func(m gen.TemplateModel) bool {
			if m.ModelName == "" {
				return false
			}
			return true
		},
		"isStruct": func(m gen.TemplateModel) bool {
			if m.Type == "struct" || strings.Contains(m.Type, "[]") {
				return true
			}
			return false
		},
	})

	bf, err := cmdtp.ParseTemplate(templateModelTxt, "templateModelTxt", cmdtp)
	if err != nil {
		return err
	}
	if cmdtp.IsDebug {
		fmt.Println(bf.String())
	}

	err = cmdtp.FormatCodeToFile(filepath.Join(folderPath, "types_generate.go"), bf)
	if err != nil {
		return err
	}

	bf, err = cmdtp.ParseTemplate(templateClientTxt, "templateClientTxt", cmdtp)
	if err != nil {
		return err
	}
	if cmdtp.IsDebug {
		fmt.Println(bf.String())
	}

	err = cmdtp.FormatCodeToFile(filepath.Join(folderPath, "client_generate.go"), bf)
	if err != nil {
		return err
	}

	err = cmdtp.ParseTemplateAndFormatToFile(cmdtp.Environments.CmdDir)
	if err != nil {
		return err
	}
	return nil
}
