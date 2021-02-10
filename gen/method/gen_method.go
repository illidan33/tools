package method

import (
	"fmt"
	"tools/common"
	"path/filepath"
)

type CmdGenMethod struct {
	IsDebug   bool
	ModelName string
	ModelFile string

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

	// for test
	if cmdtp.IsDebug {
		fmt.Printf("%#v\n", cmdtp.Environments)
		if cmdtp.Environments.PackageName == "main" {
			cmdtp.Environments.PackageName = "model_test"
			cmdtp.Environments.CmdDir = filepath.Join(common.GetGoPath(), "/github.com/illidan33/tools/example/model")
		}
	}
	if cmdtp.Environments.CmdFileName == "" {
		cmdtp.Environments.CmdFileName = common.ToLowerSnakeCase(cmdtp.ModelName) + ".go"
	}
	if cmdtp.ModelFile == "" {
		cmdtp.ModelFile = filepath.Join(cmdtp.Environments.CmdDir, cmdtp.Environments.CmdFileName)
	}
	cmdtp.Template.PackageName = cmdtp.Environments.PackageName
	cmdtp.Template.ModelName = cmdtp.ModelName

	if cmdtp.ModelName == "" {
		return fmt.Errorf("need model name")
	}
	if cmdtp.ModelFile == "" || common.IsDir(cmdtp.ModelFile) {
		return fmt.Errorf("model path need a dir")
	}

	return nil
}
func (cmdtp *CmdGenMethod) Parse() error {
	//tmpPath, err := common.GetImportPath(cmdtp.Environments.CmdDir)
	//if err != nil {
	//	return err
	//}
	//if err := cmdtp.Template.ParseImportFile(tmpPath); err != nil {
	//	fmt.Println(err) // 记录错误，不打断，退化到由语法树来解析字段
	//}

	filePath := filepath.Join(cmdtp.Environments.CmdDir, cmdtp.Environments.CmdFileName)
	dstTree, err := cmdtp.Template.GetDstTree(filePath)
	if err != nil {
		return err
	}
	if err = cmdtp.Template.ParseDstTree(dstTree); err != nil {
		return err
	}
	if err = cmdtp.Template.ParseIndexToMethod(templateIndexMap, templateIndexUniqMap, templateBaseMap); err != nil {
		return err
	}

	bf, err := cmdtp.Template.ParseTemplate(templateTxt, cmdtp.ModelName, cmdtp.Template)
	if err != nil {
		return err
	}

	dstFilePath := filepath.Join(cmdtp.Environments.CmdDir, common.ToLowerSnakeCase(cmdtp.ModelName)+"_generate.go")
	err = cmdtp.Template.FormatCodeToFile(dstFilePath, bf)
	if err != nil {
		return err
	}
	return nil
}
