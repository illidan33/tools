package modle

import (
	"fmt"
	"html/template"
	"myprojects/tools/gen"
	"myprojects/tools/gen/types"
	"sort"
)

const templateModleTxt = `package {{ .PackageName }}

import (
	{{ range $value := .PackageList }} "{{ $value }}" {{end}}
)
{{ $fields := sortField $.GormTable.Fields }}
{{$filename := var .StructName}}

//go:generate gormtools gen method -f ./{{snake $filename}}.go
{{ range $value := $.GormTable.Indexs }}{{genIndex $value}}{{end}}
type {{ .StructName }} struct{
	{{ range $value := $fields }}
	{{ type $value.Name }} {{ $value.Type }} {{ genTag $value }} // {{html $value.Comment }}
	{{end}}
}`

type TemplateModleStruct struct {
	gen.TemplateGenStruct
	GormTable gen.GormTable
}

func registeTemplateFunc(tms *TemplateModleStruct) {
	tms.Init()
	tms.Registe(map[string]interface{}{
		"genTag": func(field gen.GormField) template.HTML {
			nullStr := "NOT NULL"
			if field.IsNull {
				nullStr = "NULL"
			}
			//index := ""
			//if field.IsKeyField {
			//	switch field.KeyType {
			//	case types.INDEXTYPE__PRIMARY:
			//		index = "primary_key;"
			//	case types.INDEXTYPE__UNIQUE_INDEX:
			//		index = fmt.Sprintf("unique_key:%s;", field.KeyName)
			//	case types.INDEXTYPE__INDEX:
			//		index = fmt.Sprintf("index:%s;", field.KeyName)
			//	default:
			//	}
			//}
			return template.HTML(fmt.Sprintf("`gorm:\"column:%s;type:%s;%s;default:%s\" json:\"%s\"`", field.Name, field.SqlType, nullStr, field.Default, field.Name))
		},
		"sortField": func(fields map[string]gen.GormField) []gen.GormField {
			s := make([]gen.GormField, len(fields))
			i := 0
			for _, field := range fields {
				s[i] = field
				i++
			}
			sort.Slice(s, func(i, j int) bool {
				if s[i].Sort < s[j].Sort {
					return true
				}
				return false
			})
			return s
		},
		"genIndex": func(index gen.GormIndex) string {
			str := ""
			if len(index.Fields) > 0 {
				for i, field := range index.Fields {
					if i == 0 {
						if index.Type == types.INDEXTYPE__PRIMARY {
							str = fmt.Sprintf("\n// @def %s %s", index.Type.KeyLowerString(), gen.ToUpperCamelCase(field.Name))
						} else {
							str = fmt.Sprintf("\n// @def %s:%s %s", index.Type.KeyLowerString(), index.Name, gen.ToUpperCamelCase(field.Name))
						}
					} else {
						str = fmt.Sprintf("%s %s", str, gen.ToUpperCamelCase(field.Name))
					}
				}
			}
			return fmt.Sprintf("%s", str)
		},
	})
}
