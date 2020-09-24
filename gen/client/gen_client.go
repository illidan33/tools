package client

import (
	"fmt"
	"myprojects/tools/common"
	"os"
)

type CmdGenClient struct {
	CmdGenClientDocUrl      string
	CmdGenClientServiceName string
}

func (cgc CmdGenClient) CmdHandle() {
	tpData := TemplateGenClient{}
	tpData.InitTemplateFuncs()

	//package:= tpData.ParseFilePath()
	tpData.PackageName = common.ToLowerSnakeCase("client_" + common.ToLowerSnakeCase(cgc.CmdGenClientServiceName))
	tpData.ClientModel.ModelName = common.ToUpperCamelCase(tpData.PackageName)

	err := tpData.Parse(cgc.CmdGenClientDocUrl)
	if err != nil {
		panic(err)
	}

	exeFilePath, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	err = tpData.ParseTemplateAndFormatToFile(exeFilePath)
	if err != nil {
		panic(err)
	}

	fmt.Println("Success")
}
