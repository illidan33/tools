package dao

import (
	"fmt"
	"regexp"

	"github.com/dave/dst"
	"github.com/illidan33/tools/gen"
	"github.com/illidan33/tools/gen/method"
	"github.com/illidan33/tools/gen/util/types"
)

const templateMethodTxt = `package {{$.PackageName}}

import (
	"github.com/jinzhu/gorm"
	"github.com/m2c/kiplestar"
	{{range $pkg := $.PackageList}}
	"{{html $pkg}}"
	{{end}}
)

//go:generate tools kiple interface -i {{$.InterfaceName}} -m {{$.ModelName}}
type {{$.InterfaceName}} interface {
	{{range $funcName := $.CmdKipleDaoFuncNames}}
	{{html $funcName}}
	{{end}}
}

func {{$.InterfaceName}}Instance() {{$.InterfaceName}} {
	return &{{$.ModelName}}{
		db: kiplestar.GetKipleServerInstance().DB("replace with your db name").DB(),
	}
}

type {{$.ModelName}} struct {
	db *gorm.DB
}

{{range $func := .TemplateDataMethodFuncs}}
{{html $func}}
{{end}}

`

var templateMethodMap = map[string]string{
	"FetchBy%s": `func (d *{{$.ModelName}}) {{$.FuncName}}({{$.ConditionStr}}) (*{{$.EntityPackageName}}{{$.EntityName}}, error) {
			entt := {{$.EntityPackageName}}{{$.EntityName}}{}
			if err := d.db.Model(entt).Where("{{$.WhereStr}}", {{$.ConditionFieldStr}}).First(&entt).Error; err != nil{
				return nil, err
			}
			return &entt, nil
		}`,
	"UpdateBy%sWithStruct": `func (d *{{$.ModelName}}) {{$.FuncName}}(entt {{$.EntityPackageName}}{{$.EntityName}}) (error) {
			if err := d.db.Model(entt).Where("{{$.WhereStr}}", {{$.FieldStr}}).Updates(entt).Error; err != nil{
				return err
			}
			return nil
		}`,
	"UpdateBy%sWithMap": `func (d *{{$.ModelName}}) {{$.FuncName}}(args map[string]interface{}) (error) {
			entt := {{$.EntityPackageName}}{{$.EntityName}}{}
			if err := d.db.Model(entt).Where("{{$.WhereStr}}", {{$.FieldStr}}).Updates(args).Error; err != nil{
				return err
			}
			return nil
		}`,
	"BatchFetchBy%s": `func (d *{{$.ModelName}}) {{$.FuncName}}({{$.ConditionStr}})(dList []{{$.EntityPackageName}}{{$.EntityName}}, err error) {
			err = d.db.Where("{{$.WhereStr}}", {{$.ConditionFieldStr}}).Find(&dList).Error
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
	"Create": `func (d *{{$.ModelName}}) Create(entt {{$.EntityPackageName}}{{$.EntityName}}) (error) {
			if err := d.db.Create(entt).Error; err != nil{
				return err
			}
			return nil
		}`,
	"Delete": `func (d *{{$.ModelName}}) Delete(entt {{$.EntityPackageName}}{{$.EntityName}}) (error) {
			if err := d.db.Delete(entt).Error; err != nil {
				return err
			}
			return nil
		}`,
	"FetchList": `func (d *{{$.ModelName}}) FetchList(size int32, offset int32, args map[string]interface{})(dList []{{$.EntityPackageName}}{{$.EntityName}}, count int32, err error) {
			if size == -1 {
				err = d.db.Where(args).Offset(offset).Find(&dList).Count(&count).Error
			} else {
				err = d.db.Where(args).Offset(offset).Limit(size).Find(&dList).Count(&count).Error
			}			
			return 
		}`,
}

type KipleTemplateDao struct {
	InterfaceName        string
	EntityName           string
	EntityPackageName    string
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
	Fields []gen.TemplateModelField
}

func (tgm *KipleTemplateDao) joinConditionFields(fields []gen.TemplateModelField) string {
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

// public
func (tgm *KipleTemplateDao) ParseKipleIndexToMethod() error {
	var err error
	err = tgm.ParseIndexToMethod()
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
	tm.EntityPackageName = file.Name.Name

	err := tm.ParseDstTree(file)
	if err != nil {
		return err
	}
	tm.EntityName = tm.ModelName

	return nil
}
