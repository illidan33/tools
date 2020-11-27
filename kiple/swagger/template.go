package swagger

import (
	"errors"
	"fmt"
	"github.com/dave/dst"
	"github.com/fatih/structtag"
	"github.com/illidan33/tools/common"
	"github.com/illidan33/tools/gen"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
)

type TemplateKipleSwagger struct {
	gen.GenTemplate
	gen.TemplateModel
	ModelList      map[string]gen.TemplateModel
	ImportList     map[string]string
	LoadImportList map[string]bool
	Swagger        SwaggerRoot
}

type SwaggerRoot struct {
	Swagger     string
	Info        SwaggerInfo                       `json:"info" yaml:"info"`
	Paths       map[string]map[string]SwaggerPath `json:"paths" yaml:"paths"`
	Definitions map[string]SwaggerDefinition      `json:"definitions" yaml:"definitions"`
}

type SwaggerInfo struct {
	Definition string   `json:"definition" yaml:"definition"`
	Title      string   `json:"title" yaml:"title"`
	Contact    []string `json:"contact" yaml:"contact"`
	Version    string   `json:"version" yaml:"version"`
}

type SwaggerPath struct {
	Consumes   []string                   `json:"consumes" yaml:"consumes"`
	Produces   []string                   `json:"produces" yaml:"produces"`
	Tags       []string                   `json:"tags" yaml:"tags"`
	Summary    string                     `json:"summary" yaml:"summary"`
	Parameters []SwaggerPathParam         `json:"parameters" yaml:"parameters"`
	Response   map[string]SwaggerPathResp `json:"response" yaml:"response"`
}

type SwaggerPathParam struct {
	Description string `json:"description" yaml:"description"`
	Name        string `json:"name" yaml:"name"`
	In          string `json:"in" yaml:"in"`
	Required    bool   `json:"required" yaml:"required"`
	SwaggerDefinitionProperty
}

type SwaggerPathResp struct {
	Description string            `json:"description" yaml:"description"`
	Schema      map[string]string `json:"schema" yaml:"schema"`
}

type SwaggerDefinition struct {
	Type       string                               `json:"type" yaml:"type"`
	Properties map[string]SwaggerDefinitionProperty `json:"properties" yaml:"properties"`
}

type SwaggerDefinitionProperty struct {
	Type   string            `json:"type" yaml:"type"`
	Schema map[string]string `json:"schema" yaml:"schema"`
	Items  map[string]string `json:"items" yaml:"items"`
}

var SwagTypeMap = map[string]string{
	"int8":    "integer",
	"int16":   "integer",
	"int":     "integer",
	"int32":   "integer",
	"int64":   "integer",
	"uint8":   "integer",
	"uint16":  "integer",
	"uint":    "integer",
	"uint32":  "integer",
	"uint64":  "integer",
	"float32": "integer",
	"float64": "integer",
	"bool":    "boolean",
	"string":  "string",
}

func (tm *TemplateKipleSwagger) getSwagSchema(modelPackage string, propertyType string) SwaggerDefinitionProperty {
	tmp := SwaggerDefinitionProperty{
		Schema: map[string]string{},
		Items:  map[string]string{},
	}
	if fType, ok := SwagTypeMap[propertyType]; ok {
		tmp.Type = fType
		tmp.Schema = nil
		tmp.Items = nil
	} else if strings.HasPrefix(propertyType, "[]") {
		tmp.Type = "array"
		tmp.Items = map[string]string{
			"$ref": "#/definitions/" + modelPackage + "." + strings.TrimPrefix(propertyType, "[]"),
		}
	} else {
		tmp.Type = "object"
		tmp.Schema = map[string]string{
			"$ref": "#/definitions/" + modelPackage + "." + propertyType,
		}
	}
	return tmp
}

func (tm *TemplateKipleSwagger) parseModelList(t *dst.Ident) error {
	if path, ok := tm.ImportList[t.Name]; ok {
		if v, ok := tm.LoadImportList[t.Name]; !ok || !v {
			mList, err := tm.ParseModelPackage(path)
			if err != nil {
				return err
			}
			for _, model := range mList {
				tm.ModelList[model.ModelName] = model

				swagDef := SwaggerDefinition{
					Type:       "object",
					Properties: map[string]SwaggerDefinitionProperty{},
				}
				for _, field := range model.TemplateModelFields {
					swagDef.Properties[field.Name] = tm.getSwagSchema(model.Package, field.Type)
				}
				tm.Swagger.Definitions[fmt.Sprintf("%s.%s", model.Package, model.ModelName)] = swagDef
			}
			tm.LoadImportList[t.Name] = true
		}
	}
	return nil
}

