package model

import (
	"fmt"
	"myprojects/tools/gen"
	"path/filepath"
	"runtime/debug"
)

type CmdGenmodelFlags struct {
	CmdGenmodelName              string
	CmdGenmodelFilePath          string
	CmdGenmodelWithGormTag       bool
	CmdGenmodelWithSimpleGormTag bool
	CmdGenmodelWithJsonTag       bool
	CmdGenmodelWithDefaultTag    bool
}

var tpData gen.TemplateGenmodel

func init() {
	tpData = gen.TemplateGenmodel{
		TemplateGenStruct: gen.TemplateGenStruct{
			PackageName:   "",
			PackageList:   map[string]string{},
			StructName:    "",
			TemplateFuncs: map[string]interface{}{},
		},
		ModelStructFields: map[string]gen.TemplateGenStructField{},
	}
	registeTemplateFunc(&tpData)
}

func (flag *CmdGenmodelFlags) CmdHandle() {
	var err error
	flag.CmdGenmodelFilePath, err = filepath.Abs(flag.CmdGenmodelFilePath)
	if err != nil {
		panic(err)
	}
	if gen.IsDir(flag.CmdGenmodelFilePath) {
		panic(fmt.Errorf("file path is not a file"))
	}

	exeFilePath, packageName, err := gen.GetExeFilePath()
	if err != nil {
		panic(fmt.Errorf("panic: file path not exists; calltrace:%s", string(debug.Stack())))
	}

	// parse sql
	gormTable := gen.GormTable{
		Fields: map[string]gen.GormField{},
		Indexs: map[string]gen.GormIndex{},
	}
	err = gormTable.Parse(flag.CmdGenmodelFilePath)
	if err != nil {
		panic(err)
	}

	// init template data
	tpData.TransformGormTomodel(gormTable, flag.CmdGenmodelWithGormTag, flag.CmdGenmodelWithSimpleGormTag, flag.CmdGenmodelWithJsonTag, flag.CmdGenmodelWithDefaultTag)
	tpData.PackageName = packageName
	tpData.StructName = gen.ToUpperCamelCase(gormTable.Name)

	// create new file
	bf, err := tpData.ParseTemplate(templatemodelTxt)
	if err != nil {
		panic(err)
	}

	filename := gormTable.Name
	if flag.CmdGenmodelName != "" {
		filename = flag.CmdGenmodelName
		tpData.StructName = flag.CmdGenmodelName
	}
	filePath := fmt.Sprintf("%s/%s.go", exeFilePath, gen.ToLowerSnakeCase(filename))
	err = tpData.FormatCodeToFile(filePath, bf)
	if err != nil {
		panic(err)
	}
}
