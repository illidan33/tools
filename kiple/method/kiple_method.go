package method

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/illidan33/tools/common"
)

type CmdKipleMethod struct {
	ModelName string
	Entity    string
	IsDebug   bool
}

func (cgm *CmdKipleMethod) CmdHandle() {
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
	tpData.InterfaceName = cgm.ModelName
	tpData.EntityPath = cgm.Entity
	tpData.ModelName = cgm.ModelName + "Impl"

	// for test
	if cgm.IsDebug {
		os.Setenv("GOFILE", "user_dao_impl.go")
		os.Setenv("GOPACKAGE", "model")
		exeFilePath = os.Getenv("GOPATH") + "/src/github.com/illidan33/tools/example/model"
		cmdFile.CmdFileName = "user_profiles_dao.go"
		tpData.PackageName = "model"
	}

	excuteFilePath := fmt.Sprintf("%s/%s", exeFilePath, cmdFile.CmdFileName)
	dstFilePath := fmt.Sprintf("%s/%s", exeFilePath, cmdFile.CmdFileName)

	// update exist interface
	if common.IsExists(excuteFilePath) {
		dstfl, err := tpData.GetDstTree(excuteFilePath)
		if err != nil {
			panic(err)
		}
		flag, err := tpData.FindInterfaceAndFillMethods(dstfl, dstFilePath)
		if err != nil {
			panic(err)
		}
		if flag {
			fmt.Println("Update exist interface success")
			return
		}
	}

	// create new interface
	filePath, err := filepath.Abs(cgm.Entity)
	if err != nil {
		panic(errors.New("can not parse source to abs filepath"))
	}
	err = tpData.Parse(filePath, cgm.IsDebug)
	if err != nil {
		panic(err)
	}

	bf, err := tpData.ParseTemplate(templateMethodTxt, tpData.ModelName, tpData)
	if err != nil {
		panic(err)
	}

	err = tpData.FormatCodeToFile(dstFilePath, bf)
	if err != nil {
		panic(err)
	}

	fmt.Println(cgm.ModelName + " Success")
}
