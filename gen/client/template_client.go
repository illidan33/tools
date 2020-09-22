package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"myprojects/tools/common"
	"myprojects/tools/gen"
	"myprojects/tools/gen/util/types"
	"myprojects/tools/httptool"
	"os"
	"strings"
)

const templateClientTxt = `package {{ .PackageName }}

import (
	"myprojects/tools/httptool"
	{{ range $value := .PackageList }} "{{ $value }}" {{end}}
)
{{$ModelName := .ClientModel.ModelName}}
{{$modelName := var $ModelName}}

// {{html .ClientModel.ModelComment}}
type {{$ModelName}} struct{
	httptool.Client
}

{{range $func := $.ClientFuncs}}
{{$isModel := isModel $func.RequestModel}}
{{if $isModel}}
{{if $func.ResponseModel.IsArray}}
// {{$func.Comment}}
func (c {{$ModelName}}){{$func.Name}}(req {{$func.RequestModel.ModelName}}, headers ...map[string]string)(resp *[]{{$func.ResponseModel.Model.ModelName}}, err error) {
	resp = &[]{{$func.ResponseModel.Model.ModelName}}{}
	err = c.Request("{{$func.Method}}", "{{$func.Path}}", req, resp, headers...)
	return
}
{{else}}
// {{$func.Comment}}
func (c {{$ModelName}}){{$func.Name}}(req {{$func.RequestModel.ModelName}}, headers ...map[string]string)(resp *{{$func.ResponseModel.Model.ModelName}}, err error) {
	resp = &{{$func.ResponseModel.Model.ModelName}}{}
	err = c.Request("{{$func.Method}}", "{{$func.Path}}", req, resp, headers...)
	return
}
{{end}}
{{else}}
{{if $func.ResponseModel.IsArray}}
// {{$func.Comment}}
func (c {{$ModelName}}){{$func.Name}}(headers ...map[string]string)(resp *[]{{$func.ResponseModel.Model.ModelName}}, err error) {
	resp = &[]{{$func.ResponseModel.Model.ModelName}}{}
	err = c.Request("{{$func.Method}}", "{{$func.Path}}", nil, resp, headers...)
	return
}
{{else}}
// {{$func.Comment}}
func (c {{$ModelName}}){{$func.Name}}(headers ...map[string]string)(resp *{{$func.ResponseModel.Model.ModelName}}, err error) {
	resp = &{{$func.ResponseModel.Model.ModelName}}{}
	err = c.Request("{{$func.Method}}", "{{$func.Path}}", nil, resp, headers...)
	return
}
{{end}}
{{end}}
{{end}}
`
const templateModelTxt = `package {{ .PackageName }}
{{range $model := $.ParamModels}}
// {{$model.ModelComment}}
type {{$model.ModelName}} struct {
	{{ range $value := $model.TemplateModelFields }}
	{{ $value.Name }} {{ $value.Type }} {{ html $value.Tag }} // {{html $value.Comment }}
	{{end}}
}
{{end}}
`

type TemplateGenClient struct {
	gen.GenTemplate
	gen.TemplatePackage
	// main model
	ClientModel gen.TemplateModel
	ClientFuncs []TemplateGenClientFunc
	// response and request models
	ParamModels []gen.TemplateModel
}

type TemplateGenClientFunc struct {
	Name           string
	Comment        string
	BelongToStruct gen.TemplateModel
	Path           string
	Method         string
	RequestModel   gen.TemplateModel
	ResponseModel  TemplateGenClientFuncReturn
}

type TemplateGenClientFuncReturn struct {
	Model   gen.TemplateModel
	IsArray bool
}

type GenClientSwagger struct {
	Info        GenClientSwaggerInfo                `json:"info"`
	Paths       map[string]map[string]GenClientFunc `json:"paths"`
	Definitions map[string]GenClientDefinition      `json:"definitions"`
}

type GenClientSwaggerInfo struct {
	Description string `json:"description"`
	Title       string `json:"title"`
	Version     string `json:"version"`
}

type GenClientDefinition struct {
	Type       string                       `json:"type"`
	Properties map[string]map[string]string `json:"properties"`
}
type GenClientDefinitionPropertie struct {
	Description string `json:"description"`
	Type        string `json:"type"`
}

type GenClientFunc struct {
	Name       string                           `json:"-"`
	Method     string                           `json:"-"`
	Path       string                           `json:"-"`
	Consumes   []string                         `json:"consumes"`
	Produces   []string                         `json:"produces"`
	Tags       []string                         `json:"tags"`
	Summary    string                           `json:"summary"`
	Parameters []GenClientFuncParam             `json:"parameters"`
	Responses  map[string]GenClientFuncResponse `json:"responses"`
}

type GenClientFuncParam struct {
	Description string            `json:"description"`
	Name        string            `json:"name"`
	In          string            `json:"in"`
	Required    bool              `json:"required"`
	Schema      map[string]string `json:"schema"`
}

type GenClientFuncResponse struct {
	Description string                     `json:"description"`
	Schema      map[string]json.RawMessage `json:"schema"`
}

func (tgc *TemplateGenClient) parseRefModelName(s string) string {
	name := strings.Trim(s, "\"")
	i := strings.LastIndex(name, "/")
	if i != -1 {
		ii := strings.LastIndex(name[i+1:], ".")
		if ii != -1 {
			name = name[i+ii+2:]
		}
	}
	return name
}

