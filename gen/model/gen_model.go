package model

import (
	"bytes"
	"fmt"
	"github.com/illidan33/tools/common"
	"github.com/illidan33/tools/gen"
	"os"
	"path/filepath"
	"strings"
)

type CmdGenModel struct {
	DdlFilePath       string
	WithGormTag       bool
	WithSimpleGormTag bool
	WithJsonTag       bool
	WithDefaultTag    bool
	IsDebug           bool
}

func (cgm *CmdGenModel) CmdHandle() {
	packageFile, err := common.GetGenEnvironmentValues(cgm.IsDebug)
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
	if !filepath.IsAbs(cgm.DdlFilePath) {
		cgm.DdlFilePath = filepath.Join(rootPath, cgm.DdlFilePath)
	}
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
		if codeData, err = tpData.ParseTemplate(templateModelTxt, tpData.ModelName, tpData, map[string]interface{}{
			"hasComment": func(field gen.TemplateModelField) bool {
				if field.Comment != "" {
					return true
				}
				return false
			},
		}); err != nil {
			panic(err)
		}

		filename := common.ToLowerSnakeCase(tpData.ModelName)
		filePath := filepath.Join(rootPath, filename+".go")
		if err = tpData.FormatCodeToFile(filePath, codeData); err != nil {
			panic(err)
		}
		fmt.Println(tpData.ModelName + " Success")
	}
}
