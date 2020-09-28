package model

import "myprojects/tools/gen"

const templateModelTxt = `package {{ .PackageName }}

import (
	{{ range $value := .PackageList }} "{{ $value }}" {{end}}
)
{{$filename := var .ModelName}}

//go:generate tools gen method --name={{type .ModelName}}
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
