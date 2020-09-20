package model

import (
	"bytes"
	"fmt"
	"myprojects/tools/common"
	"myprojects/tools/gen"
	"os"
)

type CmdGenModel struct {
	ModelName         string
	DdlFilePath       string
	WithGormTag       bool
	WithSimpleGormTag bool
	WithJsonTag       bool
	WithDefaultTag    bool
}

func (cgm *CmdGenModel) CmdHandle() {
	tpData := TemplateDataModel{}
	tpData.InitTemplateFuncs()

	packageFile, err := tpData.ParseFilePath()
	if err != nil {
		panic(err)
	}
	tpData.PackageName = packageFile.PackageName

	rootPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// parse sql
	gormTable := gen.GormTable{}
	err = gormTable.Parse(cgm.DdlFilePath)
	if err != nil {
		panic(err)
	}

	// init template data
	gormFlags := gen.GormFlags{
		HasGorm:      cgm.WithGormTag,
		IsSimpleGorm: cgm.WithSimpleGormTag,
		HasJson:      cgm.WithJsonTag,
		HasDefault:   cgm.WithDefaultTag,
	}
	tpData.TemplatePackage.PackageList, err = gormTable.TransformGormToModel(&tpData.TemplateModel, gormFlags)
	if err != nil {
		panic(err)
	}

	var codeData *bytes.Buffer
	if codeData, err = tpData.ParseTemplate(templateModelTxt, tpData.ModelName, tpData); err != nil {
		panic(err)
	}

	filename := gormTable.Name
	if cgm.ModelName != "" {
		filename = common.ToLowerSnakeCase(cgm.ModelName)
		tpData.TemplateModel.ModelName = common.ToUpperCamelCase(cgm.ModelName)
	}
	filePath := fmt.Sprintf("%s/%s.go", rootPath, common.ToLowerSnakeCase(filename))
	if err = tpData.FormatCodeToFile(filePath, codeData); err != nil {
		panic(err)
	}
}
