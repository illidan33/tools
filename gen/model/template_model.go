package model

import "github.com/illidan33/tools/gen"

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
