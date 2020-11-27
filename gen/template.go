package gen

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/fatih/structtag"
	"github.com/illidan33/tools/common"
	utiltypes "github.com/illidan33/tools/gen/util/types"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	gotypes "go/types"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type GenTemplate struct {
	TemplateMapFuncs map[string]interface{}
}

type TemplateModelField struct {
	Name     string
	GormName string
	JsonName string
	Type     string
	Default  string
	Tag      string
	Comment  string
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
	if gt.TemplateMapFuncs == nil {
		gt.TemplateMapFuncs = map[string]interface{}{}
	}
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
		e = fmt.Errorf("ParseTemplate - parse error: %s\n", e.Error())
		return
	}
	e = tp.Execute(templateSource, templateData)
	if e != nil {
		e = fmt.Errorf("ParseTemplate - execute template data error: %s\n", e.Error())
		return
	}
	return
}

func (gt *GenTemplate) FormatCodeToFile(filePath string, templateData *bytes.Buffer) (err error) {
	dir := filepath.Dir(filePath)
	if !common.IsExists(dir) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}
	filePath, _ = filepath.Abs(filePath)

	var file *os.File
	file, err = os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return errors.New("FormatCodeToFile - open file error: " + err.Error())
	}
	defer file.Close()

	// write temporary data to file first
	file.Write(templateData.Bytes())

	f, e := decorator.Parse(templateData.String())
	if e != nil {
		err = fmt.Errorf("FormatCodeToFile - parse error: %s", e.Error())
		return
	}
	// clear file content
	file.Truncate(0)
	file.Seek(0, 0)

	err = decorator.Fprint(file, f)
	if err != nil {
		return errors.New("FormatCodeToFile - write file error: " + err.Error())
	}
	return
}

func (gt *GenTemplate) GetAstTree(filePath string) (*token.FileSet, *ast.File, error) {
	fset := token.NewFileSet()
	astfile, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, nil, errors.New("GetAstTree - ParseFile error: " + err.Error())
	}
	return fset, astfile, nil
}

func (gt *GenTemplate) GetDstTree(filePath string) (*dst.File, error) {
	var file *os.File
	var err error
	if !common.IsExists(filePath) {
		return nil, fmt.Errorf("GetDstTree - file not exist: %s", filePath)
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
		return nil, fmt.Errorf("GetDstTree - file empty: %s", filePath)
	}

	f, err := decorator.Parse(codes)
	if err != nil {
		return nil, errors.New("GetDstTree - decorator parse error: " + err.Error())
	}
	return f, nil
}

func (gt *GenTemplate) GetTypesPackage(filePath string) (*types.Package, error) {
	pkg, err := importer.For("source", nil).Import(filePath)
	if err != nil {
		return nil, err
	}

	return pkg, nil
}

type TemplateModel struct {
	ModelName           string
	ModelComment        string
	Package             string
	Type                string // like []Model or empty
	TemplateModelFields []TemplateModelField
}

func (tm *TemplateModel) parseTagToTokens(s string) (rs []string, e error) {
	rs = make([]string, 0)
	tmp := bytes.Buffer{}
	keyS := false
	for i := 0; i < len(s); i++ {
		if s[i] == ' ' {
			ts := tmp.String()
			if keyS {
				tmp.WriteByte(s[i])
			} else if ts != "" {
				rs = append(rs, tmp.String())
				tmp = bytes.Buffer{}
			}
		} else {
			if s[i] == '"' {
				if keyS == false {
					keyS = true
				} else {
					keyS = false
				}
			}
			tmp.WriteByte(s[i])
		}
	}
	if tmp.String() != "" {
		rs = append(rs, tmp.String())
		tmp = bytes.Buffer{}
	}
	return
}

func (tm *TemplateModel) parseTypesVar(v *gotypes.Var, tag string) []TemplateModelField {
	if v.Embedded() {
		t := v.Type()
		str := make([]TemplateModelField, 0)
		if ts, ok := t.Underlying().(*gotypes.Struct); ok {
			for i := 0; i < ts.NumFields(); i++ {
				tmp := tm.parseTypesVar(ts.Field(i), ts.Tag(i))
				str = append(str, tmp...)
			}
		}
		return str
	} else {
		templateField := TemplateModelField{
			Name:    v.Name(),
			Type:    v.Type().String(),
			Tag:     tag,
			Comment: "",
		}
		return []TemplateModelField{templateField}
	}
}

