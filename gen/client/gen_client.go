package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"myprojects/tools/common"
	"myprojects/tools/gen"
	"net/http"
	"strings"
)

type CmdGenClient struct {
	CmdGenClientDocUrl      string
	CmdGenClientServiceName string
}

func (flags CmdGenClient) CmdHandle() {
	tpData := TemplateGenClient{}
	res, err := http.Get("http://192.168.1.116:8080/swagger/swagger/doc.json")
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	bf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	data := GenClientSwagger{}
	err = json.Unmarshal(bf, &data)
	if err != nil {
		panic(err)
	}

	clientFuncs := make([]GenClientFunc, 0)
	for name, path := range data.Paths {
		for method, funcItem := range path {
			funcItem.Name = name
			funcItem.Method = method
			clientFuncs = append(clientFuncs, funcItem)
		}
	}
	tpData.PackageName = common.ToLowerSnakeCase("client_" + flags.CmdGenClientServiceName)
	tpData.ClientModel.ModelName = common.ToUpperCamelCase(tpData.PackageName)
	tpData.ClientModel.ModelComment = data.Info.Description

	for name, def := range data.Definitions {
		s := strings.Index(name, ".")
		if s != -1 {
			name = name[s+1:]
		}
		f := gen.TemplateModel{
		}
		if def.Type == "object" {
			i := 0
			for proName, pro := range def.Properties {
				fName := common.ToUpperCamelCase(proName)
				f.TemplateModelFields = append(f.TemplateModelFields, gen.TemplateModelField{
					Name:    fName,
					Type:    pro.Type,
					Tag:     fmt.Sprintf("`json:\"%s\"`", proName),
					Comment: pro.Description,
				})
				i++
			}
		}
		tpData.ParamModels = append(tpData.ParamModels, f)
	}

	for _, clientFunc := range clientFuncs {
		name := clientFunc.Name
		s := strings.LastIndex(name, "/")
		if s != -1 {
			name = name[s+1:]
		}
		f := TemplateGenClientFunc{
			TemplateGenModelFunc: gen.TemplateGenModelFunc{},
			Args:                 []TemplateGenClientFuncParam{},
			Path:                 clientFunc.Name,
		}
		for _, parameter := range clientFunc.Parameters {
			param := TemplateGenClientFuncParam{
				Description: parameter.Description,
				Name:        parameter.Name,
				In:          parameter.In,
				Required:    parameter.Required,
			}
			if len(parameter.Schema) > 0 {
				if ref, ok := parameter.Schema["$ref"]; ok {
					s := strings.LastIndex(ref, ".")
					if s == -1 {
						s = strings.LastIndex(ref, "/")
					}
					if s != -1 {
						param.SchemaName = ref[s+1:]
					} else {
						param.SchemaName = ref
					}
				}
			}
			f.Args = append(f.Args, param)
		}
		for _, response := range clientFunc.Responses {
			if response.Schema.Type == "array" {

			}
			//f.Returns = append(f.Returns, response.Schema)
		}
		tpData.ClientFuncs = append(tpData.ClientFuncs, f)
	}

	//for i, i2 := range data.{
	//
	//}

	fmt.Printf("success")
}
