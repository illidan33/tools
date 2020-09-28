package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"myprojects/tools/common"
	"myprojects/tools/gen"
	"myprojects/tools/gen/util/types"
	"myprojects/tools/httptool"
	"os"
	"regexp"
	"strings"
	"time"
)

const templateClientTxt = `// Code generated by "tools"; DO NOT EDIT

package {{ .PackageName }}

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
{{$hasReq := isModel $func.RequestModel}}
{{$hasResp := isModel $func.ResponseModel}}
{{$isStruct := isStruct $func.ResponseModel}}
{{if $hasResp}}
// {{$func.Comment}}
func (c {{$ModelName}}){{$func.Name}}({{if $hasReq}}req {{$func.RequestModel.ModelName}},{{end}} headers ...map[string]string)(resp {{$func.ResponseModel.ModelName}}, err error) {
	resp = {{$func.ResponseModel.ModelName}}{}
	var rs []byte
	rs, err = c.Request("{{$func.Method}}", "{{$func.Path}}", {{if $hasReq}}req{{else}}nil{{end}}, headers...)
	if err != nil {
		return
	}
	{{if $isStruct}}
	err = c.ParseToResult(rs, &resp)
	{{else}}
	resp = rs
	{{end}}
	return
}
{{else}}
// {{$func.Comment}}
func (c {{$ModelName}}){{$func.Name}}({{if $hasReq}}req {{$func.RequestModel.ModelName}},{{end}} headers ...map[string]string)(err error) {	
	_, err = c.Request("{{$func.Method}}", "{{$func.Path}}", {{if $hasReq}}req{{else}}nil{{end}}, headers...)
	return
}
{{end}}
{{end}}
`
const templateModelTxt = `// Code generated by "tools"; DO NOT EDIT
package {{ .PackageName }}

{{range $model := $.ParamModels}}
{{$isStruct := isStruct $model}}
{{if $isStruct}}
// {{$model.ModelComment}}
type {{$model.ModelName}} struct {
	{{ range $value := $model.TemplateModelFields }}
	{{ $value.Name }} {{ $value.Type }} {{ html $value.Tag }} // {{html $value.Comment }}
	{{end}}
}
{{else}}
// {{$model.ModelComment}}
type {{$model.ModelName}} {{$model.Type}}
{{end}}
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
	ResponseModel  gen.TemplateModel
}

type TemplateGenClientFuncReturn struct {
	Model   gen.TemplateModel
	IsArray bool
}

type GenClientSwagger struct {
	Info        GenClientSwaggerInfo                `json:"info"`
	Host        string                              `json:"host"`
	BasePath    string                              `json:"basePath"`
	Paths       map[string]map[string]GenClientFunc `json:"paths"`
	Definitions map[string]GenClientDefinition      `json:"definitions"`
}

type GenClientSwaggerInfo struct {
	Description string `json:"description"`
	Title       string `json:"title"`
	Version     string `json:"version"`
}

type GenClientDefinition struct {
	Type       string                                  `json:"type"`
	Required   []string                                `json:"required"`
	Properties map[string]GenClientDefinitionPropertie `json:"properties"`
}

type GenClientDefinitionPropertie struct {
	Description string                                  `json:"description"`
	Type        string                                  `json:"type"`
	Items       GenClientSwaggerItem                    `json:"items"`
	Properties  map[string]GenClientDefinitionPropertie `json:"properties"`
}

type GenClientFunc struct {
	Name        string                           `json:"-"`
	Method      string                           `json:"-"`
	Path        string                           `json:"-"`
	Description string                           `json:"description"`
	Consumes    []string                         `json:"consumes"`
	Produces    []string                         `json:"produces"`
	Tags        []string                         `json:"tags"`
	Summary     string                           `json:"summary"`
	OperationId string                           `json:"operationId"`
	Parameters  []GenClientFuncParam             `json:"parameters"`
	Responses   map[string]GenClientFuncResponse `json:"responses"`
}

type GenClientFuncParam struct {
	Description string              `json:"description"`
	Name        string              `json:"name"`
	In          string              `json:"in"`
	Required    bool                `json:"required"`
	Type        string              `json:"type"`
	Schema      GenClientFuncSchema `json:"schema"`
}

type GenClientFuncResponse struct {
	Description string              `json:"description"`
	Schema      GenClientFuncSchema `json:"schema"`
}

type GenClientSwaggerItem struct {
	GenClientSwaggerRef
	Type string `json:"type"`
}

type GenClientFuncSchema struct {
	GenClientSwaggerItem
	Items GenClientSwaggerRef `json:"items"`
}

type GenClientSwaggerRef struct {
	Ref string `json:"$ref"`
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

func (tgc *TemplateGenClient) transformGoType(s string) string {
	tp := ""
	switch s {
	case types.SWAGGER_PROPERTY_TYPE__STRING:
		tp = "string"
	case types.SWAGGER_PROPERTY_TYPE__INTEGER:
		tp = "int64"
	case types.SWAGGER_PROPERTY_TYPE__BOOLEAN:
		tp = "bool"
	case types.SWAGGER_PROPERTY_TYPE__NUMBER:
		tp = "int64"
	default:
	}
	return tp
}

func (tgc *TemplateGenClient) parseDefinitionProperty(modleName string, proName string, pro GenClientDefinitionPropertie) (*gen.TemplateModelField, error) {
	tp := ""
	switch pro.Type {
	case types.SWAGGER_PROPERTY_TYPE__ARRAY:
		if pro.Items.Ref != "" {
			name := pro.Items.Ref
			s := strings.LastIndex(name, "/")
			if s != -1 {
				ss := strings.LastIndex(name[s+1:], ".")
				if ss != -1 {
					name = name[s+ss+2:]
				} else {
					name = name[s+1:]
				}
			}
			tp = "[]" + name
		} else if pro.Items.Type != "" {
			tp = "[]" + tgc.transformGoType(pro.Items.Type)
		}

	case types.SWAGGER_PROPERTY_TYPE__OBJECT:
		model := gen.TemplateModel{
			ModelName:    fmt.Sprintf("%s%s", modleName, common.ToUpperCamelCase(proName)),
			ModelComment: pro.Description,
			Type:         "struct",
		}
		if len(pro.Properties) > 0 {
			for proProName, propertie := range pro.Properties {
				field, err := tgc.parseDefinitionProperty(modleName, proProName, propertie)
				if err != nil {
					return nil, err
				}
				model.TemplateModelFields = append(model.TemplateModelFields, *field)
			}
		}

		tgc.ParamModels = append(tgc.ParamModels, model)
		tp = model.ModelName
	default:
		tp = tgc.transformGoType(pro.Type)
	}
	return &gen.TemplateModelField{
		Name:    common.ToUpperCamelCase(proName),
		Type:    tp,
		Tag:     fmt.Sprintf("`json:\"%s\"`", proName),
		Comment: pro.Description,
	}, nil
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
			Type:      "struct",
		}
		switch def.Type {
		case types.SWAGGER_PROPERTY_TYPE__OBJECT:
			for proName, pro := range def.Properties {
				field, err := tgc.parseDefinitionProperty(f.ModelName, proName, pro)
				if err != nil {
					return err
				}
				f.TemplateModelFields = append(f.TemplateModelFields, *field)
			}
		default:
			return errors.New("unknow swagger property type: " + def.Type)
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
		reg := regexp.MustCompile(`[^a-zA-Z-_0-9/]`)
		funcName := reg.ReplaceAllString(clientFunc.Name, "")

		if funcName == "" {
			return errors.New("Definition path empty")
		}
		funcName = strings.Trim(funcName, "/")
		name := common.ToUpperCamelCase(strings.ReplaceAll(funcName, "/", "_"))
		f := TemplateGenClientFunc{
			Name:           name,
			Comment:        clientFunc.Summary,
			BelongToStruct: tgc.ClientModel,
			Path:           clientFunc.Path,
			Method:         clientFunc.Method,
			RequestModel:   gen.TemplateModel{},
			ResponseModel:  gen.TemplateModel{},
		}
		f.RequestModel = gen.TemplateModel{
			ModelName: fmt.Sprintf("%sRequest", f.Name),
			Type:      "struct",
		}
		for _, parameter := range clientFunc.Parameters {
			name = parameter.Name
			if parameter.Schema.Ref != "" {
				name = tgc.parseRefModelName(parameter.Schema.Ref)
			}
			field := gen.TemplateModelField{
				Name:    common.ToUpperCamelCase(name),
				Comment: parameter.Description,
				Tag:     fmt.Sprintf("`json:\"%s\" in:\"%s\"`", parameter.Name, parameter.In),
			}
			if model, ok := paramModelsMap[name]; ok {
				field.Type = model.ModelName
			} else {
				field.Type = tgc.transformGoType(parameter.Type)
			}
			f.RequestModel.TemplateModelFields = append(f.ResponseModel.TemplateModelFields, field)
		}
		tgc.ParamModels = append(tgc.ParamModels, f.RequestModel)
		for code, response := range clientFunc.Responses {
			if code != httptool.SWAGGER_TYPE__SUCCESS_CODE {
				continue
			}
			if response.Schema.Ref != "" {
				name := tgc.parseRefModelName(response.Schema.Ref)
				if model, ok := paramModelsMap[name]; ok {
					f.ResponseModel = model
				} else {
					return errors.New("response model parse error: " + name)
				}
			} else if response.Schema.Items.Ref != "" {
				name := tgc.parseRefModelName(response.Schema.Items.Ref)
				if model, ok := paramModelsMap[name]; ok {
					tmpModel := gen.TemplateModel{
						ModelName:           model.ModelName + "s",
						ModelComment:        "",
						Type:                "[]" + name,
						TemplateModelFields: nil,
					}
					tgc.ParamModels = append(tgc.ParamModels, tmpModel)

					f.ResponseModel = tmpModel
				} else {
					return errors.New("response model parse error: " + name)
				}
			} else if response.Schema.Type != "" {
				tmpModel := gen.TemplateModel{
					ModelName:           "[]byte",
					ModelComment:        "",
					Type:                "[]byte",
					TemplateModelFields: nil,
				}

				f.ResponseModel = tmpModel
			}
		}
		tgc.ClientFuncs = append(tgc.ClientFuncs, f)
	}
	return nil
}

func (tgc *TemplateGenClient) Parse(url string, isDebug bool) error {
	var body []byte
	var err error
	if !isDebug {
		req := httptool.NewHttpRequest(url, nil)
		body, err = req.SetTimeout(time.Second * 30).Get()
		if err != nil {
			panic(err)
		}
	} else {
		// for test
		file, err := os.Open("/data/golang/go/src/myprojects/tools/example/clients/gkspg-staging.json")
		if err != nil {
			return err
		}
		body, _ = ioutil.ReadAll(file)
	}

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

	tgc.InitTemplateFuncs()
	tgc.RegisteTemplateFunc(map[string]interface{}{
		"isModel": func(m gen.TemplateModel) bool {
			if m.ModelName == "" {
				return false
			}
			return true
		},
		"isStruct": func(m gen.TemplateModel) bool {
			if m.Type == "struct" || strings.Contains(m.Type, "[]") {
				return true
			}
			return false
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
