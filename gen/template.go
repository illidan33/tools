package gen

import (
	"bytes"
	"fmt"
	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
)

type CmdGen interface {
	CmdHandle(cwd string)
}

type TemplateGenStruct struct {
	PackageName   string
	PackageList   map[string]string
	StructName    string
	StructComment string
	TemplateFuncs map[string]interface{}
}

type TemplateGenStructField struct {
	Name    string
	Type    string
	Tag     string
	Comment string
}

func (tfm *TemplateGenStruct) Registe(data map[string]interface{}) {
	for name, fc := range data {
		tfm.TemplateFuncs[name] = fc
	}
}

func (tfm *TemplateGenStruct) Init() {
	d := map[string]interface{}{
		"var":  func(s string) string { return ToLowerCamelCase(s) },
		"type": func(s string) string { return ToUpperCamelCase(s) },
		"snake": func(s string) string { return ToLowerSnakeCase(s) },
		"printf": func(s string, args ...interface{}) string {
			return fmt.Sprintf(s, args...)
		},
		"html": func(s string) template.HTML {
			return template.HTML(s)
		},
	}
	tfm.Registe(d)
}

func (tfm *TemplateGenStruct) ParseTemplate(templateTxt string, templateData interface{}) (templateSource *bytes.Buffer, e error) {
	templateSource = bytes.NewBuffer([]byte(""))

	tp := template.New("sql_to_modle")
	tp.Funcs(tfm.TemplateFuncs)
	if tp, e = tp.Parse(templateTxt); e != nil {
		return
	}
	e = tp.Execute(templateSource, templateData)
	return
}

func (tfm *TemplateGenStruct) FormatCodeToFile(filePath string, templateSource *bytes.Buffer) (err error) {
	var file *os.File
	filePath, _ = filepath.Abs(filePath)
	file, err = os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0744)
	if err != nil {
		return err
	}
	defer file.Close()

	f, e := decorator.Parse(templateSource.String())
	if e != nil {
		err = e
		return
	}
	err = decorator.Fprint(file, f)
	return
}

func (tfm *TemplateGenStruct) GetDstTree(filePath string) (*dst.File, error) {
	var file *os.File
	var err error
	if IsExists(filePath) {
		file, err = os.Open(filePath)
		if err != nil {
			return nil, err
		}
	} else {
		file, err = os.Create(filePath)
		if err != nil {
			return nil, err
		}
	}
	defer file.Close()
	templateSource, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	f, err := decorator.Parse(templateSource)
	if err != nil {
		panic(err)
	}
	return f, nil
}

func (tfm *TemplateGenStruct) SortFields(fields map[string]GormField) []GormField {
	arr := make([]GormField, len(fields))
	i := 0
	for _, field := range fields {
		arr[i] = field
	}
	sort.Slice(arr, func(i, j int) bool {
		if arr[i].Sort < arr[j].Sort {
			return true
		}
		return false
	})
	return arr
}
