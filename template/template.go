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
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"tools/common"
	ttypes "tools/template/types"
)

type GenTemplate struct {
	TemplateMapFuncs map[string]interface{}
}

type TemplateModelField struct {
	Name        string
	GormName    string
	JsonName    string
	Package     string // 只有匿名struct不为空
	PackagePath string // 同上
	Type        string
	Default     string
	Tag         string
	Comment     string
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

	//codes, err := ioutil.ReadAll(file)
	//if err != nil {
	//	return nil, err
	//}
	//if len(codes) == 0 {
	//	return nil, fmt.Errorf("GetDstTree - file empty: %s", filePath)
	//}

	f, err := decorator.Parse(file)
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
	PackagePath         string
	Type                string
	HasAnonymousField   bool
	TemplateModelFields []*TemplateModelField
	TemplateModelFuncs  []*TemplateGenFunc
}

type TemplateDstFileImport struct {
	Name string
	Path string
}

type TemplateDstFile struct {
	Imports     []*TemplateDstFileImport
	Package     string
	PackagePath string
	Models      map[string]*TemplateModel
	Funcs       map[string]*TemplateGenFunc
}

func (tm *TemplateDstFile) parseModelFieldFromDst(field *dst.Field) (templateFields []TemplateModelField, err error) {
	if field.Names == nil {
		switch ftp := field.Type.(type) {
		case *dst.Ident:
			if ftp.Obj != nil {
				sonStructTp, ok := ftp.Obj.Decl.(*dst.TypeSpec)
				if ok {
					for _, f := range sonStructTp.Type.(*dst.StructType).Fields.List {
						fields, err := tm.parseModelFieldFromDst(f)
						if err != nil {
							return templateFields, err
						}
						templateFields = append(templateFields, fields...)
					}
				}
			} else {
				templateField := TemplateModelField{
					Package:     tm.Package,
					PackagePath: tm.PackagePath,
					Type:        ftp.Name,
				}
				templateFields = append(templateFields, templateField)
			}
		case *dst.SelectorExpr:
			pkg := ftp.X.(*dst.Ident).Name
			structName := ftp.Sel.Name
			templateField := TemplateModelField{
				Package: pkg,
				Type:    structName,
			}
			// find package path
			for _, fileImport := range tm.Imports {
				if fileImport.Name == pkg {
					templateField.PackagePath = fileImport.Path
				}
			}
			templateFields = append(templateFields, templateField)
		}
	} else if len(field.Names) > 0 {
		templateField := TemplateModelField{}
		templateField.Name = field.Names[0].Name
		if field.Tag != nil {
			templateField.Tag = strings.Trim(field.Tag.Value, "`")
		}
		tags, err := structtag.Parse(templateField.Tag)
		if err != nil {
			return templateFields, err
		}
		for _, tag := range tags.Tags() {
			switch tag.Key {
			case ttypes.MODEL_TAG_TYPE__JSON:
				templateField.JsonName = tag.Name
			case ttypes.MODEL_TAG_TYPE__GORM:
				templateField.GormName = tag.Name
			}
		}
		if len(field.Decs.NodeDecs.Start) > 0 {
			templateField.Comment = field.Decs.NodeDecs.Start[0]
		}
		if len(field.Decs.End) > 0 {
			templateField.Comment = field.Decs.End[0]
		}
		templateField.Type = ParseDstNodeType(field.Type)
		templateFields = []TemplateModelField{templateField}
	}
	return
}