func (tgc *TemplateGenClient) parseDefinitions(defs map[string]GenClientDefinition) error {
	paramModelsMap := map[string]gen.TemplateModel{}
	for name, def := range defs {
		s := strings.Index(name, ".")
		if s != -1 {
			name = name[s+1:]
		}
		f := gen.TemplateModel{
			ModelName: common.ToUpperCamelCase(name),
		}
		if def.Type == "object" {
			for proName, pro := range def.Properties {
				if pro["type"] == "object" {
					name := tgc.parseRefModelName(pro[types.SWAGGER_TYPE__REF])
					f.TemplateModelFields = append(f.TemplateModelFields, gen.TemplateModelField{
						Name:    common.ToUpperCamelCase(proName),
						Type:    name,
						Tag:     fmt.Sprintf("`json:\"%s\"`", proName),
						Comment: pro["Description"],
					})
				} else {
					fName := common.ToUpperCamelCase(proName)
					f.TemplateModelFields = append(f.TemplateModelFields, gen.TemplateModelField{
						Name:    common.ToUpperCamelCase(fName),
						Type:    pro["type"],
						Tag:     fmt.Sprintf("`json:\"%s\"`", proName),
						Comment: pro["Description"],
					})
				}
			}
		}
		tgc.ParamModels = append(tgc.ParamModels, f)
		paramModelsMap[f.ModelName] = f
	}
	return nil
}

func (tgc *TemplateGenClient) parseFuncs(funcs []GenClientFunc) error {
	paramModelsMap := map[string]gen.TemplateModel{}
	for _, model := range tgc.ParamModels {
		paramModelsMap[model.ModelName] = model
	}
	for _, clientFunc := range funcs {
		name := clientFunc.Name
		s := strings.LastIndex(name, "/")
		if s != -1 {
			name = name[s+1:]
		}
		f := TemplateGenClientFunc{
			Name:           common.ToUpperCamelCase(name),
			Comment:        clientFunc.Summary,
			BelongToStruct: tgc.ClientModel,
			Path:           clientFunc.Path,
			Method:         clientFunc.Method,
			RequestModel:   gen.TemplateModel{},
			ResponseModel:  TemplateGenClientFuncReturn{},
		}
		for _, parameter := range clientFunc.Parameters {
			if len(parameter.Schema) > 0 {
				if ref, ok := parameter.Schema[types.SWAGGER_TYPE__REF]; ok {
					name := tgc.parseRefModelName(ref)
					if model, ok := paramModelsMap[name]; ok {
						f.RequestModel = model
						break
					} else {
						return errors.New("param schema model parse error: " + name)
					}
				}
			}
		}
		for _, response := range clientFunc.Responses {
			for k, v := range response.Schema {
				if k == types.SWAGGER_TYPE__REF {
					name := tgc.parseRefModelName(string(v))
					if model, ok := paramModelsMap[name]; ok {
						f.ResponseModel.Model = model
					} else {
						return errors.New("respnse model parse error: " + name)
					}
				} else if k == types.SWAGGER_TYPE__ITEMS {
					data := map[string]string{}
					err := json.Unmarshal([]byte(v), &data)
					if err != nil {
						return err
					}
					name := tgc.parseRefModelName(data[types.SWAGGER_TYPE__REF])
					if model, ok := paramModelsMap[name]; ok {
						f.ResponseModel.Model = model
					} else {
						return errors.New("respnse model parse error: " + name)
					}
					f.ResponseModel.IsArray = true
				}
			}
		}
		tgc.ClientFuncs = append(tgc.ClientFuncs, f)
	}
	return nil
}

func (tgc *TemplateGenClient) Parse(url string) error {
	req := httptool.InitHttpRequest(url)
	body, err := req.Get()
	if err != nil {
		panic(err)
	}

	// for test
	//file, err := os.Open("/data/golang/go/src/myprojects/tools/example/fpx.json")
	//if err != nil {
	//	return err
	//}
	//body, _ := ioutil.ReadAll(file)

	data := GenClientSwagger{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return err
	}

	err = tgc.parseDefinitions(data.Definitions)
	if err != nil {
		return err
	}

	clientFuncs := make([]GenClientFunc, 0)
	for uri, path := range data.Paths {
		for method, funcItem := range path {
			funcItem.Name = uri
			funcItem.Method = strings.ToUpper(method)
			funcItem.Path = uri
			clientFuncs = append(clientFuncs, funcItem)
		}
	}
	err = tgc.parseFuncs(clientFuncs)
	if err != nil {
		return err
	}
	tgc.ClientModel.ModelComment = data.Info.Description

	return nil
}

func (tgc *TemplateGenClient) ParseTemplateAndFormatToFile(exeFilePath string) error {
	folderPath := fmt.Sprintf("%s/%s", exeFilePath, tgc.PackageName)
	if !common.IsExists(folderPath) {
		err := os.MkdirAll(folderPath, os.ModePerm)
		if err != nil {
			return err
		}
	}

	tgc.RegisteTemplateFunc(map[string]interface{}{
		"isModel": func(m gen.TemplateModel) bool {
			if m.ModelName == "" {
				return false
			}
			return true
		},
	})

	bf, err := tgc.ParseTemplate(templateModelTxt, "templateModelTxt", tgc)
	if err != nil {
		return err
	}

	err = tgc.FormatCodeToFile(fmt.Sprintf("%s/types_generate.go", folderPath), bf)
	if err != nil {
		return err
	}

	bf, err = tgc.ParseTemplate(templateClientTxt, "templateClientTxt", tgc)
	if err != nil {
		return err
	}

	err = tgc.FormatCodeToFile(fmt.Sprintf("%s/client_generate.go", folderPath), bf)
	if err != nil {
		return err
	}
	return nil
}
