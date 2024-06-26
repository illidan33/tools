package dao_create

import (
	"bytes"
	"fmt"
	"regexp"
	ttempl "text/template"
	"tools/common"
	"tools/template"

	"github.com/dave/dst"
	"tools/gen/method"
	"tools/template/types"
)

var templateDaoTxt = `package {{$.PackageName}}

import (
	"github.com/jinzhu/gorm"
	"github.com/m2c/kiplestar"
	{{range $pkg := $.PackageList}}
	"{{html $pkg}}"
	{{end}}
)

//go:generate tools kiple methodsync -i {{$.InterfaceName}} -m {{$.ModelName}}
type {{$.InterfaceName}} interface {
	{{range $funcName := $.CmdKipleDaoFuncNames}}
	{{html $funcName}}
	{{end}}
}

func {{$.InterfaceName}}Instance() {{$.InterfaceName}} {
	return &{{$.ModelName}}{
		db: kiplestar.GetKipleServerInstance().DB("{{$.DbName}}").DB(),
	}
}

type {{$.ModelName}} struct {
	db *gorm.DB
}

`
var templateDaoGenTxt = `package {{$.PackageName}}

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"reflect"
	{{range $pkg := $.PackageList}}
	"{{html $pkg}}"
	{{end}}
)

{{range $func := .TemplateDataMethodFuncs}}
{{html $func}}
{{end}}

`

var templateMethodMap = map[string]string{
	"FetchBy%s": `func (d *{{$.ModelName}}) {{$.FuncName}}({{$.ConditionStr}}) ({{var $.EntityName}} {{$.EntityPackageName}}{{$.EntityName}},err error) {
		err = d.db.Model({{var $.EntityName}}).Where("{{$.WhereStr}}", {{$.ConditionFieldStr}}).First(&{{var $.EntityName}}).Error
		return 
	}`,
	"BatchFetchBy%s": `func (d *{{$.ModelName}}) {{$.FuncName}}({{$.ConditionStr}}) ({{var $.EntityName}}List []{{$.EntityPackageName}}{{$.EntityName}},err error) {
		err = d.db.Model(&{{$.EntityPackageName}}{{$.EntityName}}{}).Where("{{$.WhereStr}}", {{$.ConditionFieldStr}}).Find(&{{var $.EntityName}}List).Error
		return 
	}`,
	"UpdateBy%sWithStruct": `func (d *{{$.ModelName}}) {{$.FuncName}}({{var $.EntityName}} *{{$.EntityPackageName}}{{$.EntityName}}) (err error) {
		err = d.db.Model({{var $.EntityName}}).Where("{{$.WhereStr}}", {{$.FieldStr}}).Updates({{var $.EntityName}}).Error
		return 
	}`,
	"UpdateBy%sWithMap": `func (d *{{$.ModelName}}) {{$.FuncName}}({{$.ConditionStr}}, args map[string]interface{}) (err error) {
		entt := {{$.EntityPackageName}}{{$.EntityName}}{}
		err = d.db.Model(entt).Where("{{$.WhereStr}}", {{$.ConditionFieldStr}}).Updates(args).Error
		return 
	}`,
}

var templateMethodFiedUniqMap = map[string]string{
	"BatchFetchBy%sList": `func (d *{{$.ModelName}}) {{$.FuncName}}({{var $.UniqFieldName}}List []{{$.UniqFieldType}})(dList []{{$.EntityPackageName}}{{$.EntityName}}, err error) {
		err = d.db.Where("{{snake $.UniqFieldName}} in (?)", {{var $.UniqFieldName}}List).Find(&dList).Error
		return 
	}`,
}

