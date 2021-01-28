package model

import (
	"bytes"
	"strings"
	"tools/template"
)

const templateModelTxt = `package {{ .PackageName }}

import (
	{{ range $value := .PackageList }} 
	"{{ $value }}"
	{{end}}
)
{{$filename := var .ModelName}}

{{html .ModelComment}}
type {{ type .ModelName }} struct{
	{{ range $value := $.TemplateModelFields }}
	{{ $value.Name }} {{ $value.Type }} {{ html $value.Tag }} {{if hasComment $value}}// {{html $value.Comment }} {{end}}
	{{end}}
}`

type TemplateDataModel struct {
	template.GenTemplate
	template.TemplatePackage
	template.TemplateModel
}

func (tp *TemplateDataModel) Parse() (*bytes.Buffer, error) {
	for _, field := range tp.TemplateModelFields {
		if strings.Contains(field.Type, "time") {
			tp.AddPackage("time", "time")
		} else if strings.Contains(field.Type, "BitBool") {
			tp.AddPackage("bool", "github.com/m2c/kiplestar/kipledb/types")
		}
	}

	codeData, err := template.DefaultGenTemplate.ParseTemplate(templateModelTxt, tp.ModelName, tp, map[string]interface{}{
		"hasComment": func(field template.TemplateModelField) bool {
			if field.Comment != "" {
				return true
			}
			return false
		},
	})
	if err != nil {
		return nil, err
	}
	return codeData, nil
}