func (tm *TemplateModel) ParseDstTree(file *dst.File) (modelList []TemplateModel, err error) {
	for _, i := range file.Decls {
		gd, ok := i.(*dst.GenDecl)
		if !ok {
			continue
		}
		if gd.Tok != token.TYPE {
			continue
		}
		model := TemplateModel{}
		tf, ok := gd.Specs[0].(*dst.TypeSpec)
		if !ok {
			err = fmt.Errorf("can not change to TypeSpec: %#v", gd.Specs)
			return
		}
		// this entity model name
		model.ModelName = tf.Name.Name
		model.Package = file.Name.Name

		fieldMap := map[string]TemplateModelField{}
		if len(model.TemplateModelFields) == 0 {
			st, ok := tf.Type.(*dst.StructType)
			if !ok {
				continue
			}

			for _, field := range st.Fields.List {
				templateField := TemplateModelField{}

				if len(field.Names) > 0 {
					templateField.Name = field.Names[0].Name
					templateField.Tag = strings.Trim(field.Tag.Value, "`")
					tags, e := structtag.Parse(templateField.Tag)
					if e != nil {
						err = e
						return
					}
					for _, tag := range tags.Tags() {
						switch tag.Key {
						case utiltypes.MODEL_TAG_TYPE__JSON:
							templateField.JsonName = tag.Name
						case utiltypes.MODEL_TAG_TYPE__GORM:
							templateField.GormName = tag.Name
						}
					}
				}

				if len(field.Decs.NodeDecs.Start) > 0 {
					templateField.Comment = field.Decs.NodeDecs.Start[0]
				}
				if len(field.Decs.End) > 0 {
					templateField.Comment = field.Decs.End[0]
				}

				IdentType, ok := field.Type.(*dst.Ident)
				if ok {
					templateField.Type = IdentType.Name
				} else {
					ok = true
				}
				ExprType, ok := field.Type.(*dst.SelectorExpr)
				if ok {
					ExprXType, ok := ExprType.X.(*dst.Ident)
					if ok {
						templateField.Type = ExprXType.Name
					}
					templateField.Type += "." + ExprType.Sel.Name
				}
				fieldMap[templateField.Name] = templateField
				model.TemplateModelFields = append(model.TemplateModelFields, templateField)
			}
		}
		// comment def of struct
		if gd.Decs.NodeDecs.Start != nil {
			model.ModelComment = strings.Join(gd.Decs.NodeDecs.Start.All(), ",")
		}

		modelList = append(modelList, model)
	}
	return
}

func (tm *TemplateModel) ParseModelFile(filepath string) (modelList []TemplateModel, err error) {
	pkg, e := common.GetImportPackageName(filepath)
	if e != nil {
		err = e
		return
	}
	tm.Package = pkg
	var dstFile *dst.File
	dstFile, err = DefaultGenTemplate.GetDstTree(filepath)
	if err != nil {
		return
	}
	modelList, err = tm.ParseDstTree(dstFile)
	return
}

func (tm *TemplateModel) ParseModelPackage(packagePath string) (modelList []TemplateModel, err error) {
	dirPath, e := common.GetImportByPackage(packagePath)
	if e != nil {
		err = e
		return
	}

	modelList, err = tm.ParseModelDir(dirPath)
	return
}

func (tm *TemplateModel) ParseModelDir(dirPath string) (modelList []TemplateModel, err error) {
	var files []os.FileInfo
	files, err = ioutil.ReadDir(dirPath)
	if err != nil {
		return
	}
	for _, file := range files {
		if file.IsDir() {
			mList, e := tm.ParseModelDir(filepath.Join(dirPath, file.Name()))
			if e != nil {
				err = e
				return
			}
			modelList = append(modelList, mList...)
		} else if filepath.Ext(file.Name()) == ".go" {
			mList, e := tm.ParseModelFile(filepath.Join(dirPath, file.Name()))
			if e != nil {
				err = e
				return
			}
			modelList = append(modelList, mList...)
		}
	}
	return
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
