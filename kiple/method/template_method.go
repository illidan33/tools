package method

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/illidan33/tools/common"
	"github.com/illidan33/tools/gen"
	"github.com/illidan33/tools/gen/util/types"
	"go/importer"
	"go/token"
	gotypes "go/types"
	"os"
	"regexp"
	"strings"
	"text/template"
)

const templateMethodTxt = `package {{$.PackageName}}

import (
	"github.com/m2c/kiplestar/kipledb"
)

//go:generate tools kiple method --name UserDao --entity={{$.EntityPath}}
type {{$.InterfaceName}} interface {
	{{range $funcName := $.TemplateDataMethodFuncNames}}
	{{html $funcName}}
	{{end}}
}

func {{$.InterfaceName}}Instance() {{$.InterfaceName}} {
	return &{{$.ModelName}}{
		db: nil, // need to replace nil
	}
}

type {{$.ModelName}} struct {
	db *kipledb.KipleDB
}

{{range $func := .TemplateDataMethodFuncs}}
{{html $func}}
{{end}}

`

var templateMethodMap = map[string]string{
	"FetchBy%s": `func (d *{{$.ModelName}}) {{$.FuncName}}() (*{{$.EntityPackageName}}{{$.EntityName}}, error) {
			entt := {{$.EntityPackageName}}{{$.EntityName}}{}
			if err := d.db.DB().Model(entt).Where("{{$.WhereStr}}", {{$.FieldStr}}).First(&entt).Error; err != nil{
				return nil, err
			}
			return &entt, nil
		}`,
	"UpdateBy%sWithStruct": `func (d *{{$.ModelName}}) {{$.FuncName}}(entt {{$.EntityPackageName}}{{$.EntityName}}) (error) {
			if err := d.db.DB().Model(entt).Where("{{$.WhereStr}}", {{$.FieldStr}}).Updates(entt).Error; err != nil{
				return err
			}
			return nil
		}`,
	"UpdateBy%sWithMap": `func (d *{{$.ModelName}}) {{$.FuncName}}(args map[string]interface{}) (error) {
			entt := {{$.EntityPackageName}}{{$.EntityName}}{}
			if err := d.db.DB().Model(entt).Where("{{$.WhereStr}}", {{$.FieldStr}}).Updates(args).Error; err != nil{
				return err
			}
			return nil
		}`,
	"BatchFetchBy%s": `func (d *{{$.ModelName}}) {{$.FuncName}}({{$.ConditionStr}})(dList []{{$.EntityPackageName}}{{$.EntityName}}, err error) {
			err = d.db.DB().Where("{{$.WhereStr}}", {{$.ConditionFieldStr}}).Find(&dList).Error
			return 
		}`,
}

var templateMethodFiedUniqMap = map[string]string{
	"BatchFetchBy%sList": `func (d *{{$.ModelName}}) {{$.FuncName}}({{var $.UniqFieldName}}List []{{$.UniqFieldType}})(dList []{{$.EntityPackageName}}{{$.EntityName}}, err error) {
			err = d.db.DB().Where("{{snake $.UniqFieldName}} in (?)", {{var $.UniqFieldName}}List).Find(&dList).Error
			return 
		}`,
}

var templateMethodUniqMap = map[string]string{
	"Create": `func (d *{{$.ModelName}}) Create(entt {{$.EntityPackageName}}{{$.EntityName}}) (error) {
			if err := d.db.DB().Create(entt).Error; err != nil{
				return err
			}
			return nil
		}`,
	"Delete": `func (d *{{$.ModelName}}) Delete(entt {{$.EntityPackageName}}{{$.EntityName}}) (error) {
			if err := d.db.DB().Delete(entt).Error; err != nil {
				return err
			}
			return nil
		}`,
	"FetchList": `func (d *{{$.ModelName}}) FetchList(size int32, offset int32, args map[string]interface{})(dList []{{$.EntityPackageName}}{{$.EntityName}}, count int32, err error) {
			err = d.db.DB().Where(args).Offset(offset).Limit(size).Find(&dList).Count(&count).Error
			return 
		}`,
}

type TemplateDataMethod struct {
	InterfaceName     string
	EntityPath        string
	EntityName        string
	EntityPackageName string
	gen.GenTemplate
	gen.TemplatePackage
	gen.TemplateModel
	TemplateDataMethodFuncNames []string
	TemplateDataMethodFuncs     []string
	TemplateDataMethodIndexs    []TemplateDataMethodIndex
}

