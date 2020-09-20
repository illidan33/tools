package gen

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"html/template"
	"io/ioutil"
	"myprojects/tools/common"
	"os"
	"path/filepath"
)

type GenTemplate struct {
	TemplateMapFuncs map[string]interface{}
}

type TemplateModel struct {
	ModelName           string
	ModelComment        string
	TemplateModelFields []TemplateModelField
}

type TemplatePackage struct {
	PackageName string
	PackageList []string
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
		"type": func(s string) string { return common.ToUpperCamelCase(s) },
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

func (gt *GenTemplate) ParseFilePath() (path common.CmdFilePath, err error) {
	path.CmdFileName = os.Getenv("GOFILE")
	//path.CmdFileName = "booking.go"
	if path.CmdFileName == "" {
		err = errors.New("fileName empty")
		return
	}
	path.PackageName = os.Getenv("GOPACKAGE")
	// TODO(illidan/2020/9/20): remove
	//path.PackageName = "main"
	if path.PackageName == "" {
		err = errors.New("packageName empty")
		return
	}
	return
}

func (gt *GenTemplate) ParseTemplate(templateTxt string, templateName string, templateData interface{}) (templateSource *bytes.Buffer, e error) {
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
	var file *os.File
	filePath, _ = filepath.Abs(filePath)
	file, err = os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	f, e := decorator.Parse(templateData.String())
	if e != nil {
		err = e
		return
	}
	err = decorator.Fprint(file, f)
	return
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

	f, err := decorator.Parse(codes)
	if err != nil {
		panic(err)
	}
	return f, nil
}