func parseTagToTokens(s string) (rs []string, e error) {
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

func parseTypesVar(v *gotypes.Var, tag string) []TemplateModelField {
	if v.Embedded() {
		t := v.Type()
		str := make([]TemplateModelField, 0)
		if ts, ok := t.Underlying().(*gotypes.Struct); ok {
			for i := 0; i < ts.NumFields(); i++ {
				tmp := parseTypesVar(ts.Field(i), ts.Tag(i))
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

func getCommentByFilterTag(arr []string) string {
	buf := bytes.Buffer{}
	for _, s := range arr {
		if !strings.HasPrefix(s, "// @") && !strings.HasPrefix(s, "//@") {
			buf.WriteString(s)
		}
	}
	return buf.String()
}

func ParseDstCommentFromNode(node dst.NodeDecs, isFilterTag bool) string {
	comment := make([]string, 0)
	if node.Start != nil {
		if !isFilterTag {
			comment = append(comment, node.Start.All()...)
		} else {
			comment = append(comment, getCommentByFilterTag(node.Start.All()))
		}
	}
	if node.End != nil {
		if !isFilterTag {
			comment = append(comment, node.End.All()...)
		} else {
			comment = append(comment, getCommentByFilterTag(node.End.All()))
		}
	}
	return strings.Join(comment, ",")
}

func ParseDstNodeType(tp dst.Expr) string {
	switch tpVal := tp.(type) {
	case *dst.Ident:
		return tpVal.Name
	case *dst.StarExpr:
		return ParseDstNodeType(tpVal.X)
	case *dst.ArrayType:
		return "[]" + ParseDstNodeType(tpVal.Elt)
	case *dst.SliceExpr:
		return "[]" + ParseDstNodeType(tpVal.X)
	case *dst.SelectorExpr:
		name := ""
		if tpVal.X != nil {
			name = tpVal.X.(*dst.Ident).Name
		}
		return name + "." + tpVal.Sel.Name
	case *dst.MapType:
		return fmt.Sprintf("map[%s]%s", ParseDstNodeType(tpVal.Key), ParseDstNodeType(tpVal.Value))
	default:
		// TODO(illidan/2020/11/30): other types
	}
	return ""
}

//func ParseDstTree(file *dst.File) (modelList []TemplateModel, err error) {
//	for _, decl := range file.Decls {
//		gd, ok := decl.(*dst.GenDecl)
//		if !ok {
//			continue
//		}
//		if gd.Tok != token.TYPE {
//			continue
//		}
//		model := TemplateModel{}
//		tf, ok := gd.Specs[0].(*dst.TypeSpec)
//		if !ok {
//			err = fmt.Errorf("can not change to TypeSpec: %#v", gd.Specs)
//			return
//		}
//		// this entity model name
//		model.ModelName = tf.Name.Name
//		model.Package = file.Name.Name
//
//		fieldMap := map[string]TemplateModelField{}
//		st, ok := tf.Type.(*dst.StructType)
//		if !ok {
//			continue
//		}
//
//		for _, field := range st.Fields.List {
//			templateFields, e := parseModelFieldFromDst(field)
//			if e != nil {
//				err = e
//				return
//			}
//			for _, templateField := range templateFields {
//				fieldMap[templateField.Name] = templateField
//				model.TemplateModelFields = append(model.TemplateModelFields, templateField)
//			}
//		}
//
//		// comment def of struct
//		if gd.Decs.NodeDecs.Start != nil {
//			model.ModelComment = strings.Join(gd.Decs.NodeDecs.Start.All(), ",")
//		}
//
//		modelList = append(modelList, model)
//	}
//
//	for _, decl := range file.Decls {
//		if ffv, ok := decl.(*dst.FuncDecl); ok {
//			if len(ffv.Recv.List) == 0 || ffv.Recv.List[0].Type == nil {
//				continue
//			}
//			fmt.Println(ffv)
//		}
//	}
//	return
//}

func ParseFuncDstTree(file *dst.File) (funcList []TemplateGenFunc, err error) {
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

type TemplateDstDir struct {
	Package string
	Path    string
	Imports []*TemplateDstFileImport
	Models  map[string]*TemplateModel   // key: modelName
	Funcs   map[string]*TemplateGenFunc // key: funcName
}

type TemplateDstDirList map[string]*TemplateDstDir

func (tm TemplateDstDirList) ParseDstFile(filePath string) (err error) {
	var dstFile *dst.File
	dstFile, err = DefaultGenTemplate.GetDstTree(filePath)
	if err != nil {
		return
	}
	packagePath := common.GetImportPathFromFile(filePath)

	tmFile := TemplateDstFile{}
	tmFile.Package = dstFile.Name.Name
	tmFile.PackagePath = packagePath
	tmFile.Models = map[string]*TemplateModel{}
	for _, v := range dstFile.Imports {
		pathValue, _ := strconv.Unquote(v.Path.Value)
		name := ""
		if v.Name != nil {
			name = v.Name.Name
		}
		if name == "" {
			name = path.Base(pathValue)
		}
		tmFile.Imports = append(tmFile.Imports, &TemplateDstFileImport{
			Name: name,
			Path: pathValue,
		})
	}
	for _, decl := range dstFile.Decls {
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
		model.Package = tmFile.Package
		model.PackagePath = tmFile.PackagePath

		//fieldMap := map[string]TemplateModelField{}
		switch tfTyp := tf.Type.(type) {
		case *dst.StructType:
			for _, field := range tfTyp.Fields.List {
				templateFields, err := tmFile.parseModelFieldFromDst(field)
				if err != nil {
					return err
				}
				for i := 0; i < len(templateFields); i++ {
					templateField := templateFields[i]
					if templateField.Package != "" && templateField.PackagePath != "" {
						dir, err := common.GetDirFromImport(templateField.PackagePath)
						if err != nil {
							return err
						}
						// 非同一个包才继续解析包，否则会造成循环解析
						if !strings.HasSuffix(dir, packagePath) {
							err = tm.ParseDstDir(dir, 1)
							if err != nil {
								return err
							}
						}
						tmpModel := tm.parseAnonymousFields(templateField.PackagePath, templateField.Type)
						// when package is belong to golang system package, get nil response.
						if tmpModel != nil {
							model.TemplateModelFields = append(model.TemplateModelFields, tmpModel.TemplateModelFields...)
						} else {
							model.TemplateModelFields = append(model.TemplateModelFields, &templateField)
						}
					} else {
						//fieldMap[templateField.Name] = templateField
						model.TemplateModelFields = append(model.TemplateModelFields, &templateField)
					}
				}
			}
			model.Type = model.ModelName
		default:
			model.Type = ParseDstNodeType(tf.Type)
		}
		// comment def of struct
		if gd.Decs.NodeDecs.Start != nil {
			model.ModelComment = strings.Join(gd.Decs.NodeDecs.Start.All(), "\n")
		}
		for _, f := range model.TemplateModelFields {
			if f.Name == "" && f.Type != "" {
				model.HasAnonymousField = true
			}
		}
		tmFile.Models[model.ModelName] = &model
	}
	// parse function
	for _, decl := range dstFile.Decls {
		if ffv, ok := decl.(*dst.FuncDecl); ok {
			fc := TemplateGenFunc{
				Name:           ffv.Name.Name,
				Comment:        "",
				BelongToStruct: nil,
				Package:        tmFile.Package,
				Args:           nil,
				Returns:        nil,
			}
			if ffv.Recv == nil || len(ffv.Recv.List) == 0 {
				continue
			}
			belongModelName := ParseDstNodeType(ffv.Recv.List[0].Type)
			if ffv.Type.Params.List != nil && len(ffv.Type.Params.List) > 0 {
				for _, field := range ffv.Type.Params.List {
					if len(field.Names) > 0 {
						for _, name := range field.Names {
							templateField := TemplateModelField{}
							templateField.Name = name.Name
							templateField.Type = ParseDstNodeType(field.Type)
							fc.Args = append(fc.Args, &templateField)
						}
					} else {
						templateField := TemplateModelField{}
						templateField.Name = field.Type.(*dst.Ident).Name
						templateField.Type = ParseDstNodeType(field.Type)
						fc.Args = append(fc.Args, &templateField)
					}
				}
			}
			if ffv.Type.Results != nil && ffv.Type.Results.List != nil && len(ffv.Type.Results.List) > 0 {
				for _, field := range ffv.Type.Results.List {
					if len(field.Names) > 0 {
						for _, name := range field.Names {
							templateField := TemplateModelField{}
							templateField.Name = name.Name
							templateField.Type = ParseDstNodeType(field.Type)
							fc.Args = append(fc.Args, &templateField)
						}
					} else {
						templateField := TemplateModelField{}
						templateField.Name = field.Type.(*dst.Ident).Name
						templateField.Type = ParseDstNodeType(field.Type)
						fc.Args = append(fc.Args, &templateField)
					}
				}
			}
			if tmp, ok := tmFile.Models[belongModelName]; ok {
				fc.BelongToStruct = tmp
				tmp.TemplateModelFuncs = append(tmp.TemplateModelFuncs, &fc)
				tmFile.Models[belongModelName] = tmp
				tmFile.Funcs[fc.Name] = &fc
			}
		}
	}
	tm.AddDstFile(packagePath, tmFile)
	return
}

func (tm TemplateDstDirList) AddDstFile(dirPath string, dstFile TemplateDstFile) {
	dirPath = strings.TrimPrefix(dirPath, common.GetGoPath()+"/")
	if v, ok := tm[dirPath]; ok {
		v.Imports = append(v.Imports, dstFile.Imports...)
		for k, model := range dstFile.Models {
			v.Models[k] = model
		}
		for k, genFunc := range dstFile.Funcs {
			v.Funcs[k] = genFunc
		}
		tm[dirPath] = v
	} else {
		tm[dirPath] = &TemplateDstDir{
			Package: dstFile.Package,
			Path:    dirPath,
			Imports: dstFile.Imports,
			Models:  dstFile.Models,
			Funcs:   dstFile.Funcs,
		}
	}
	return
}

func (tm TemplateDstDirList) parseAnonymousFields(pkgPath, name string) *TemplateModel {
	if strings.Index(pkgPath, "/") == -1 {
		// not to parse golang system package
		return nil
	} else {
		var err error
		pkgPath, err = common.GetTotalImportPathFromImport(pkgPath)
		if err != nil {
			return nil
		}
	}
	if dstDir, ok := tm[pkgPath]; ok {
		if model, ok := dstDir.Models[name]; ok {
			if model.HasAnonymousField {
				fields := []*TemplateModelField{}
				for _, mField := range model.TemplateModelFields {
					if mField.Name == "" && mField.Type != "" {
						tmpModel := tm.parseAnonymousFields(mField.Package, mField.Type)
						fields = append(fields, tmpModel.TemplateModelFields...)
					} else {
						fields = append(fields, mField)
					}
				}
				model.TemplateModelFields = fields
			}
			// check anonymous field
			model.HasAnonymousField = false
			for _, field := range model.TemplateModelFields {
				if field.Name == "" {
					model.HasAnonymousField = true
				}
			}
			return model
		}
	}
	return nil
}

func (tm TemplateDstDirList) ParseAnonymousModel(oldModel *TemplateModel) {
	if oldModel.HasAnonymousField {
		fields := []*TemplateModelField{}
		for _, field := range oldModel.TemplateModelFields {
			if field.Name == "" && field.Type != "" {
				tmpModel := tm.parseAnonymousFields(field.PackagePath, field.Type)
				if tmpModel != nil {
					fields = append(fields, tmpModel.TemplateModelFields...)
				} else {
					fields = append(fields, field)
				}
			} else {
				fields = append(fields, field)
			}
		}
		oldModel.TemplateModelFields = fields
		// check anonymous field
		oldModel.HasAnonymousField = false
		for _, field := range oldModel.TemplateModelFields {
			if field.Name == "" {
				oldModel.HasAnonymousField = true
			}
		}
	}
}

func (tm TemplateDstDirList) ParseDstDir(dirPath string, depth int) (err error) {
	err = tm.parseDstDir(dirPath, depth, 0)
	if err != nil {
		return
	}
	return
}

func (tm TemplateDstDirList) parseDstDir(dirPath string, depth, currentDepth int) (err error) {
	if !common.IsExists(dirPath) {
		err = errors.New(dirPath + " is not exist when ParseModelDir")
		return
	}
	imDir := common.GetImportPathFromDir(dirPath)
	if _, ok := tm[imDir]; ok {
		// 已经解析过的不再解析
		return
	}
	currentDepth++
	var files []os.FileInfo
	files, err = ioutil.ReadDir(dirPath)
	if err != nil {
		return
	}
	// load go file first
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".go" {
			err := tm.ParseDstFile(filepath.Join(dirPath, file.Name()))
			if err != nil {
				return err
			}
		}
	}
	// check anonymous field of model
	for _, dstDir := range tm {
		for _, model := range dstDir.Models {
			if model.HasAnonymousField {
				tm.ParseAnonymousModel(model)
			}
		}
	}
	if depth > currentDepth || depth == -1 {
		for _, file := range files {
			if file.IsDir() {
				err = tm.parseDstDir(filepath.Join(dirPath, file.Name()), depth, currentDepth)
				if err != nil {
					return
				}
			}
		}
	}
	return
}

func ParseDstDir(dirPath string, depth int) (TemplateDstDirList, error) {
	tm := TemplateDstDirList{}
	err := tm.ParseDstDir(dirPath, depth)
	return tm, err
}

//func parseModelDir(dirPath string, depth int, currentDepth int) (modelDirs []TemplateDstDir, err error) {
//	currentDepth++
//	if !common.IsExists(dirPath) {
//		err = errors.New(dirPath + " is not exist when ParseModelDir")
//		return
//	}
//	var files []os.FileInfo
//	modelDirs = []TemplateDstDir{}
//	files, err = ioutil.ReadDir(dirPath)
//	if err != nil {
//		return
//	}
//	// load go file first
//	currentModelDir := TemplateDstDir{}
//	for _, file := range files {
//		if !file.IsDir() && filepath.Ext(file.Name()) == ".go" {
//			err = currentModelDir.ParseModelDstTree(filepath.Join(dirPath, file.Name()))
//			if err != nil {
//				return
//			}
//		}
//	}
//	// check anonymous field of model
//	for _, dstFiles := range modelDir {
//		for _, dstFile := range dstFiles {
//			for _, model := range dstFile.Models {
//				fields := make([]TemplateModelField, 0)
//				for _, field := range model.TemplateModelFields {
//					if field.Name == "" && field.Type != "" {
//						tmpFields := modelDir.ParseModelAnonymousFields(field.Package, field.Type)
//						if len(tmpFields) > 0 {
//							fields = append(fields, tmpFields...)
//						} else {
//							fields = append(fields, field)
//						}
//					} else {
//						fields = append(fields, field)
//					}
//				}
//				model.TemplateModelFields = fields
//			}
//		}
//	}
//	if depth > currentDepth || depth == -1 {
//		for _, file := range files {
//			if file.IsDir() {
//				tmDirs, e := parseModelDir(filepath.Join(dirPath, file.Name()), depth, currentDepth)
//				if e != nil {
//					err = e
//					return
//				}
//				for _, dir := range tmDirs {
//					modelDirs = append(modelDirs, dir)
//				}
//			}
//		}
//	}
//	return
//}

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