var templateMethodUniqMap = map[string]string{
	"Create": `func (d *{{$.ModelName}}) Create({{var $.EntityName}} *{{$.EntityPackageName}}{{$.EntityName}}) (err error) {
		err = d.db.Create({{var $.EntityName}}).Error
		return 
	}`,
	"Delete": `func (d *{{$.ModelName}}) Delete({{var $.EntityName}} *{{$.EntityPackageName}}{{$.EntityName}}) (err error) {
		err = d.db.Delete({{var $.EntityName}}).Error
		return 
	}`,
	"BatchFetchBySql": `func (d *{{$.ModelName}}) BatchFetchBySql(sql string, args ...interface{}) ({{var $.EntityName}}List []{{$.EntityPackageName}}{{$.EntityName}}, err error) {
		{{var $.EntityName}}List, err = d.FetchList(-1, 0, nil, &sql, args...)
		return
	}`,
	"FetchOneBySql": `func (d *{{$.ModelName}}) FetchOneBySql(sql string,args ...interface{}) ({{var $.EntityName}} {{$.EntityPackageName}}{{$.EntityName}}, err error) {
		err = d.db.Model(&{{$.EntityPackageName}}{{$.EntityName}}{}).Where(sql, args...).First(&{{var $.EntityName}}).Error
		return
	}`,
	"BatchUpdateBySql": `func (d *{{$.ModelName}}) BatchUpdateBySql(sql string, args []interface{}, updates map[string]interface{}) (err error) {
		err = d.db.Model(&{{$.EntityPackageName}}{{$.EntityName}}{}).Where(sql, args...).Updates(updates).Error
		return
	}`,
	"FetchList": `func (d *{{$.ModelName}}) FetchList(size int32, offset int32, count *int32, sql *string, args ...interface{}) ({{var $.EntityName}}List []{{$.EntityPackageName}}{{$.EntityName}}, err error) {
		m := {{$.EntityPackageName}}{{$.EntityName}}{}
		var db *gorm.DB
		if sql != nil {
			if size == -1 {
				db = d.db.Model(&m).Where(*sql, args...).Offset(offset).Find(&{{var $.EntityName}}List)
			} else {
				db = d.db.Model(&m).Where(*sql, args...).Offset(offset).Limit(size).Find(&{{var $.EntityName}}List)
			}
		} else {
			if len(args) > 1 {
				err = fmt.Errorf("args should be one,when sql is nil.")
				return
			}
			t := reflect.TypeOf(args[0])
			switch t.Kind() {
			case reflect.Struct:
				if size == -1 {
					db = d.db.Model(&m).Where(args[0]).Offset(offset).Find(&{{var $.EntityName}}List)
				} else {
					db = d.db.Model(&m).Where(args[0]).Offset(offset).Limit(size).Find(&{{var $.EntityName}}List)
				}
			case reflect.Map:
				if size == -1 {
					db = d.db.Model(&m).Where(args[0]).Offset(offset).Find(&{{var $.EntityName}}List)
				} else {
					db = d.db.Model(&m).Where(args[0]).Offset(offset).Limit(size).Find(&{{var $.EntityName}}List)
				}
			default:
				err = fmt.Errorf("args should be map or struct when sql is nil, but give %s.", t.Kind().String())
				return
			}
		}
		err = db.Error
		if err != nil {
			if gorm.IsRecordNotFoundError(err) {
				count = nil
				err = nil
			}
		} else if count != nil {
			db.Count(count)
		}
		return
	}`,
}

type KipleTemplateDao struct {
	InterfaceName        string
	EntityName           string
	EntityPackageName    string
	DbName               string
	CmdKipleDaoFuncNames []string

	method.TemplateDataMethod
}

