package dao_create

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/dave/dst"
	"github.com/illidan33/tools/common"
	"path/filepath"
)

type CmdKipleDao struct {
	InterfaceName string
	Entity        string
	IsDebug       bool

	Environments common.CmdFilePath
	Template     KipleTemplateDao
}

func (cmdtp *CmdKipleDao) String() string {
	return cmdtp.InterfaceName
}
func (cmdtp *CmdKipleDao) Init() (err error) {
	cmdtp.Environments, err = common.GetGenEnvironmentValues()
	if err != nil {
		return
	}

	if cmdtp.IsDebug {
		fmt.Printf("environment: %#v\n", cmdtp.Environments)
		if cmdtp.Environments.PackageName == "main" {
			cmdtp.Environments.PackageName = "entity"
			cmdtp.Environments.CmdDir = filepath.Join(common.GetGoPath(), "/src/github.com/illidan33/gotest/tools_test/example/entity")
		}
	}
	cmdtp.Template.PackageName = cmdtp.Environments.PackageName
	cmdtp.Template.InterfaceName = cmdtp.InterfaceName
	return
}

func (cmdtp *CmdKipleDao) Parse() (err error) {
	var path string
	path, err = filepath.Abs(cmdtp.Entity)
	if err != nil {
		err = errors.New("can not parse source to abs filepath")
		return
	}

	pathDir := filepath.Dir(path)
	var pathPackage string
	pathPackage, err = common.GetImportPackageName(pathDir)
	if err != nil {
		return
	}
	if cmdtp.Template.PackageName != pathPackage {
		pkgPath, err := common.GetImportPath(pathDir)
		if err != nil {
			return err
		}
		cmdtp.Template.AddPackage("entity", pkgPath)
	}

	var dstTree *dst.File
	dstTree, err = cmdtp.Template.GetDstTree(path)
	if err != nil {
		return
	}
	if err = cmdtp.Template.ParseKipleDstTree(dstTree); err != nil {
		return
	}
	if err = cmdtp.Template.ParseKipleIndexToMethod(); err != nil {
		return
	}

	var bf *bytes.Buffer
	bf, err = cmdtp.Template.ParseTemplate(templateMethodTxt, cmdtp.Template.ModelName, cmdtp.Template)
	if err != nil {
		return
	}

	dstFilePath := filepath.Join(cmdtp.Environments.CmdDir, common.ToLowerSnakeCase(cmdtp.InterfaceName)+".go")
	err = cmdtp.Template.FormatCodeToFile(dstFilePath, bf)
	if err != nil {
		return
	}

	return
}
