package method

import (
	"bytes"
	"fmt"
	"github.com/dave/dst"
	"go/token"
	"myprojects/tools/gen"
	"myprojects/tools/gen/types"
	"path/filepath"
	"strings"
	"text/template"
)

const templateMethodTxt = `// Code generated by "gen method"; DO NOT EDIT

package {{ .PackageName }}

import (
	"github.com/jinzhu/gorm"
)
{{ $StructNames := printf "%ss" .StructName }}

type {{ $StructNames }} []{{ .StructName }}

{{ range $func := .TemplateGenMethodFuncs }}
{{ $func }}
{{end}}

`

var templateMethodFetchByIndexTxt = `func (%s *%s) FetchBy%s(db *gorm.DB) error {
	if err := db.Where("%s", %s).First(%s).Error; err != nil{
		return err
	}
	return nil
}`

var templateMethodUpdateByIndexWithStructTxt = `func (%s *%s) UpdateBy%sWithStruct(db *gorm.DB) error {
	if err := db.Model(%s).Where("%s", %s).Updates(%s).Error; err != nil{
		return err
	}
	return nil
}`

var templateMethodUpdateByIndexWithMapTxt = `func (%s *%s) UpdateBy%sWithMap(db *gorm.DB, args map[string]interface{}) error {
	if err := db.Model(%s).Where("%s", %s).Updates(args).Error; err != nil{
		return err
	}
	return nil
}`

var templateMethodCreateTxt = `func (%s *%s) Create(db *gorm.DB) error {
	if err := db.Create(%s).Error; err != nil{
		return err
	}
	return nil
}`

var templateMethodDeleteTxt = `func (%s *%s) Delete(db *gorm.DB) error {
	if err := db.Delete(%s).Error; err != nil {
		return err
	}
	return nil
}`

var templateMethodFetchListTxt = `func (%s *%s) FetchList(db *gorm.DB, args map[string]interface{})(%s %s, err error) {
	err = db.Where(args).Find(&%s).Error
	return 
}`

var templateMethodBatchFetchByIndexTxt = `func (%s *%s) BatchFetchBy%s(db *gorm.DB)(%s %s, err error) {
	err = db.Where("%s", %s).Find(&%s).Error
	return 
}`

var templateMethodBatchFetchByIndexListTxt = `func (%s *%s) BatchFetchBy%sList(db *gorm.DB, %ss []%s)(%s %s, err error) {
	err = db.Where("%s", %ss).Find(&%s).Error
	return 
}`

type TemplateGenMethod struct {
	gen.TemplateGenModle
	CommentIndexs          []TemplateGenMethodCommentIndex
	TemplateGenMethodFuncs []string
}

type TemplateGenMethodCommentIndex struct {
	Name   string
	Type   types.IndexType
	Fields []gen.TemplateGenStructField
}

func registeTemplateFunc(tms *TemplateGenMethod) {
	tms.Init()
	tms.Registe(map[string]interface{}{
	})
}

func (tgm *TemplateGenMethod) joinFields(structName string, fields []gen.TemplateGenStructField) string {
	rs := ""
	for i, arg := range fields {
		if i == 0 {
			rs = fmt.Sprintf("%s.%s", structName, arg.Name)
		} else {
			rs = fmt.Sprintf("%s, %s.%s", rs, structName, arg.Name)
		}
	}
	return rs
}

func (tgm *TemplateGenMethod) joinWhere(fields []gen.TemplateGenStructField) string {
	rs := ""
	for i, arg := range fields {
		name := gen.ToLowerSnakeCase(arg.Name)
		if i == 0 {
			rs = fmt.Sprintf("%s=?", name)
		} else {
			rs = fmt.Sprintf("%s AND %s=?", rs, name)
		}
	}
	return rs
}