type TemplateDataMethodFunc struct {
	ModelName         string
	ModelNames        string
	EntityName        string
	EntityPackageName string
	FuncName          string
	WhereStr          string
	FieldStr          string
	ConditionStr      string
	ConditionFieldStr string
	UniqFieldName     string
	UniqFieldType     string
}

type TemplateDataMethodIndex struct {
	Name   string
	Type   types.IndexType
	Fields []gen.TemplateModelField
}

func (gt *TemplateDataMethod) genFuncName(fields []gen.TemplateModelField) string {
	str := ""
	for i, f := range fields {
		if i == 0 {
			str = f.Name
		} else {
			str += "_And" + f.Name
		}
	}
	return common.ToUpperCamelCase(str)
}

func (tgm *TemplateDataMethod) joinFields(modelName string, fields []gen.TemplateModelField) string {
	rs := ""
	for i, arg := range fields {
		if i == 0 {
			rs = fmt.Sprintf("%s.%s", modelName, arg.Name)
		} else {
			rs = fmt.Sprintf("%s, %s.%s", rs, modelName, arg.Name)
		}
	}
	return rs
}

func (tgm *TemplateDataMethod) joinConditionFields(fields []gen.TemplateModelField) string {
	rs := ""
	for i, arg := range fields {
		if i == 0 {
			rs = fmt.Sprintf("%s", arg.Name)
		} else {
			rs = fmt.Sprintf("%s, %s", rs, arg.Name)
		}
	}
	return rs
}

func (tgm *TemplateDataMethod) joinWhere(fields []gen.TemplateModelField) string {
	rs := ""
	for i, arg := range fields {
		name := common.ToLowerSnakeCase(arg.Name)
		if i == 0 {
			rs = fmt.Sprintf("%s=?", name)
		} else {
			rs = fmt.Sprintf("%s AND %s=?", rs, name)
		}
	}
	return rs
}

func (tgm *TemplateDataMethod) joinConditions(fields []gen.TemplateModelField) string {
	rs := ""
	for i, arg := range fields {
		if i == 0 {
			rs = fmt.Sprintf("%s %s", arg.Name, arg.Type)
		} else {
			rs = fmt.Sprintf("%s, %s %s", rs, arg.Name, arg.Type)
		}
	}
	return rs
}

func (tgm *TemplateDataMethod) parseMethodFuncsToTemplate(tp *template.Template, reg *regexp.Regexp, td TemplateDataMethodFunc, templateTxt string) (err error) {
	templateSource := &bytes.Buffer{}
	tp, err = tp.Parse(templateTxt)
	if err != nil {
		return
	}
	err = tp.Execute(templateSource, td)
	if err != nil {
		return
	}
	templateString := templateSource.String()
	tgm.TemplateDataMethodFuncs = append(tgm.TemplateDataMethodFuncs, templateString)
	if tgm.TemplateDataMethodFuncNames == nil {
		tgm.TemplateDataMethodFuncNames = make([]string, 0)
	}
	s := reg.FindAllStringSubmatch(templateString, -1)
	if len(s) > 0 && len(s[0]) > 1 {
		tgm.TemplateDataMethodFuncNames = append(tgm.TemplateDataMethodFuncNames, s[0][1])
	}

	return nil
}

