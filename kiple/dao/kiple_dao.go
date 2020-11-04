package dao

import (
	"fmt"
	"github.com/illidan33/tools/common"
	"github.com/illidan33/tools/gen"
	"os"
	"path/filepath"
)

type CmdKipleDao struct {
	InterfaceName string
	Entity        string
	IsDebug       bool

	EntityName        string
	EntityPackageName string
	gen.GenTemplate
	gen.TemplatePackage
	gen.TemplateModel
	CmdKipleDaoFuncNames []string
	CmdKipleDaoFuncs     []string
	CmdKipleDaoIndexs    []CmdKipleDaoIndex
}

const implNameFlag = "Impl"

func (tpData *CmdKipleDao) CmdHandle() {
	tpData.InitTemplateFuncs()

	environValues, err := common.GetGenEnvironmentValues(tpData.IsDebug)
	if err != nil {
		panic(err)
	}

	tpData.PackageName = environValues.PackageName
	tpData.ModelName = tpData.InterfaceName + implNameFlag
	tpData.PackageList = map[string]string{}

	// for test
	if tpData.IsDebug {
		os.Setenv("GOFILE", "user_dao_impl.go")
		os.Setenv("GOPACKAGE", "model")
		environValues.CmdDir = filepath.Join(common.GetGoPath(), "/src/github.com/illidan33/tools/example/model")
		environValues.CmdFileName = "user_profiles_dao.go"
		tpData.PackageName = "model"
	}

	// create new interface
	err = tpData.Parse(tpData.Entity)
	if err != nil {
		panic(err)
	}

	bf, err := tpData.ParseTemplate(templateMethodTxt, tpData.ModelName, tpData)
	if err != nil {
		panic(err)
	}

	dstFilePath := filepath.Join(environValues.CmdDir, common.ToLowerSnakeCase(tpData.InterfaceName)+".go")
	err = tpData.FormatCodeToFile(dstFilePath, bf)
	if err != nil {
		panic(err)
	}

	fmt.Println(tpData.ModelName + " Success")
}
