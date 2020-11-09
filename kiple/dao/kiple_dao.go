package dao

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/dave/dst"
	"github.com/illidan33/tools/common"
	"github.com/illidan33/tools/gen"
	"os"
	"path/filepath"
)

type CmdKipleDao struct {
	InterfaceName string
	Entity        string
	IsDebug       bool

	Environments      common.CmdFilePath
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

func (tpData *CmdKipleDao) String() string {
	return tpData.InterfaceName
}
func (tpData *CmdKipleDao) Init() (err error) {
	tpData.InitTemplateFuncs()

	tpData.Environments, err = common.GetGenEnvironmentValues()
	if err != nil {
		return
	}
	fmt.Printf("environment: %#v\n", tpData.Environments)

	tpData.PackageName = tpData.Environments.PackageName
	tpData.ModelName = tpData.InterfaceName + implNameFlag

	// for test
	if tpData.IsDebug {
		os.Setenv("GOFILE", "user_dao_impl.go")
		os.Setenv("GOPACKAGE", "model")
		tpData.Environments.CmdDir = filepath.Join(common.GetGoPath(), "/src/github.com/illidan33/tools/example/entity")
	}

	return
}

func (tpData *CmdKipleDao) Parse() (err error) {
	var path string
	path, err = filepath.Abs(tpData.Entity)
	if err != nil {
		err = errors.New("can not parse source to abs filepath")
		return
	}

	pathDir := filepath.Dir(path)
	var pathPackage string
	pathPackage, err = common.GetImportPackageName(pathDir)
	if err != nil {
		return
	}
	if tpData.PackageName != pathPackage {
		pkgPath, err := common.GetImportPath(pathDir)
		if err != nil {
			return err
		}
		tpData.AddPackage("entity", pkgPath)
	}

	var dstTree *dst.File
	dstTree, err = tpData.GetDstTree(path)
	if err != nil {
		return err
	}
	if err = tpData.ParseDstTree(dstTree); err != nil {
		return
	}
	if err = tpData.ParseIndexToMethod(); err != nil {
		return
	}

	var bf *bytes.Buffer
	bf, err = tpData.ParseTemplate(templateMethodTxt, tpData.ModelName, tpData)

	dstFilePath := filepath.Join(tpData.Environments.CmdDir, common.ToLowerSnakeCase(tpData.InterfaceName)+".go")
	err = tpData.FormatCodeToFile(dstFilePath, bf)
	if err != nil {
		return
	}

	return
}