func (tgm *TemplateDataMethod) ParseIndexToMethod() error {
	var err error
	td := TemplateDataMethodFunc{
		ModelName:         tgm.ModelName,
		ModelNames:        tgm.ModelName + "List",
		EntityName:        tgm.EntityName,
		EntityPackageName: tgm.EntityPackageName,
	}
	if tgm.PackageName == td.EntityPackageName {
		td.EntityPackageName = ""
	} else {
		td.EntityPackageName += "."
	}
	if len(tgm.TemplateMapFuncs) == 0 {
		tgm.InitTemplateFuncs()
	}
	tp := template.New("kiple method")
	tp.Funcs(tgm.TemplateMapFuncs)
	reg := regexp.MustCompile("func \\([^\\(^\\)]*\\) (.*) {")
	for _, index := range tgm.TemplateDataMethodIndexs {
		// TODO(illidan/2020/9/28): foreign index not include
		if index.Type == types.INDEX_TYPE__FOREIGN_INDEX {
			continue
		}
		baseFuncName := tgm.genFuncName(index.Fields)
		joinFields := tgm.joinFields(common.ToLowerCamelCase("entt"), index.Fields)
		joinWhere := tgm.joinWhere(index.Fields)
		joinConditions := tgm.joinConditions(index.Fields)
		joinFieldConditions := tgm.joinConditionFields(index.Fields)
		td.FuncName = baseFuncName
		td.WhereStr = joinWhere
		td.FieldStr = joinFields
		td.ConditionStr = joinConditions
		td.ConditionFieldStr = joinFieldConditions
		for name, templ := range templateMethodMap {
			td.FuncName = fmt.Sprintf(name, baseFuncName)
			if err = tgm.parseMethodFuncsToTemplate(tp, reg, td, templ); err != nil {
				return err
			}
		}
		if (index.Type == types.INDEX_TYPE__PRIMARY || index.Type == types.INDEX_TYPE__UNIQUE_INDEX) && len(index.Fields) == 1 {
			td.UniqFieldName = index.Fields[0].Name
			td.UniqFieldType = index.Fields[0].Type
			for name, templ := range templateMethodFiedUniqMap {
				td.FuncName = fmt.Sprintf(name, baseFuncName)
				if err = tgm.parseMethodFuncsToTemplate(tp, reg, td, templ); err != nil {
					return err
				}
			}
		}
	}
	for name, templ := range templateMethodUniqMap {
		td.FuncName = fmt.Sprintf(name, td.FuncName)
		if err = tgm.parseMethodFuncsToTemplate(tp, reg, td, templ); err != nil {
			return err
		}
	}
	return nil
}

