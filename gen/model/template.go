package model

import (
	"bytes"
	"tools/gen"
	"strings"
)

const templateModelTxt = `package {{ .PackageName }}

import (
	{{ range $value := .PackageList }} "{{ $value }}" {{end}}
)
{{$filename := var .ModelName}}

{{html .ModelComment}}
type {{ type .ModelName }} struct{
	{{ range $value := $.TemplateModelFields }}
	{{ $value.Name }} {{ $value.Type }} {{ html $value.Tag }} {{if hasComment $value}}// {{html $value.Comment }} {{end}}
	{{end}}
}`

type TemplateDataModel struct {
	gen.GenTemplate
	gen.TemplatePackage
	gen.TemplateModel
}

func (tp *TemplateDataModel) Parse() (*bytes.Buffer, error) {
	for _, field := range tp.TemplateModelFields {
		if strings.Contains(field.Type, "time") {
			tp.AddPackage("time", "time")
		}
	}

	codeData, err := gen.DefaultGenTemplate.ParseTemplate(templateModelTxt, tp.ModelName, tp, map[string]interface{}{
		"hasComment": func(field gen.TemplateModelField) bool {
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
