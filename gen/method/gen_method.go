package method

import (
	"errors"
	"fmt"
	"myprojects/tools/common"
	"os"
)

type CmdGenMethod struct {
	ModelName string
	IsDebug   bool
}

func (cgm *CmdGenMethod) CmdHandle() {
	tpData := TemplateDataMethod{}
	tpData.InitTemplateFuncs()

	cmdFile, err := common.ParseFilePath(cgm.IsDebug)
	if err != nil {
		panic(err)
	}
	exeFilePath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	tpData.PackageName = cmdFile.PackageName

	// for test
	if cgm.IsDebug {
		exeFilePath = "./example/model"
		cmdFile.CmdFileName = common.ToLowerSnakeCase(cgm.ModelName) + ".go"
		tpData.PackageName = "example"
	}

	filePath := fmt.Sprintf("%s/%s", exeFilePath, cmdFile.CmdFileName)
	if !common.IsExists(filePath) {
		panic(errors.New("File not found"))
	}

	err = tpData.Parse(filePath)
	if err != nil {
		panic(err)
	}

	if tpData.ModelName != cgm.ModelName {
		panic(errors.New("Struct not found: " + cgm.ModelName))
	}

	bf, err := tpData.ParseTemplate(templateMethodTxt, tpData.ModelName, tpData)
	if err != nil {
		panic(err)
	}

	dstFilePath := fmt.Sprintf("%s/%s_generate.go", exeFilePath, common.ToLowerSnakeCase(tpData.ModelName))
	err = tpData.FormatCodeToFile(dstFilePath, bf)
	if err != nil {
		panic(err)
	}

	fmt.Println(cgm.ModelName + " Success")
}
