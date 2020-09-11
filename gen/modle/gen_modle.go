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

var tpData TemplateModleStruct

func init() {
	tpData = TemplateModleStruct{
		TemplateGenStruct: gen.TemplateGenStruct{
			PackageName:   "",
			PackageList:   map[string]string{},
			StructName:     "",
			TemplateFuncs: map[string]interface{}{},
		},
		GormTable: gen.GormTable{},
	}
	registeTemplateFunc(&tpData)
}

func (flag *CmdGenModleFlags) CmdHandle() {
	var err error
	flag.CmdGenModleFilePath, err = filepath.Abs(flag.CmdGenModleFilePath)
	if err != nil {
		panic(err)
	}

	if !gen.IsExists(flag.CmdGenModleFilePath) {
		panic(fmt.Errorf("modle file path not exists: %s", flag.CmdGenModleFilePath))
	}

	exeFilePath, packageName, err := gen.GetExeFilePath()
	if err != nil {
		panic(fmt.Errorf("panic: file path not exists; calltrace:%s", string(debug.Stack())))
	}

	if gen.IsDir(flag.CmdGenModleFilePath) {
		panic(fmt.Errorf("file path is not a file"))

	}
	_, fileName := filepath.Split(flag.CmdGenModleFilePath)
	if err = flag.walkSqlPath(exeFilePath, packageName, flag.CmdGenModleFilePath, fileName); err != nil {
		panic(err)
	}
}

func (flag *CmdGenModleFlags) walkSqlPath(createPath string, pkName string, sqlPath string, filename string) (err error) {
	// parse sql
	gormTable := gen.GormTable{
		Fields: map[string]gen.GormField{},
		Indexs: map[string]gen.GormIndex{},
	}
	err = gormTable.Parse(sqlPath)
	if err != nil {
		return
	}

	// init template data
	tpData.PackageName = pkName
	tpData.StructName = gen.ToUpperCamelCase(gormTable.Name)
	tpData.GormTable = gormTable
	flag.CmdGenModleName = tpData.StructName

	// create new file
	bf, e := tpData.ParseTemplate(templateModleTxt, tpData)
	if e != nil {
		err = e
		return
	}
	filePath := fmt.Sprintf("%s/%s.go", createPath, gen.ToLowerSnakeCase(gormTable.Name))
	err = tpData.FormatCodeToFile(filePath, bf)

	return
}
