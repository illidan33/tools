package method

import (
	"fmt"
	"myprojects/tools/gen"
	"path/filepath"
	"strings"
)

type CmdGenMethodFlags struct {
	CmdGenmodelName     string
	CmdGenmodelFilePath string
}

var tpData TemplateGenMethod

func init() {
	tpData = TemplateGenMethod{
		TemplateGenmodel: gen.TemplateGenmodel{
			TemplateGenStruct: gen.TemplateGenStruct{
				PackageName:   "",
				PackageList:   map[string]string{},
				StructName:    "",
				StructComment: "",
				TemplateFuncs: map[string]interface{}{},
			},
			ModelStructFields: map[string]gen.TemplateGenStructField{},
		},
		CommentIndexs: []TemplateGenMethodCommentIndex{},
		TemplateGenMethodFuncs: []string{},
	}
	registeTemplateFunc(&tpData)
}

func (flags CmdGenMethodFlags) CmdHandle() {
	var err error
	flags.CmdGenmodelFilePath, err = filepath.Abs(flags.CmdGenmodelFilePath)
	if err != nil {
		panic(err)
	}

	if !gen.IsExists(flags.CmdGenmodelFilePath) {
		panic(fmt.Errorf("model file path not exists: %s", flags.CmdGenmodelFilePath))
	}

	// init template data
	if err := tpData.Parse(flags); err != nil {
		panic(err)
	}

	bf, err := tpData.ParseTemplate(templateMethodTxt)
	if err != nil {
		panic(err)
	}

	basePath, fileName := filepath.Split(flags.CmdGenmodelFilePath)
	fns := strings.LastIndex(fileName, ".")
	err = tpData.FormatCodeToFile(fmt.Sprintf("%s/%s_generate.go", basePath, fileName[:fns]), bf)
	if err != nil {
		panic(err)
	}
}
