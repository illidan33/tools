package method

import (
	"fmt"
	"github.com/illidan33/tools/common"
	"path/filepath"
)

type CmdGenMethod struct {
	IsDebug   bool
	ModelName string

	Environments common.CmdFilePath
	Template     TemplateDataMethod
}

func (cmdtp *CmdGenMethod) String() string {
	return cmdtp.ModelName
}
func (cmdtp *CmdGenMethod) Init() error {
	var err error
	cmdtp.Environments, err = common.GetGenEnvironmentValues()
	if err != nil {
		return err
	}
	if cmdtp.Environments.CmdFileName == "" {
		cmdtp.Environments.CmdFileName = common.ToLowerSnakeCase(cmdtp.ModelName) + ".go"
	}

	// for test
	if cmdtp.IsDebug {
		fmt.Printf("%#v\n", cmdtp.Environments)
		if cmdtp.Environments.PackageName == "main" {
			cmdtp.Environments.PackageName = "model_test"
			cmdtp.Environments.CmdDir = filepath.Join(common.GetGoPath(), "/src/github.com/illidan33/gotest/tools_test/example/model")
		}
	}
	cmdtp.Template.PackageName = cmdtp.Environments.PackageName
	cmdtp.Template.ModelName = cmdtp.ModelName

	return nil
}
func (cmdtp *CmdGenMethod) Parse() error {
	tmpPath, err := common.GetImportPath(cmdtp.Environments.CmdDir)
	if err != nil {
		return err
	}
	if err := cmdtp.Template.ParseImportFile(tmpPath); err != nil {
		fmt.Println(err) // 记录错误，不打断，退化到由语法树来解析字段
	}

	filePath := filepath.Join(cmdtp.Environments.CmdDir, cmdtp.Environments.CmdFileName)
	dstTree, err := cmdtp.Template.GetDstTree(filePath)
	if err != nil {
		return err
	}
	if err = cmdtp.Template.ParseDstTree(dstTree); err != nil {
		return err
	}
	if err = cmdtp.Template.ParseIndexToMethod(); err != nil {
		return err
	}

	bf, err := cmdtp.Template.ParseTemplate(templateMethodTxt, cmdtp.ModelName, cmdtp.Template)
	if err != nil {
		return err
	}
	if cmdtp.IsDebug {
		fmt.Printf(bf.String())
	}

	dstFilePath := filepath.Join(cmdtp.Environments.CmdDir, common.ToLowerSnakeCase(cmdtp.ModelName)+"_generate.go")
	err = cmdtp.Template.FormatCodeToFile(dstFilePath, bf)
	if err != nil {
		return err
	}
	return nil
}