func (tgm *TemplateGenMethod) joinConditions(fields []gen.TemplateGenStructField) string {
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

func (tgm *TemplateGenMethod) parseLineToTokens(s string) (rs []string, e error) {
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

func (tgm *TemplateGenMethod) parseDstTree(file *dst.File) error {
	tgm.PackageName = file.Name.Name
	if tgm.PackageName == "" {
		return fmt.Errorf("packageName empty: %v", file.Name)
	}
	for _, i := range file.Decls {
		gd, ok := i.(*dst.GenDecl)
		if !ok {
			continue
		}
		if gd.Tok != token.TYPE {
			continue
		}
		tf := gd.Specs[0].(*dst.TypeSpec)
		if tgm.StructName == "" {
			tgm.StructName = tf.Name.Name
		}

		st := tf.Type.(*dst.StructType)
		for k, field := range st.Fields.List {
			templateField := gen.TemplateGenStructField{
				Name: field.Names[0].Name,
				Tag:  field.Tag.Value,
				Sort: k,
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
			tgm.ModleStructFields[templateField.Name] = templateField
		}

		// comment def of struct
		if gd.Decs.NodeDecs.Start != nil {
			decs := gd.Decs.NodeDecs.Start
			for _, dec := range decs {
				if strings.Contains(dec, "@def") {
					arr := strings.Split(dec, " ")
					if arr[0] == "//" && arr[1] == "@def" {
						tgmci := TemplateGenMethodCommentIndex{}
						names := strings.Split(arr[2], ":")
						if len(names) > 1 {
							tgmci.Name = names[1]
						}
						switch names[0] {
						case types.INDEXTYPE__PRIMARY.KeyLowerString():
							tgmci.Type = types.INDEXTYPE__PRIMARY
						case types.INDEXTYPE__UNIQUE_INDEX.KeyLowerString():
							tgmci.Type = types.INDEXTYPE__UNIQUE_INDEX
						case types.INDEXTYPE__INDEX.KeyLowerString():
							tgmci.Type = types.INDEXTYPE__INDEX
						default:
						}
						tgmci.Fields = make([]gen.TemplateGenStructField, 0)
						for i := 3; i < len(arr); i++ {
							if f, ok := tgm.ModleStructFields[arr[i]]; !ok {
								return fmt.Errorf("index field of comment def is not struct field: %s", arr[i])
							} else {
								tgmci.Fields = append(tgmci.Fields, f)
							}
						}
						tgm.CommentIndexs = append(tgm.CommentIndexs, tgmci)
					}
				}
			}
		}
	}

	// TODO(illidan/2020/9/13):
	//for _, index := range tgm.CommentIndexs {
	//	tgmfFetch := gen.TemplateGenModleFunc{
	//		Name:           "",
	//		BelongToStruct: tgm.StructName,
	//		Args:           index.Fields,
	//		Returns:        []interface{}{"error"},
	//	}
	//	name := ""
	//	for i2, field := range index.Fields {
	//		if i2 == 0 {
	//			name = fmt.Sprintf("%s", field.Name)
	//		} else {
	//			name = fmt.Sprintf("%s_%s", name, field.Name)
	//		}
	//	}
	//	tgmfFetch.Name = "FetchBy" + gen.ToUpperCamelCase(name)
	//	tgm.ModleFuncs = append(tgm.ModleFuncs, tgmfFetch)
	//
	//	if index.Type == types.INDEXTYPE__INDEX {
	//		tgmfNew := gen.TemplateGenModleFunc{
	//			Name:           "",
	//			BelongToStruct: tgm.StructName,
	//			Args:           index.Fields,
	//			Returns:        []interface{}{"error"},
	//		}
	//		tgmfNew.Name = "BatchFetchBy" + gen.ToUpperCamelCase(name)
	//		tgm.ModleFuncs = append(tgm.ModleFuncs, tgmfNew)
	//	}
	//}

	return nil
}

func (tgm *TemplateGenMethod) parseToMethods() error {
	structName := gen.ToLowerCamelCase(tgm.StructName)
	structNames := gen.ToLowerCamelCase(tgm.StructName) + "s"
	StructNames := tgm.StructName + "s"
	for _, index := range tgm.CommentIndexs {
		baseFuncName := tgm.GenFuncName(index.Fields)
		joinFields := tgm.joinFields(structName, index.Fields)
		joinWhere := tgm.joinWhere(index.Fields)
		//joinConditions := tgm.joinConditions(index.Fields)
		tgmf1 := fmt.Sprintf(templateMethodFetchByIndexTxt, structName, tgm.StructName, baseFuncName, joinWhere, joinFields, structName)
		tgmf2 := fmt.Sprintf(templateMethodUpdateByIndexWithStructTxt, structName, tgm.StructName, baseFuncName, structName, joinWhere, joinFields, structName)
		tgmf3 := fmt.Sprintf(templateMethodUpdateByIndexWithMapTxt, structName, tgm.StructName, baseFuncName, structName, joinWhere, joinFields)
		tgmf4 := fmt.Sprintf(templateMethodBatchFetchByIndexTxt, structName, tgm.StructName, baseFuncName, structNames, StructNames, joinWhere, joinFields, structNames)
		tgm.TemplateGenMethodFuncs = append(tgm.TemplateGenMethodFuncs, tgmf1, tgmf2, tgmf3, tgmf4)
		if (index.Type == types.INDEXTYPE__PRIMARY || index.Type == types.INDEXTYPE__UNIQUE_INDEX) && len(index.Fields) == 1 {
			uniqField := index.Fields[0]
			tgmf5 := fmt.Sprintf(templateMethodBatchFetchByIndexListTxt, structName, tgm.StructName, baseFuncName, gen.ToLowerCamelCase(uniqField.Name), uniqField.Type, structNames, StructNames, joinWhere, gen.ToLowerCamelCase(uniqField.Name), structNames)
			tgm.TemplateGenMethodFuncs = append(tgm.TemplateGenMethodFuncs, tgmf5)
		}
	}
	tgmf6 := fmt.Sprintf(templateMethodCreateTxt, structName, tgm.StructName, structName)
	tgmf7 := fmt.Sprintf(templateMethodDeleteTxt, structName, tgm.StructName, structName)
	tgmf8 := fmt.Sprintf(templateMethodFetchListTxt, structName, tgm.StructName, structNames, StructNames, structNames)
	tgm.TemplateGenMethodFuncs = append(tgm.TemplateGenMethodFuncs, tgmf6, tgmf7, tgmf8)
	return nil
}

func (tgm *TemplateGenMethod) Parse(flags CmdGenMethodFlags) error {
	basePath, _ := filepath.Split(flags.CmdGenModleFilePath)
	s := strings.LastIndex(basePath, "/")
	if s == -1 {
		return fmt.Errorf("basepath error: %s", basePath)
	}

	tgm.StructName = flags.CmdGenModleName

	dstTree, err := tgm.GetDstTree(flags.CmdGenModleFilePath)
	if err != nil {
		return err
	}
	if err = tgm.parseDstTree(dstTree); err != nil {
		return err
	}
	if err = tgm.parseToMethods(); err != nil {
		return err
	}

	return nil
}

func (tgm *TemplateGenMethod) ParseTemplate(templateTxt string) (templateSource *bytes.Buffer, e error) {
	templateSource = bytes.NewBuffer([]byte(""))

	tp := template.New("gen_method")
	tp.Funcs(tgm.TemplateFuncs)
	if tp, e = tp.Parse(templateTxt); e != nil {
		return
	}
	e = tp.Execute(templateSource, *tgm)
	return
}
