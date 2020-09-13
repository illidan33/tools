package modle

import (
	"myprojects/tools/gen"
	"sort"
)

const templateModleTxt = `package {{ .PackageName }}

import (
	{{ range $value := .PackageList }} "{{ $value }}" {{end}}
)
{{$filename := var .StructName}}
{{$fields := sort $.ModleStructFields}}

//go:generate tools gen method -f ./{{snake $filename}}.go
{{html .StructComment}}
type {{ .StructName }} struct{
	{{ range $value := $fields }}
	{{ $value.Name }} {{ $value.Type }} {{ html $value.Tag }} // {{html $value.Comment }}
	{{end}}
}`

func registeTemplateFunc(tms *gen.TemplateGenModle) {
	tms.Init()
	tms.Registe(map[string]interface{}{
		"sortField": func(fields map[string]gen.TemplateGenStructField) []gen.TemplateGenStructField {
			s := make([]gen.TemplateGenStructField, len(fields))
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
	})
}
