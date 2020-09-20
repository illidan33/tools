package model

import "myprojects/tools/gen"

const templateModelTxt = `package {{ .PackageName }}

import (
	{{ range $value := .PackageList }} "{{ $value }}" {{end}}
)
{{$filename := var .ModelName}}

//go:generate tools gen method
{{html .ModelComment}}
type {{ type .ModelName }} struct{
	{{ range $value := $.TemplateModelFields }}
	{{ $value.Name }} {{ $value.Type }} {{ html $value.Tag }} // {{html $value.Comment }}
	{{end}}
}`


type TemplateDataModel struct {
	gen.GenTemplate
	gen.TemplatePackage
	gen.TemplateModel
}