type CmdKipleDaoFunc struct {
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

type CmdKipleDaoIndex struct {
	Name   string
	Type   types.IndexType
	Fields []template.TemplateModelField
}

func (tgm *KipleTemplateDao) joinConditionFields(fields []*template.TemplateModelField) string {
	rs := ""
	for i, arg := range fields {
		if i == 0 {
			rs = fmt.Sprintf("%s", common.ToLowerCamelCase(arg.Name))
		} else {
			rs = fmt.Sprintf("%s, %s", rs, common.ToLowerCamelCase(arg.Name))
		}
	}
	return rs
}

// public
func (tgm *KipleTemplateDao) ParseKipleIndexToMethod() error {
	var err error
	err = tgm.ParseIndexToMethod(templateMethodMap, templateMethodFiedUniqMap, templateMethodUniqMap)
	if err != nil {
		return err
	}

	if tgm.CmdKipleDaoFuncNames == nil {
		tgm.CmdKipleDaoFuncNames = make([]string, 0)
	}
	reg := regexp.MustCompile("func \\([^\\(^\\)]*\\) (.*) {")
	for _, tpFunc := range tgm.TemplateDataMethodFuncs {
		s := reg.FindAllStringSubmatch(tpFunc, -1)
		if len(s) > 0 && len(s[0]) > 1 {
			tgm.CmdKipleDaoFuncNames = append(tgm.CmdKipleDaoFuncNames, s[0][1])
		}
	}

	return nil
}

func (tm *KipleTemplateDao) ParseKipleDstTree(file *dst.File) error {
	sourceModelName := tm.ModelName
	tm.ModelName = ""

	err := tm.ParseDstTree(file)
	if err != nil {
		return err
	}

	tm.EntityPackageName = file.Name.Name
	tm.EntityName = tm.ModelName
	tm.ModelName = sourceModelName

	return nil
}

func (tgm *KipleTemplateDao) parseMethodFuncsToTemplate(tp *ttempl.Template, td CmdKipleDaoFunc, templateTxt string, templateName string) (err error) {
	td.FuncName = fmt.Sprintf(templateName, td.FuncName)
	templateSource := &bytes.Buffer{}
	tp, err = tp.Parse(templateTxt)
	if err != nil {
		err = fmt.Errorf("parse [%s] template error: %s\n", templateName, err.Error())
		return
	}
	err = tp.Execute(templateSource, td)
	if err != nil {
		err = fmt.Errorf("execute [%s] tmplate data error: %s\n", templateName, err.Error())
		return
	}
	tgm.TemplateDataMethodFuncs = append(tgm.TemplateDataMethodFuncs, templateSource.String())
	return nil
}

func (tgm *KipleTemplateDao) ParseIndexToMethod(templateMethodMap, templateMethodFieldUniqMap, templateMethodUniqMap map[string]string) error {
	var err error
	td := CmdKipleDaoFunc{
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
	tp := ttempl.New("kiple method")
	tp.Funcs(tgm.TemplateMapFuncs)
	for _, index := range tgm.TemplateDataMethodIndexs {
		// TODO(illidan/2020/9/28): foreign index not include
		if index.Type == types.INDEX_TYPE__FOREIGN_INDEX {
			continue
		}
		td.FuncName = tgm.GenFuncName(index.Fields)
		td.WhereStr = tgm.JoinWhere(index.Fields)
		td.FieldStr = tgm.JoinFields(common.ToLowerCamelCase(tgm.EntityName), index.Fields)
		td.ConditionStr = tgm.JoinConditions(index.Fields)
		td.ConditionFieldStr = tgm.joinConditionFields(index.Fields)

		for k, tpMethod := range templateMethodMap {
			if err = tgm.parseMethodFuncsToTemplate(tp, td, tpMethod, k); err != nil {
				return err
			}
		}

		if (index.Type == types.INDEX_TYPE__PRIMARY || index.Type == types.INDEX_TYPE__UNIQUE_INDEX) && len(index.Fields) == 1 {
			td.UniqFieldName = index.Fields[0].Name
			td.UniqFieldType = index.Fields[0].Type
			for k, tpMethod := range templateMethodFieldUniqMap {
				if err = tgm.parseMethodFuncsToTemplate(tp, td, tpMethod, k); err != nil {
					return err
				}
			}
		}
	}

	for k, tpMethod := range templateMethodUniqMap {
		if err = tgm.parseMethodFuncsToTemplate(tp, td, tpMethod, k); err != nil {
			return err
		}
	}
	return nil
}
