package dao_sync

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/illidan33/tools/common"
)

type CmdKipleInterfaceCheck struct {
	InterfaceName string
	ModelName     string
	IsDebug       bool

	Environments common.CmdFilePath
	Template     KipleTemplateDaoSync
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
			cmdtp.Environments.CmdDir = filepath.Join(common.GetGoPath(), "/src/github.com/m2c/kiplelive-agent/service")
			cmdtp.Environments.CmdFileName = "user_service.go"
		}
	}
	cmdtp.Template.InterfaceName = cmdtp.InterfaceName
	cmdtp.Template.ModelName = cmdtp.ModelName

	return nil
}

func (cmdtp *CmdKipleInterfaceCheck) Parse() error {
	excuteFilePath := filepath.Join(cmdtp.Environments.CmdDir, cmdtp.Environments.CmdFileName)
	if !common.IsExists(excuteFilePath) {
		return errors.New("file not exist: " + excuteFilePath)
	}
	dstfl, err := cmdtp.Template.GetDstTree(excuteFilePath)
	if err != nil {
		return err
	}
	dstfl, err = cmdtp.Template.FindInterfaceMethods(dstfl)
	if err != nil {
		return err
	}

	err = cmdtp.Template.ParseToFile(excuteFilePath, dstfl)
	if err != nil {
		return err
	}

	return nil
}
