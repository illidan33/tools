package method

import (
	"errors"
	"fmt"
	"github.com/illidan33/tools/common"
	"os"
)

type CmdGenMethod struct {
	ModelName string
	IsDebug   bool
}

func (cgm *CmdGenMethod) CmdHandle() {
	tpData := TemplateDataMethod{}
	tpData.InitTemplateFuncs()

	cmdFile, err := common.GetGenEnvironmentValues()
	if err != nil {
		panic(err)
	}
	exeFilePath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	tpData.PackageName = cmdFile.PackageName
	tpData.ModelName = cgm.ModelName

	// for test
	if cgm.IsDebug {
		os.Setenv("GOFILE", "mp_orders.go")
		os.Setenv("GOPACKAGE", "model")
		exeFilePath = os.Getenv("GOPATH") + "/src/github.com/illidan33/tools/example/model"
		cmdFile.CmdFileName = "mp_orders.go"
		tpData.PackageName = "model"
	}

	filePath := fmt.Sprintf("%s/%s", exeFilePath, cmdFile.CmdFileName)
	err = tpData.Parse(filePath, cgm.IsDebug)
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
