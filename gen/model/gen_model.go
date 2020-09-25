package model

import (
	"bytes"
	"fmt"
	"myprojects/tools/common"
	"myprojects/tools/gen"
	"os"
	"strings"
)

type CmdGenModel struct {
	DdlFilePath       string
	WithGormTag       bool
	WithSimpleGormTag bool
	WithJsonTag       bool
	WithDefaultTag    bool
}

func (cgm *CmdGenModel) CmdHandle() {
	packageFile, err := common.ParseFilePath()
	if err != nil {
		panic(err)
	}

	rootPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// parse sql
	gormFlags := gen.GormFlags{
		HasGorm:      cgm.WithGormTag,
		IsSimpleGorm: cgm.WithSimpleGormTag,
		HasJson:      cgm.WithJsonTag,
		HasDefault:   cgm.WithDefaultTag,
	}
	gormTable := gen.GormTableList{}
	rs, err := gormTable.Parse(cgm.DdlFilePath, gormFlags)
	if err != nil {
		panic(err)
	}

	for _, tpmData := range rs {
		tpData := TemplateDataModel{
			GenTemplate: gen.GenTemplate{},
			TemplatePackage: gen.TemplatePackage{
				PackageName: packageFile.PackageName,
				PackageList: map[string]string{},
			},
			TemplateModel: tpmData,
		}
		for _, field := range tpData.TemplateModelFields {
			if strings.Contains(field.Type, "time") {
				tpData.PackageList["time"] = "time"
			}
		}

		var codeData *bytes.Buffer
		if codeData, err = tpData.ParseTemplate(templateModelTxt, tpData.ModelName, tpData); err != nil {
			panic(err)
		}

		filename := common.ToLowerSnakeCase(tpData.ModelName)
		filePath := fmt.Sprintf("%s/%s.go", rootPath, filename)
		if err = tpData.FormatCodeToFile(filePath, codeData); err != nil {
			panic(err)
		}
		fmt.Println(tpData.ModelName + " Success")
	}
}
