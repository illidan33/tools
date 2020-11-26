package method_sync

import (
	"fmt"
	"github.com/dave/dst"
	"path/filepath"
	"strings"

	"github.com/illidan33/tools/common"
)

type CmdKipleInterfaceCheck struct {
	InterfaceName string
	ModelName     string
	IsDebug       bool

	Environments common.CmdFilePath
	Template     KipleTemplatemethodsync
}

func (cmdtp *CmdKipleInterfaceCheck) String() string {
	return cmdtp.InterfaceName
}

func (cmdtp *CmdKipleInterfaceCheck) Init() error {
	cmdtp.Template.InitTemplateFuncs()

	var err error
	cmdtp.Environments, err = common.GetGenEnvironmentValues()
	if err != nil {
		return err
	}
	if cmdtp.Environments.CmdFileName == "" {
		cmdtp.Environments.CmdFileName = fmt.Sprintf("%s.go", common.ToLowerSnakeCase(cmdtp.InterfaceName))
	}

	// for test
	if cmdtp.IsDebug {
		fmt.Printf("%#v\n", cmdtp.Environments)
		if cmdtp.Environments.PackageName == "main" {
			cmdtp.Environments.CmdDir = filepath.Join(common.GetGoPath(), "/src/github.com/illidan33/gotest/tools_test/example/entity")
			cmdtp.Environments.CmdFileName = "user_profiles_dao.go"
		}
	}
	cmdtp.Template.InterfaceName = cmdtp.InterfaceName
	cmdtp.Template.ModelName = cmdtp.ModelName

	return nil
}

func (cmdtp *CmdKipleInterfaceCheck) Parse() error {
	excuteFilePath := filepath.Join(cmdtp.Environments.CmdDir, cmdtp.Environments.CmdFileName)
	dstfl, err := cmdtp.Template.GetDstTree(excuteFilePath)
	if err != nil {
		return err
	}
	interfaceNode, err := cmdtp.Template.FindSourceInterfaceNode(dstfl)
	if err != nil {
		return err
	}
	interfaceNode.Methods.List = make([]*dst.Field, 0)

	// file interface method
	err = cmdtp.Template.FindInterfaceMethods(dstfl, interfaceNode)
	if err != nil {
		return err
	}
	var genDstfl *dst.File
	end := strings.TrimSuffix(excuteFilePath, ".go")
	genDstfl, err = cmdtp.Template.GetDstTree(end + "_generate.go")
	if err != nil {
		return err
	}
	err = cmdtp.Template.FindInterfaceMethods(genDstfl, interfaceNode)
	if err != nil {
		return err
	}

	// format code
	err = cmdtp.Template.ParseToFile(excuteFilePath, dstfl)
	if err != nil {
		return err
	}

	return nil
}
