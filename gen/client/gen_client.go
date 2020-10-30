package client

import (
	"errors"
	"fmt"
	"github.com/illidan33/tools/common"
	"os"
)

type CmdGenClient struct {
	DocUrl      string
	ServiceName string
	IsDebug     bool
}

func (cgc CmdGenClient) CmdHandle() {
	tpData := TemplateGenClient{}

	if cgc.DocUrl == "" {
		panic(errors.New("DocUrl required"))
	}
	if cgc.ServiceName == "" {
		panic(errors.New("ServiceName required"))
	}

	//package:= tpData.GetGenEnvironmentValues()
	tpData.PackageName = "client_" + common.ToLowerSnakeCase(cgc.ServiceName)
	tpData.ClientModel.ModelName = common.ToUpperCamelCase(tpData.PackageName)

	err := tpData.Parse(cgc.DocUrl, cgc.IsDebug)
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

	fmt.Println(cgc.ServiceName + " Success")
}
