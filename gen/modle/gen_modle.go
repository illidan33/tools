package modle

import (
	"fmt"
	"myprojects/tools/gen"
	"path/filepath"
	"runtime/debug"
)

type CmdGenModleFlags struct {
	CmdGenModleName              string
	CmdGenModleFilePath          string
	CmdGenModleWithGormTag       bool
	CmdGenModleWithSimpleGormTag bool
	CmdGenModleWithJsonTag       bool
	CmdGenModleWithDefaultTag    bool
}

var tpData gen.TemplateGenModle

func init() {
	tpData = gen.TemplateGenModle{
		TemplateGenStruct: gen.TemplateGenStruct{
			PackageName:   "",
			PackageList:   map[string]string{},
			StructName:    "",
			TemplateFuncs: map[string]interface{}{},
		},
		ModleStructFields: map[string]gen.TemplateGenStructField{},
	}
	registeTemplateFunc(&tpData)
}

func (flag *CmdGenModleFlags) CmdHandle() {
	var err error
	flag.CmdGenModleFilePath, err = filepath.Abs(flag.CmdGenModleFilePath)
	if err != nil {
		panic(err)
	}
	if gen.IsDir(flag.CmdGenModleFilePath) {
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
	err = gormTable.Parse(flag.CmdGenModleFilePath)
	if err != nil {
		panic(err)
	}

	// init template data
	tpData.TransformGormToModle(gormTable, flag.CmdGenModleWithGormTag, flag.CmdGenModleWithSimpleGormTag, flag.CmdGenModleWithJsonTag, flag.CmdGenModleWithDefaultTag)
	tpData.PackageName = packageName
	tpData.StructName = gen.ToUpperCamelCase(gormTable.Name)

	// create new file
	bf, err := tpData.ParseTemplate(templateModleTxt)
	if err != nil {
		panic(err)
	}

	filename := gormTable.Name
	if flag.CmdGenModleName != "" {
		filename = flag.CmdGenModleName
		tpData.StructName = flag.CmdGenModleName
	}
	filePath := fmt.Sprintf("%s/%s.go", exeFilePath, gen.ToLowerSnakeCase(filename))
	err = tpData.FormatCodeToFile(filePath, bf)
	if err != nil {
		panic(err)
	}
}
