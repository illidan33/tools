package method

import (
	"fmt"
	"myprojects/tools/gen"
	"path/filepath"
	"strings"
)

type CmdGenMethodFlags struct {
	CmdGenModleName     string
	CmdGenModleFilePath string
}

var tpData TemplateMethod

func init() {
	tpData = TemplateMethod{
		TemplateGenStruct: gen.TemplateGenStruct{
			PackageName:   "",
			PackageList:   map[string]string{},
			StructName:     "",
			TemplateFuncs: map[string]interface{}{},
		},
		GormTable: gen.GormTable{
			Name:          "",
			Engine:        "",
			Charset:       "",
			Collate:       "",
			Comment:       "",
			AutoIncrement: 0,
			Fields:        map[string]gen.GormField{},
			Indexs:        map[string]gen.GormIndex{},
		},
	}
	registeTemplateFunc(&tpData)
}

func (flags CmdGenMethodFlags) CmdHandle() {
	var err error
	flags.CmdGenModleFilePath, err = filepath.Abs(flags.CmdGenModleFilePath)
	if err != nil {
		panic(err)
	}

	if !gen.IsExists(flags.CmdGenModleFilePath) {
		panic(fmt.Errorf("modle file path not exists: %s", flags.CmdGenModleFilePath))
	}

	// init template data
	if err := tpData.Parse(flags); err != nil {
		panic(err)
	}

	bf, err := tpData.ParseTemplate(templateMethodTxt, tpData)
	if err != nil {
		panic(err)
	}

	basePath, fileName := filepath.Split(flags.CmdGenModleFilePath)
	fns := strings.LastIndex(fileName, ".")
	err = tpData.FormatCodeToFile(fmt.Sprintf("%s/%s_generate.go", basePath, fileName[:fns]), bf)
	if err != nil {
		panic(err)
	}
}
