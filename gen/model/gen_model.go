package model

import (
	"fmt"
	"github.com/illidan33/tools/common"
	"github.com/illidan33/tools/gen"
	"path/filepath"
)

type CmdGenModel struct {
	DdlFilePath       string
	WithGormTag       bool
	WithSimpleGormTag bool
	WithJsonTag       bool
	WithDefaultTag    bool
	IsDebug           bool

	Environments common.CmdFilePath
	NameString   string
	Template     gen.GormTableList
}

func (cmdtp *CmdGenModel) String() string {
	return cmdtp.NameString
}

func (cmdtp *CmdGenModel) Init() error {
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
			cmdtp.Environments.CmdDir = filepath.Join(common.GetGoPath(), "/src/github.com/illidan33/gotest/tools_test/example/model")
		}
	}

	return nil
}
func (cmdtp *CmdGenModel) Parse() error {
	var err error

	// parse sql
	gormFlags := gen.GormFlags{
		HasGorm:      cmdtp.WithGormTag,
		IsSimpleGorm: cmdtp.WithSimpleGormTag,
		HasJson:      cmdtp.WithJsonTag,
		HasDefault:   cmdtp.WithDefaultTag,
	}
	if !filepath.IsAbs(cmdtp.DdlFilePath) {
		cmdtp.DdlFilePath, err = filepath.Abs(cmdtp.DdlFilePath)
		if err != nil {
			return err
		}
	}
	rs, err := cmdtp.Template.Parse(cmdtp.DdlFilePath, gormFlags)
	if err != nil {
		return err
	}

	for _, tpmData := range rs {
		cmdtpTmp := TemplateDataModel{
			GenTemplate: gen.GenTemplate{},
			TemplatePackage: gen.TemplatePackage{
				PackageName: cmdtp.Environments.PackageName,
				PackageList: map[string]string{},
			},
			TemplateModel: tpmData,
		}

		codeData, err := cmdtpTmp.Parse()
		if err != nil {
			return err
		}

		filePath := filepath.Join(cmdtp.Environments.CmdDir, common.ToLowerSnakeCase(cmdtpTmp.ModelName)+".go")
		if err = cmdtpTmp.FormatCodeToFile(filePath, codeData); err != nil {
			return err
		}
		cmdtp.NameString += cmdtpTmp.ModelName + "\r\n"
	}
	return nil
}
