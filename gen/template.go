package gen

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/illidan33/tools/common"
	"go/ast"
	"go/parser"
	"go/token"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
)

type GenTemplate struct {
	TemplateMapFuncs map[string]interface{}
}

type TemplateModel struct {
	ModelName           string
	ModelComment        string
	Type                string // like []Model or empty
	TemplateModelFields []TemplateModelField
}

type TemplateModelField struct {
	Name    string
	Type    string
	Tag     string
	Comment string
}

type TemplateGenModelFunc struct {
	Name           string
	Comment        string
	BelongToStruct TemplateModel
	Args           []TemplateModelField
	Returns        []TemplateModelField
}

var DefaultGenTemplate = &GenTemplate{}

func init() {
	DefaultGenTemplate.InitTemplateFuncs()
}

func (gt *GenTemplate) RegisteTemplateFunc(data map[string]interface{}) {
	for name, fc := range data {
		gt.TemplateMapFuncs[name] = fc
	}
}

func (gt *GenTemplate) InitTemplateFuncs() {
	if gt.TemplateMapFuncs == nil {
		gt.TemplateMapFuncs = map[string]interface{}{}
	}
	d := map[string]interface{}{
		"var":   func(s string) string { return common.ToLowerCamelCase(s) },
		"type":  func(s string) string { return common.ToUpperCamelCase(s) },
		"snake": func(s string) string { return common.ToLowerSnakeCase(s) },
		"printf": func(s string, args ...interface{}) string {
			return fmt.Sprintf(s, args...)
		},
		"html": func(s string) template.HTML {
			return template.HTML(s)
		},
	}
	gt.RegisteTemplateFunc(d)
}

func (gt *GenTemplate) ParseTemplate(templateTxt string, templateName string, templateData interface{}, templateFuncMap ...map[string]interface{}) (templateSource *bytes.Buffer, e error) {
	if len(gt.TemplateMapFuncs) == 0 {
		gt.InitTemplateFuncs()
	}
	if len(templateFuncMap) > 0 {
		for _, funcMap := range templateFuncMap {
			gt.RegisteTemplateFunc(funcMap)
		}
	}

	templateSource = &bytes.Buffer{}
	tp := template.New(templateName)
	tp.Funcs(gt.TemplateMapFuncs)
	tp, e = tp.Parse(templateTxt)
	if e != nil {
		return
	}
	e = tp.Execute(templateSource, templateData)
	return
}

func (gt *GenTemplate) FormatCodeToFile(filePath string, templateData *bytes.Buffer) (err error) {
	filePath, _ = filepath.Abs(filePath)

	f, e := decorator.Parse(templateData.String())
	if e != nil {
		fmt.Println(templateData.String())
		err = fmt.Errorf("format code error: %s", e.Error())
		return
	}

	var file *os.File
	file, err = os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	err = decorator.Fprint(file, f)
	return
}

func (gt *GenTemplate) GetAstTree(filePath string) (*token.FileSet, *ast.File, error) {
	fset := token.NewFileSet()
	astfile, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, nil, err
	}
	return fset, astfile, nil
}

func (gt *GenTemplate) GetDstTree(filePath string) (*dst.File, error) {
	var file *os.File
	var err error
	if !common.IsExists(filePath) {
		return nil, errors.New("file not exist")
	}
	file, err = os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	codes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	if len(codes) == 0 {
		return nil, errors.New("file empty")
	}

	f, err := decorator.Parse(codes)
	if err != nil {
		return nil, errors.New("decorator parse error: " + err.Error())
	}
	return f, nil
}

type TemplatePackage struct {
	PackageName string
	PackageList map[string]string
}

func (tpkg *TemplatePackage) AddPackage(name, val string) {
	if tpkg.PackageList == nil {
		tpkg.PackageList = map[string]string{}
	}
	tpkg.PackageList[name] = val
}
