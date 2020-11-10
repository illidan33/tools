package model

import (
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
	//gen.GenTemplate
	//gen.TemplatePackage
	//gen.TemplateModel
}

func (tpData *CmdGenModel) String() string {
	return tpData.NameString
}

func (tpData *CmdGenModel) Init() error {
	//tpData.InitTemplateFuncs()

	var err error
	tpData.Environments, err = common.GetGenEnvironmentValues()
	if err != nil {
		return err
	}

	// for test
	if tpData.IsDebug {
		tpData.Environments.PackageName = "model_test"
		tpData.Environments.CmdDir = filepath.Join(common.GetGoPath(), "/src/github.com/illidan33/tools/example/model")
	}

	return nil
}
func (tpData *CmdGenModel) Parse() error {
	var err error

	// parse sql
	gormFlags := gen.GormFlags{
		HasGorm:      tpData.WithGormTag,
		IsSimpleGorm: tpData.WithSimpleGormTag,
		HasJson:      tpData.WithJsonTag,
		HasDefault:   tpData.WithDefaultTag,
	}
	gormTable := gen.GormTableList{}
	if !filepath.IsAbs(tpData.DdlFilePath) {
		tpData.DdlFilePath, err = filepath.Abs(tpData.DdlFilePath)
		if err != nil {
			return err
		}
	}
	rs, err := gormTable.Parse(tpData.DdlFilePath, gormFlags)
	if err != nil {
		return err
	}

	for _, tpmData := range rs {
		tpDataTmp := TemplateDataModel{
			GenTemplate: gen.GenTemplate{},
			TemplatePackage: gen.TemplatePackage{
				PackageName: tpData.Environments.PackageName,
				PackageList: map[string]string{},
			},
			TemplateModel: tpmData,
		}

		codeData, err := tpDataTmp.Parse()
		if err != nil {
			return err
		}

		filePath := filepath.Join(tpData.Environments.CmdDir, common.ToLowerSnakeCase(tpDataTmp.ModelName)+".go")
		if err = tpDataTmp.FormatCodeToFile(filePath, codeData); err != nil {
			return err
		}
		tpData.NameString += tpDataTmp.ModelName + "\r\n"
	}
	return nil
}
