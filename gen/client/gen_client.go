package client

import (
	"errors"
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
	tpData.PackageName = common.ToLowerSnakeCase("client_" + cgc.CmdGenClientServiceName)
	tpData.ClientModel.ModelName = common.ToUpperCamelCase(tpData.PackageName)

	//test
	//cgc.CmdGenClientDocUrl = "http://192.168.1.116:8080/swagger/swagger/doc.json"

	if cgc.CmdGenClientDocUrl == "" {
		panic(errors.New("required doc url"))
	}
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
}
