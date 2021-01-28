package template

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/fatih/structtag"
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
	"tools/common"
	utiltypes "tools/gen/util/types"
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

type TemplateGenFunc struct {
	Name           string
	Comment        string
	BelongToStruct *TemplateModel
	Package        string
	Args           []*TemplateModelField
	Returns        []*TemplateModelField
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

func (gt *GenTemplate) WriteToFile(filePath string, data []byte) (err error) {
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

	// clear file content
	file.Truncate(0)
	file.Seek(0, 0)

	// write temporary data to file first
	_, err = file.Write(data)
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
	if filepath.IsAbs(filePath) {
		//var err error
		filePat, err := common.GetBuildPackageFromDir(filePath)
		if err != nil {
			return nil, err
		}
		filePath = filePat.ImportPath
	}
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
	Type                string
	TemplateModelFields []TemplateModelField
	TemplateModelFuncs  []TemplateGenFunc
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

func (tm *TemplateModel) parseModelFieldFromDst(field *dst.Field) (templateField TemplateModelField, err error) {
	if len(field.Names) > 0 {
		templateField.Name = field.Names[0].Name
		if field.Tag != nil {
			templateField.Tag = strings.Trim(field.Tag.Value, "`")
		}
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
	templateField.Type = tm.ParseDstNodeType(field.Type)
	return
}

func (tm *TemplateModel) getCommentByFilterTag(arr []string) string {
	buf := bytes.Buffer{}
	for _, s := range arr {
		if !strings.HasPrefix(s, "// @") && !strings.HasPrefix(s, "//@") {
			buf.WriteString(s)
		}
	}
	return buf.String()
}

func (tm *TemplateModel) ParseDstCommentFromNode(node dst.NodeDecs, isFilterTag bool) string {
	comment := make([]string, 0)
	if node.Start != nil {
		if !isFilterTag {
			comment = append(comment, node.Start.All()...)
		} else {
			comment = append(comment, tm.getCommentByFilterTag(node.Start.All()))
		}
	}
	if node.End != nil {
		if !isFilterTag {
			comment = append(comment, node.End.All()...)
		} else {
			comment = append(comment, tm.getCommentByFilterTag(node.End.All()))
		}
	}
	return strings.Join(comment, ",")
}

func (tm *TemplateModel) ParseDstNodeType(tp dst.Expr) string {
	switch tpVal := tp.(type) {
	case *dst.Ident:
		return tpVal.Name
	case *dst.StarExpr:
		return tm.ParseDstNodeType(tpVal.X)
	case *dst.ArrayType:
		return "[]" + tm.ParseDstNodeType(tpVal.Elt)
	case *dst.SliceExpr:
		return "[]" + tm.ParseDstNodeType(tpVal.X)
	case *dst.SelectorExpr:
		name := ""
		if tpVal.X != nil {
			name = tpVal.X.(*dst.Ident).Name
		}
		return name + "." + tpVal.Sel.Name
	case *dst.MapType:
		return fmt.Sprintf("map[%s]%s", tm.ParseDstNodeType(tpVal.Key), tm.ParseDstNodeType(tpVal.Value))
	default:
		// TODO(illidan/2020/11/30): other types
	}
	return ""
}

func (tm *TemplateModel) ParseDstTree(file *dst.File) (modelList []TemplateModel, err error) {
	for _, decl := range file.Decls {
		gd, ok := decl.(*dst.GenDecl)
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
		st, ok := tf.Type.(*dst.StructType)
		if !ok {
			continue
		}

		for _, field := range st.Fields.List {
			templateField, e := tm.parseModelFieldFromDst(field)
			if e != nil {
				err = e
				return
			}
			fieldMap[templateField.Name] = templateField
			model.TemplateModelFields = append(model.TemplateModelFields, templateField)
		}

		// comment def of struct
		if gd.Decs.NodeDecs.Start != nil {
			model.ModelComment = strings.Join(gd.Decs.NodeDecs.Start.All(), ",")
		}

		modelList = append(modelList, model)
	}

	for _, decl := range file.Decls {
		if ffv, ok := decl.(*dst.FuncDecl); ok {
			if len(ffv.Recv.List) == 0 || ffv.Recv.List[0].Type == nil {
				continue
			}
			fmt.Println(ffv)
		}
	}
	return
}

func (tm *TemplateModel) ParseFuncDstTree(file *dst.File) (funcList []TemplateGenFunc, err error) {
	for _, decl := range file.Decls {
		if ffv, ok := decl.(*dst.FuncDecl); ok {
			fc := TemplateGenFunc{
				Name:           "",
				Comment:        "",
				BelongToStruct: nil,
				Package:        "",
				Args:           nil,
				Returns:        nil,
			}
			if ffv.Recv != nil && len(ffv.Recv.List) != 0 {
				continue
			}
			if ffv.Type.Params.List != nil && len(ffv.Type.Params.List) > 0 {

			}
			fmt.Println(fc)
		}
	}
	return
}

func (tm *TemplateModel) ParseModelDstTree(file *dst.File, parseFuncs bool) (modelList map[string]TemplateModel, err error) {
	modelList = map[string]TemplateModel{}
	for _, decl := range file.Decls {
		gd, ok := decl.(*dst.GenDecl)
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
		switch tfTyp := tf.Type.(type) {
		case *dst.StructType:
			for _, field := range tfTyp.Fields.List {
				templateField, e := tm.parseModelFieldFromDst(field)
				if e != nil {
					err = e
					return
				}
				fieldMap[templateField.Name] = templateField
				model.TemplateModelFields = append(model.TemplateModelFields, templateField)
			}
			model.Type = model.ModelName
		default:
			model.Type = tm.ParseDstNodeType(tf.Type)
		}

		// comment def of struct
		if gd.Decs.NodeDecs.Start != nil {
			model.ModelComment = strings.Join(gd.Decs.NodeDecs.Start.All(), "\n")
		}

		modelList[model.Package+"."+model.ModelName] = model
	}
	if parseFuncs {
		for _, decl := range file.Decls {
			if ffv, ok := decl.(*dst.FuncDecl); ok {
				fc := TemplateGenFunc{
					Name:           ffv.Name.Name,
					Comment:        "",
					BelongToStruct: nil,
					Package:        "",
					Args:           nil,
					Returns:        nil,
				}
				if ffv.Recv == nil || len(ffv.Recv.List) == 0 {
					continue
				}
				belongModelName := tm.ParseDstNodeType(ffv.Recv.List[0].Type)

				if ffv.Type.Params.List != nil && len(ffv.Type.Params.List) > 0 {
					for _, field := range ffv.Type.Params.List {
						if len(field.Names) > 0 {
							for _, name := range field.Names {
								templateField := TemplateModelField{
									Name:     name.Name,
									GormName: "",
									JsonName: "",
									Type:     "",
									Default:  "",
									Tag:      "",
									Comment:  "",
								}
								templateField.Type = tm.ParseDstNodeType(field.Type)
								fc.Args = append(fc.Args, &templateField)
							}
						} else {
							templateField := TemplateModelField{
								Name:     "",
								GormName: "",
								JsonName: "",
								Type:     field.Type.(*dst.Ident).Name,
								Default:  "",
								Tag:      "",
								Comment:  "",
							}
							templateField.Type = tm.ParseDstNodeType(field.Type)
							fc.Args = append(fc.Args, &templateField)
						}
					}
				}
				if ffv.Type.Results != nil && ffv.Type.Results.List != nil && len(ffv.Type.Results.List) > 0 {
					for _, field := range ffv.Type.Results.List {
						if len(field.Names) > 0 {
							for _, name := range field.Names {
								templateField := TemplateModelField{
									Name:     name.Name,
									GormName: "",
									JsonName: "",
									Type:     "",
									Default:  "",
									Tag:      "",
									Comment:  "",
								}
								templateField.Type = tm.ParseDstNodeType(field.Type)
								fc.Args = append(fc.Args, &templateField)
							}
						} else {
							templateField := TemplateModelField{
								Name:     "",
								GormName: "",
								JsonName: "",
								Type:     field.Type.(*dst.Ident).Name,
								Default:  "",
								Tag:      "",
								Comment:  "",
							}
							templateField.Type = tm.ParseDstNodeType(field.Type)
							fc.Args = append(fc.Args, &templateField)
						}
					}
				}

				if tmp, ok := modelList[belongModelName]; ok {
					tmp.TemplateModelFuncs = append(tmp.TemplateModelFuncs, fc)
					modelList[belongModelName] = tmp
				}
			}
		}
	}
	return
}

func (tm *TemplateModel) ParseModelsFromFile(filepath string, parseFunc bool) (modelList map[string]TemplateModel, err error) {
	//pkg, e := common.GetPackageNameFromPath(filepath)
	//if e != nil {
	//	err = e
	//	return
	//}
	//tm.Package = pkg
	var dstFile *dst.File
	dstFile, err = DefaultGenTemplate.GetDstTree(filepath)
	if err != nil {
		return
	}
	modelList, err = tm.ParseModelDstTree(dstFile, parseFunc)
	return
}

func (tm *TemplateModel) ParseModelDir(dirPath string, parseFunc bool) (modelList map[string]TemplateModel, err error) {
	var files []os.FileInfo
	modelList = map[string]TemplateModel{}
	files, err = ioutil.ReadDir(dirPath)
	if err != nil {
		return
	}
	for _, file := range files {
		if file.IsDir() {
			mList, e := tm.ParseModelDir(filepath.Join(dirPath, file.Name()), parseFunc)
			if e != nil {
				err = e
				return
			}
			for k, model := range mList {
				modelList[k] = model
			}
		} else if filepath.Ext(file.Name()) == ".go" {
			mList, e := tm.ParseModelsFromFile(filepath.Join(dirPath, file.Name()), parseFunc)
			if e != nil {
				err = e
				return
			}
			for k, model := range mList {
				modelList[k] = model
			}
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
