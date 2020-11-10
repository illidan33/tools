package method

import (
	"fmt"
	"github.com/illidan33/tools/common"
	"github.com/illidan33/tools/gen"
	"path/filepath"
)

type CmdGenMethod struct {
	IsDebug bool
	// ModelName string // init in cmd flags

	gen.GenTemplate
	gen.TemplatePackage
	gen.TemplateModel
	Environments             common.CmdFilePath
	TemplateDataMethodFuncs  []string
	TemplateDataMethodIndexs []TemplateDataMethodIndex
}

func (cmdtp *CmdGenMethod) String() string {
	return cmdtp.ModelName
}
func (cmdtp *CmdGenMethod) Init() error {
	cmdtp.InitTemplateFuncs()

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
			cmdtp.Environments.CmdDir = filepath.Join(common.GetGoPath(), "/src/github.com/illidan33/tools/example/model")
		}
	}
	cmdtp.PackageName = cmdtp.Environments.PackageName

	return nil
}
func (cmdtp *CmdGenMethod) Parse() error {
	tmpPath, err := common.GetImportPath(cmdtp.Environments.CmdDir)
	if err != nil {
		return err
	}
	if err := cmdtp.ImportFile(tmpPath); err != nil {
		fmt.Println(err) // 记录错误，不打断，退化到由语法树来解析字段
	}

	filePath := filepath.Join(cmdtp.Environments.CmdDir, cmdtp.Environments.CmdFileName)
	dstTree, err := cmdtp.GetDstTree(filePath)
	if err != nil {
		return err
	}
	if err = cmdtp.ParseDstTree(dstTree); err != nil {
		return err
	}
	if err = cmdtp.ParseIndexToMethod(); err != nil {
		return err
	}

	bf, err := cmdtp.ParseTemplate(templateMethodTxt, cmdtp.ModelName, cmdtp)
	if err != nil {
		return err
	}
	if cmdtp.IsDebug {
		fmt.Printf(bf.String())
	}

	dstFilePath := filepath.Join(cmdtp.Environments.CmdDir, common.ToLowerSnakeCase(cmdtp.ModelName)+"_generate.go")
	err = cmdtp.FormatCodeToFile(dstFilePath, bf)
	if err != nil {
		return err
	}
	return nil
}