func (tm *TemplateKipleSwagger) parseSwaggerParam(t *dst.SelectorExpr) (paramList []SwaggerPathParam, err error) {
	pkgIdent, ok := t.X.(*dst.Ident)
	if !ok || pkgIdent == nil {
		err = errors.New("model package is empty, should not put model in service file.")
		return
	}
	err = tm.parseModelList(pkgIdent)
	if err != nil {
		return
	}

	model := tm.ModelList[t.Sel.Name]
	for _, field := range model.TemplateModelFields {
		tags, e := structtag.Parse(field.Tag)
		if e != nil {
			err = e
			return
		}
		inStr := ""
		in, _ := tags.Get("in")
		if in != nil {
			inStr = in.Name
		}
		requiredTag, _ := tags.Get("require")
		requiredStr := ""
		required := false
		if requiredTag != nil {
			requiredStr = requiredTag.Name
			required, _ = strconv.ParseBool(requiredStr)
		}
		swaggerPm := SwaggerPathParam{
			Description:               field.Comment,
			Name:                      field.Name,
			In:                        inStr,
			Required:                  required,
			SwaggerDefinitionProperty: tm.getSwagSchema(model.Package, field.Type),
		}

		if field.Type == "struct" {
			swaggerPm.Schema = map[string]string{
				"$ref": "#/definitions/request." + field.Type,
			}
		}
		paramList = append(paramList, swaggerPm)
	}
	return
}

func (tm *TemplateKipleSwagger) parseSwaggerResp(t *dst.SelectorExpr) (respList map[string]SwaggerPathResp, err error) {
	pkgIdent, ok := t.X.(*dst.Ident)
	if !ok || pkgIdent == nil {
		return
	}
	err = tm.parseModelList(pkgIdent)
	if err != nil {
		return
	}
	model, ok := tm.ModelList[t.Sel.Name]
	if !ok {
		err = errors.New(t.Sel.Name + " not in model list")
		return
	}
	respList = map[string]SwaggerPathResp{
		"200": {
			Description: "OK",
			Schema: map[string]string{
				"$ref": "#/definitions/" + model.Package + "." + model.ModelName,
			},
		},
	}

	return
}

func (tm *TemplateKipleSwagger) parseFuncDef(fc *dst.FuncDecl) error {
	api := SwaggerPath{
		Consumes:   []string{"application/json"},
		Produces:   []string{"application/json"},
		Tags:       []string{},
		Summary:    strings.Join(fc.Decs.NodeDecs.Start, ","),
		Parameters: []SwaggerPathParam{},
		Response:   map[string]SwaggerPathResp{},
	}
	for _, pm := range fc.Type.Params.List {
		if t, ok := pm.Type.(*dst.SelectorExpr); ok {
			if strings.Contains(t.Sel.Name, "Request") {
				paramList, err := tm.parseSwaggerParam(t)
				if err != nil {
					return err
				}
				api.Parameters = paramList
			} else if strings.Contains(t.Sel.Name, "Response") {
				var err error
				api.Response, err = tm.parseSwaggerResp(t)
				if err != nil {
					return err
				}
			}
		}
	}
	tm.Swagger.Paths["url"] = map[string]SwaggerPath{
		"post": api,
	}
	return nil
}

func (tm *TemplateKipleSwagger) ParseServiceDir(dir string) error {
	if !common.IsExists(dir) {
		return errors.New(dir + " is not exist")
	}

	rd, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, fi := range rd {
		if fi.IsDir() {
			tm.ParseServiceDir(fi.Name())
		} else {
			dstFile, err := tm.GetDstTree(filepath.Join(dir, fi.Name()))
			if err != nil {
				return err
			}

			for _, spec := range dstFile.Imports {
				path := strings.Trim(spec.Path.Value, "\"")
				base := filepath.Base(path)
				tm.ImportList[base] = path
			}
			for _, decl := range dstFile.Decls {
				if fc, ok := decl.(*dst.FuncDecl); ok {
					if fc.Name != nil && fc.Name.Name != "" && common.IsUpperLetter(rune(fc.Name.Name[0])) {
						err = tm.parseFuncDef(fc)
						if err != nil {
							return err
						}
					}
				}
			}
		}
	}

	return nil
}
