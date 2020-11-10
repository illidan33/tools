package dao

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"os"
	"path/filepath"

	"github.com/illidan33/tools/common"
	"github.com/illidan33/tools/gen"
)

type CmdKipleInterfaceCheck struct {
	InterfaceName string
	IsDebug       bool

	Environments common.CmdFilePath
	gen.GenTemplate
	gen.TemplateModel
}

func (tpData *CmdKipleInterfaceCheck) String() string {
	return tpData.InterfaceName
}

func (tpData *CmdKipleInterfaceCheck) Init() error {
	tpData.InitTemplateFuncs()

	var err error
	tpData.Environments, err = common.GetGenEnvironmentValues()
	if err != nil {
		return err
	}
	if tpData.Environments.CmdFileName == "" {
		tpData.Environments.CmdFileName = fmt.Sprintf("%s.go", common.ToLowerSnakeCase(tpData.InterfaceName))
	}

	// for test
	if tpData.IsDebug {
		os.Setenv("GOFILE", "user_dao_impl.go")
		os.Setenv("GOPACKAGE", "model")
		tpData.Environments.CmdDir = filepath.Join(common.GetGoPath(), "/src/github.com/illidan33/gotest/tools_test/example/model")
		tpData.Environments.CmdFileName = "user_profiles_dao.go"
	}
	return nil
}

func (tpData *CmdKipleInterfaceCheck) Parse() error {
	excuteFilePath := filepath.Join(tpData.Environments.CmdDir, tpData.Environments.CmdFileName)
	if !common.IsExists(excuteFilePath) {
		panic(errors.New("file not exist: " + excuteFilePath))
	}
	fset, dstfl, err := tpData.GetAstTree(excuteFilePath)
	if err != nil {
		panic(err)
	}
	err = tpData.FindInterfaceAndFillMethods(fset, dstfl, excuteFilePath)
	if err != nil {
		panic(err)
	}

	return nil
}

func (tm *CmdKipleInterfaceCheck) FindInterfaceAndFillMethods(fset *token.FileSet, dstfile *ast.File, dstFilePath string) error {
	var interfaceNode *ast.InterfaceType
	//userdaoFuncMap := map[string]*ast.Field{}
	for _, decl := range dstfile.Decls {
		if declv, ok := decl.(*ast.GenDecl); ok && declv.Tok == token.TYPE {
			if len(declv.Specs) == 0 {
				return errors.New("FindInterfaceAndFillMethods - GenDecl has no Specs")
			}
			if typespec, ok := declv.Specs[0].(*ast.TypeSpec); ok && typespec.Name.Name == tm.InterfaceName {
				if interfaceNode, ok = typespec.Type.(*ast.InterfaceType); ok {
					//for _, field := range interfaceNode.Methods.List {
					//	userdaoFuncMap[field.Names[0].Name] = field
					//}
					break
				}
			}
		}
	}
	if interfaceNode == nil {
		return errors.New(fmt.Sprintf("Interface %s not found", tm.InterfaceName))
	}

	newList := make([]*ast.Field, 0)
	for i := 0; i < len(dstfile.Decls); i++ {
		decl := dstfile.Decls[i]
		if ffv, ok := decl.(*ast.FuncDecl); ok && ffv.Recv != nil {
			if len(ffv.Recv.List) == 0 || ffv.Recv.List[0].Type == nil {
				continue
			}
			if ffvse, ok := (ffv.Recv.List[0].Type).(*ast.StarExpr); ok && ffvse.X.(*ast.Ident).Name == tm.ModelName {
				ffnew := ast.Field{}
				ffnew.Names = []*ast.Ident{ast.NewIdent(ffv.Name.Name)}
				ffnew.Type = ffv.Type
				newList = append(newList, &ffnew)
			}
		}
	}

	interfaceNode.Methods.List = newList

	cmap := ast.NewCommentMap(fset, dstfile, dstfile.Comments)
	dstfile.Comments = cmap.Filter(dstfile).Comments()

	var output []byte
	buffer := bytes.NewBuffer(output)
	err := format.Node(buffer, fset, dstfile)
	if err != nil {
		return err
	}
	var file *os.File
	file, err = os.OpenFile(dstFilePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	file.Write(buffer.Bytes())

	return nil
}