func (tm *TemplateDataMethod) parseTagToTokens(s string) (rs []string, e error) {
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

func (tm *TemplateDataMethod) parseDecsToIndex(decs dst.Decorations, fieldMap *map[string]gen.TemplateModelField) error {
	for _, dec := range decs {
		if strings.Contains(dec, "@def") {
			arr := strings.Split(dec, " ")
			if arr[0] == "//" && arr[1] == "@def" {
				tgmci := TemplateDataMethodIndex{}
				names := strings.Split(arr[2], ":")
				if len(names) > 1 {
					tgmci.Name = names[1]
				}
				switch names[0] {
				case types.INDEX_TYPE__PRIMARY.KeyLowerString():
					tgmci.Type = types.INDEX_TYPE__PRIMARY
				case types.INDEX_TYPE__UNIQUE_INDEX.KeyLowerString():
					tgmci.Type = types.INDEX_TYPE__UNIQUE_INDEX
				case types.INDEX_TYPE__INDEX.KeyLowerString():
					tgmci.Type = types.INDEX_TYPE__INDEX
				case types.INDEX_TYPE__FOREIGN_INDEX.KeyLowerString():
					tgmci.Type = types.INDEX_TYPE__FOREIGN_INDEX
				default:
				}
				tgmci.Fields = make([]gen.TemplateModelField, 0)
				for i := 3; i < len(arr); i++ {
					if f, ok := (*fieldMap)[arr[i]]; !ok {
						return fmt.Errorf("index field of comment def is not struct field: %s", arr[i])
					} else {
						tgmci.Fields = append(tgmci.Fields, f)
					}
				}
				tm.TemplateDataMethodIndexs = append(tm.TemplateDataMethodIndexs, tgmci)
			}
		}
	}
	return nil
}

func (tm *TemplateDataMethod) ParseDstTree(file *dst.File) error {
	tm.EntityPackageName = file.Name.Name
	for _, i := range file.Decls {
		gd, ok := i.(*dst.GenDecl)
		if !ok {
			continue
		}
		if gd.Tok != token.TYPE {
			continue
		}
		tf, ok := gd.Specs[0].(*dst.TypeSpec)
		if !ok {
			return fmt.Errorf("can not change to TypeSpec: %#v", gd.Specs)
		}
		if tm.EntityName == "" {
			tm.EntityName = tf.Name.Name
		}

		fieldMap := map[string]gen.TemplateModelField{}
		if len(tm.TemplateModelFields) == 0 {
			st, ok := tf.Type.(*dst.StructType)
			if !ok {
				return fmt.Errorf("can not change to StructType: %#v", tf.Type)
			}

			for _, field := range st.Fields.List {
				templateField := gen.TemplateModelField{}

				if len(field.Names) > 0 {
					templateField.Name = field.Names[0].Name
					templateField.Tag = field.Tag.Value
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
				tm.TemplateModelFields = append(tm.TemplateModelFields, templateField)
			}
		} else {
			for _, field := range tm.TemplateModelFields {
				fieldMap[field.Name] = field
			}
		}

		// comment def of struct
		if gd.Decs.NodeDecs.Start != nil {
			decs := gd.Decs.NodeDecs.Start
			err := tm.parseDecsToIndex(decs, &fieldMap)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (tm *TemplateDataMethod) FindInterfaceAndFillMethods(astfile *dst.File, dstFilePath string) (bool, error) {
	var userdao *dst.TypeSpec
	flag := false
	userdaoFuncMap := map[string]bool{}
	for _, decl := range astfile.Decls {
		if declv, ok := decl.(*dst.GenDecl); ok {
			if declv.Tok == token.TYPE {
				if len(declv.Specs) == 0 {
					return false, errors.New("GenDecl has no Specs")
				}
				if typespec, ok := declv.Specs[0].(*dst.TypeSpec); ok && typespec.Name.Name == tm.ModelName {
					userdao = typespec
					flag = true
					if userdaointerface, ok := userdao.Type.(*dst.InterfaceType); ok {
						for _, field := range userdaointerface.Methods.List {
							userdaoFuncMap[field.Names[0].Name] = true
						}
					}
					break
				}
			}
		}
	}
	if userdao == nil {
		return flag, nil
	}
	for _, decl := range astfile.Decls {
		if ffv, ok := decl.(*dst.FuncDecl); ok && ffv.Recv != nil && ffv.Recv.List[0].Names[0].Name == tm.InterfaceName {
			ffnew := dst.Field{}
			ffnew.Names = []*dst.Ident{ffv.Name}
			ffnew.Type = ffv.Type
			if userdaointerface, ok := userdao.Type.(*dst.InterfaceType); ok {
				if _, ok := userdaoFuncMap[ffv.Name.Name]; !ok {
					userdaointerface.Methods.List = append(userdaointerface.Methods.List, &ffnew)
					userdaoFuncMap[ffv.Name.Name] = true
				}
			}
		}
	}

	var file *os.File
	file, err := os.OpenFile(dstFilePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return flag, err
	}
	defer file.Close()
	err = decorator.Fprint(file, astfile)
	return flag, nil
}

func (tm *TemplateDataMethod) Parse(filePath string, isDebug bool) error {
	gofile := os.Getenv("GOFILE")
	importerFilePath := ""
	if isDebug {
		importerFilePath = strings.TrimPrefix(filePath, os.Getenv("GOPATH")+"/src/")
		importerFilePath = strings.TrimSuffix(importerFilePath, gofile)
	} else {
		importerFilePath = gofile
	}

	dstTree, err := tm.GetDstTree(filePath)
	if err != nil {
		return err
	}
	if err = tm.ParseDstTree(dstTree); err != nil {
		return err
	}
	if err = tm.ParseIndexToMethod(); err != nil {
		return err
	}
	return nil
}

func (tm *TemplateDataMethod) ImportFile(filePath string) error {
	pkg, err := importer.For("source", nil).Import(filePath)
	if err != nil {
		return err
	}

	elem := pkg.Scope().Lookup(tm.ModelName)
	strArr := make([]gen.TemplateModelField, 0)
	if named, ok := elem.Type().(*gotypes.Named); ok {
		if ts, ok := named.Underlying().(*gotypes.Struct); ok {
			for i := 0; i < ts.NumFields(); i++ {
				tmp := tm.parseTypesVar(ts.Field(i), ts.Tag(i))
				strArr = append(strArr, tmp...)
			}
		}
	}
	tm.TemplateModelFields = strArr
	return nil
}

func (tm *TemplateDataMethod) parseTypesVar(v *gotypes.Var, tag string) []gen.TemplateModelField {
	if v.Embedded() {
		t := v.Type()
		str := make([]gen.TemplateModelField, 0)
		if ts, ok := t.Underlying().(*gotypes.Struct); ok {
			for i := 0; i < ts.NumFields(); i++ {
				tmp := tm.parseTypesVar(ts.Field(i), ts.Tag(i))
				str = append(str, tmp...)
			}
		}
		return str
	} else {
		templateField := gen.TemplateModelField{
			Name:    v.Name(),
			Type:    v.Type().String(),
			Tag:     tag,
			Comment: "",
		}
		return []gen.TemplateModelField{templateField}
	}
}
