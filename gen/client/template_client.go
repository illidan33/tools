package client

import (
	"myprojects/tools/gen"
)

const templateClientTxt = ``

type TemplateGenClient struct {
	gen.GenTemplate
	gen.TemplatePackage
	ClientModel gen.TemplateModel
	ClientFuncs []TemplateGenClientFunc
	ParamModels []gen.TemplateModel
}

type TemplateGenClientFunc struct {
	gen.TemplateGenModelFunc
	Args []TemplateGenClientFuncParam
	Path string
}

type TemplateGenClientFuncParam struct {
	Description string
	Name        string
	In          string
	Required    bool
	SchemaName  string
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
	Type       string                                  `json:"type"`
	Properties map[string]GenClientDefinitionPropertie `json:"properties"`
}
type GenClientDefinitionPropertie struct {
	Description string `json:"description"`
	Type        string `json:"type"`
}

type GenClientFunc struct {
	Name       string                           `json:"-"`
	Method     string                           `json:"-"`
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
	Description string                      `json:"description"`
	Schema      GenClientFuncResponseSchema `json:"schema"`
}
type GenClientFuncResponseSchema struct {
	Type  string            `json:"type"`
	Items map[string]string `json:"items"`
}
