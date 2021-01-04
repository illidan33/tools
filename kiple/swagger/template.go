package swagger

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/fatih/structtag"
	"github.com/ghodss/yaml"
	"go/format"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"tools/common"
	"tools/gen"
)

var templateSwagTxt = `// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = {{html $.Docs}}

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "{{$.Version}}",
	Host:        "",
	BasePath:    "",
	Schemes:     []string{},
	Title:       "{{$.Title}}",
	Description: "{{$.Description}}",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
`
var irisMethodMap = map[string]bool{
	"get":    true,
	"post":   true,
	"patch":  true,
	"delete": true,
	"put":    true,
	"option": true,
	"GET":    true,
	"POST":   true,
	"PATCH":  true,
	"DELETE": true,
	"PUT":    true,
	"OPTION": true,
}

var irisTagInMap = map[string]bool{
	"query":  true,
	"path":   true,
	"body":   true,
	"header": true,
}

type TemplateKipleSwagger struct {
	gen.GenTemplate
	gen.TemplateModel

	ModelList    map[string]gen.TemplateModel
	TemplateIris TemplateIris
	Swagger      SwaggerDocRoot
	IsInit       uint8
}

type TemplateIris struct {
	Parties      map[string]*TemplateIrisParty
	Applications map[string]*TemplateIrisApplication
	Controllers  map[string]TemplateIrisController
	Funcs        map[string]TemplateIrisFunc
}

type TemplateIrisParty struct {
	Name        string
	Url         string
	ParentParty *TemplateIrisParty
}

type TemplateIrisApplication struct {
	Name  string
	Party *TemplateIrisParty
}

type TemplateIrisController struct {
	Name string
	Url  string
}

type TemplateIrisFunc struct {
	BelongController *TemplateIrisController
	Url              string
	Method           string
	Consumes         []string
	Produces         []string
	Tag              string
	Summary          string
	Description      string
	FuncName         string
	Parameters       []SwaggerDefinitionProperty
	Responses        []SwaggerPathResp
}

type SwaggerRoot struct {
	Swagger     string                            `json:"swagger"`
	Info        SwaggerInfo                       `json:"info" yaml:"info"`
	Paths       map[string]map[string]SwaggerPath `json:"paths" yaml:"paths"`
	Definitions map[string]SwaggerDefinition      `json:"definitions" yaml:"definitions"`
}

type SwaggerDocRoot struct {
	Schemes  string `json:"schemes"`
	Host     string `json:"host"`
	BasePath string `json:"basePath"`
	SwaggerRoot
}

type SwaggerInfo struct {
	Description string      `json:"description" yaml:"description"`
	Title       string      `json:"title" yaml:"title"`
	Contact     interface{} `json:"contact" yaml:"contact"`
	Version     string      `json:"version" yaml:"version"`
}

type SwaggerPath struct {
	Consumes    []string                          `json:"consumes" yaml:"consumes"`
	Produces    []string                          `json:"produces" yaml:"produces"`
	Tags        []string                          `json:"tags" yaml:"tags"`
	Summary     string                            `json:"summary" yaml:"summary"`
	Description string                            `json:"description" yaml:"description"`
	Parameters  []map[string]interface{}          `json:"parameters" yaml:"parameters"`
	Responses   map[string]map[string]interface{} `json:"responses" yaml:"responses"`
}

// replaced by map[string]interface{}
type SwaggerPathParam struct {
	Description string `json:"description" yaml:"description"`
	Name        string `json:"name" yaml:"name"`
	In          string `json:"in" yaml:"in"`
	Required    bool   `json:"required" yaml:"required"`
	SwaggerDefinitionProperty
}

type SwaggerPathResp struct {
	ResponseCode string `json:"-" yaml:"-"`
	Description  string `json:"description" yaml:"description"`
	Model        gen.TemplateModel
}

type SwaggerDefinition struct {
	Type       string                 `json:"type" yaml:"type"`
	Properties map[string]interface{} `json:"properties" yaml:"properties"`
}

type SwaggerDefinitionProperty struct {
	Name        string
	In          string
	Required    bool
	Description string
	Type        string
	Package     string
}

var SwagTypeMap = map[string]string{
	"int8":      "integer",
	"int16":     "integer",
	"int":       "integer",
	"int32":     "integer",
	"int64":     "integer",
	"uint8":     "integer",
	"uint16":    "integer",
	"uint":      "integer",
	"uint32":    "integer",
	"uint64":    "integer",
	"float32":   "integer",
	"float64":   "integer",
	"bool":      "boolean",
	"string":    "string",
	"time.Time": "string",
}

func (tm *TemplateKipleSwagger) getSwagDefPropertyType(modelName, tp string) string {
	if fType, ok := SwagTypeMap[tp]; ok {
		return fType
	} else if strings.HasPrefix(tp, "[]") {
		sourceType := strings.TrimPrefix(tp, "[]")
		return tm.getSwagDefPropertyType(modelName, sourceType)
	} else {
		return "#/definitions/" + modelName + "." + tp
	}
}

func (tm *TemplateKipleSwagger) getSwagDefPropertyTagType(modelName, tp string) string {
	if fType, ok := SwagTypeMap[tp]; ok {
		return fType
	} else if strings.HasPrefix(tp, "[]") {
		sourceType := strings.TrimPrefix(tp, "[]")
		return "[]" + tm.getSwagDefPropertyTagType(modelName, sourceType)
	} else {
		return modelName + "." + tp
	}
}

func (tm *TemplateKipleSwagger) getSwagSchemaMap(modelName, propertyType string) map[string]interface{} {
	sch := map[string]interface{}{}
	if fType, ok := SwagTypeMap[propertyType]; ok {
		sch["type"] = fType
	} else if strings.HasPrefix(propertyType, "[]") {
		sourceType := strings.TrimPrefix(propertyType, "[]")
		sch["type"] = "array"
		sch["items"] = tm.getSwagSchemaMap(modelName, sourceType)
	} else if strings.HasPrefix(propertyType, "map[") {
		i := strings.Index(propertyType, "]")
		sourceType := propertyType[i+1:]
		sch["type"] = "object"
		sch["additionalProperties"] = tm.getSwagSchemaMap(modelName, sourceType)
	} else {
		//sch["type"] = "object"
		if v, ok := tm.ModelList[propertyType]; ok {
			swagDef := SwaggerDefinition{
				Type:       "object",
				Properties: map[string]interface{}{},
			}
			for _, field := range v.TemplateModelFields {
				if field.JsonName == "-" {
					continue
				}
				swagDef.Properties[field.JsonName] = tm.getSwagPropertity(v.Package, field)
			}
			tm.Swagger.Definitions[fmt.Sprintf("%s.%s", v.Package, v.ModelName)] = swagDef
		}

		return map[string]interface{}{
			"$ref": "#/definitions/" + modelName + "." + propertyType,
		}
	}
	return sch
}

func (tm *TemplateKipleSwagger) getSwagPropertity(pkg string, field gen.TemplateModelField) map[string]interface{} {
	pro := map[string]interface{}{}
	if fType, ok := SwagTypeMap[field.Type]; ok {
		pro["type"] = fType
	} else {
		pro["schema"] = tm.getSwagSchemaMap(pkg, field.Type)
	}
	pro["description"] = field.Comment
	tags, _ := structtag.Parse(field.Tag)
	val, _ := tags.Get("validate")
	if val != nil && strings.Trim(val.Value(), "\"") == "required" {
		pro["required"] = true
	}

	return pro
}

func (tm *TemplateKipleSwagger) getSwagReqParam(param SwaggerDefinitionProperty) map[string]interface{} {
	pro := map[string]interface{}{
		"name":        param.Name,
		"in":          param.In,
		"description": param.Description,
		"required":    param.Required,
	}
	if v, ok := tm.ModelList[param.Type]; ok {
		pro["schema"] = tm.getSwagSchemaMap(v.Package, v.Type)
	} else {
		pro["schema"] = tm.getSwagSchemaMap(param.Package, param.Type)
	}
	return pro
}

func (tm *TemplateKipleSwagger) getSwagResp(resp SwaggerPathResp) map[string]interface{} {
	pro := map[string]interface{}{}
	pro["schema"] = tm.getSwagSchemaMap(resp.Model.Package, resp.Model.Type)
	pro["description"] = resp.Description
	return pro
}

func (tm *TemplateKipleSwagger) parseFuncDef(controllerName, method, apiUrl, name, tag, comment string) error {
	method = strings.ToLower(strings.Trim(method, "\""))
	apiUrl = strings.Trim(apiUrl, "\"")
	name = strings.Trim(name, "\"")
	api := TemplateIrisFunc{
		BelongController: nil,
		Url:              apiUrl,
		Method:           method,
		Consumes:         []string{"application/json"},
		Produces:         []string{"application/json"},
		Tag:              tag,
		Summary:          name,
		Description:      comment,
		FuncName:         name,
		Parameters:       []SwaggerDefinitionProperty{},
		Responses:        nil,
	}
	if v, ok := tm.TemplateIris.Controllers[controllerName]; ok {
		api.BelongController = &v
	}
	var req gen.TemplateModel
	reqModelName := name + "Request"
	if v, ok := tm.ModelList[reqModelName]; ok {
		req = v
	}
	if req.ModelName != "" {
		for _, field := range req.TemplateModelFields {
			if field.JsonName == "-" {
				continue
			}
			tags, err := structtag.Parse(field.Tag)
			if err != nil {
				return err
			}
			inStr := ""
			in, _ := tags.Get("in")
			if in != nil {
				inStr = in.Name
				if _, ok := irisTagInMap[inStr]; !ok {
					return fmt.Errorf("tag [%s] of field [%s] is not support", inStr, field.Name)
				}
			} else if strings.ToUpper(method) == http.MethodGet {
				inStr = "query"
			} else {
				inStr = "body"
			}
			requiredTag, _ := tags.Get("validate")
			required := false
			if requiredTag != nil {
				if strings.Contains(requiredTag.Name, "required") {
					required = true
				}
			}
			swaggerPm := SwaggerDefinitionProperty{
				Name:        field.JsonName,
				In:          inStr,
				Required:    required,
				Description: field.Comment,
				Type:        field.Type,
				Package:     req.Package,
			}
			api.Parameters = append(api.Parameters, swaggerPm)
		}
	} else {
		reqModelName = name + "Body"
		if v, ok := tm.ModelList[reqModelName]; ok {
			req = v
			swaggerPm := SwaggerDefinitionProperty{
				Name:        "data",
				In:          "body",
				Required:    true,
				Description: req.ModelComment,
				Type:        req.ModelName,
				Package:     req.Package,
			}
			api.Parameters = append(api.Parameters, swaggerPm)
		}
	}

	api.Responses = make([]SwaggerPathResp, 0)
	respModelName := name + "Response"
	resp, ok := tm.ModelList[respModelName]
	if ok {
		tmpResp := SwaggerPathResp{
			ResponseCode: "200",
			Description:  "OK",
			Model:        resp,
		}
		api.Responses = append(api.Responses, tmpResp)
	}

	tm.TemplateIris.Funcs[fmt.Sprintf("%s-%s", controllerName, api.FuncName)] = api
	return nil
}

func (tm *TemplateKipleSwagger) ParsePojoDir(dir string) error {
	if !common.IsExists(dir) {
		return errors.New(dir + " is not exist")
	}

	mList, err := tm.ParseModelDir(dir, false)
	if err != nil {
		return err
	}
	for _, model := range mList {
		if _, ok := tm.ModelList[model.ModelName]; ok {
			return errors.New("ParsePojoDir - model name repeat.")
		}
		tm.ModelList[model.ModelName] = model
	}

	return nil
}

func (tm *TemplateKipleSwagger) genDstFileToFile(dstFilePath string, node *dst.File) (err error) {
	fset, af, e := decorator.RestoreFile(node)
	if e != nil {
		err = e
		return
	}

	var file *os.File
	file, err = os.OpenFile(dstFilePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	err = format.Node(file, fset, af)

	return
}

func (tm *TemplateKipleSwagger) ParseSwagTitle(file string) error {
	dstFile, err := tm.GetDstTree(file)
	if err != nil {
		return err
	}
	for _, decl := range dstFile.Decls {
		fun, ok := decl.(*dst.FuncDecl)
		if !ok {
			continue
		}
		if fun.Name.Name != "main" {
			continue
		}
		ttl, ver, desc := false, false, false
		for _, s := range fun.Decs.Start.All() {
			i := strings.Index(s, "@")
			if i == -1 {
				continue
			}
			s = s[i+1:]
			space := strings.Index(s, " ")
			if space == -1 {
				continue
			}
			v := strings.TrimSpace(s[space+1:])
			switch s[:space] {
			case "title":
				ttl = true
				tm.Swagger.Info.Title = v
			case "version":
				ver = true
				tm.Swagger.Info.Version = v
			case "description":
				desc = true
				tm.Swagger.Info.Description = v
			default:

			}
		}
		rootName := filepath.Base(filepath.Dir(file))
		if !ttl {
			fun.Decs.Start.Append(fmt.Sprintf("// @title %s", rootName))
		}
		if !ver {
			fun.Decs.Start.Append("// @version 1.0.0")
		}
		if !desc {
			fun.Decs.Start.Append(fmt.Sprintf("// @description This is api document of %s service.", rootName))
		}
	}

	flag := false
	for _, imp := range dstFile.Imports {
		if imp.Name == nil || imp.Name.Name != "_" {
			continue
		}
		path := strings.Trim(imp.Path.Value, "\"")
		if strings.HasSuffix(filepath.Dir(file), filepath.Dir(path)) && filepath.Base(path) == "docs" {
			flag = true
			break
		}
	}
	if !flag {
		projectPkg, err := common.GetPackageFromDir(filepath.Dir(file))
		if err != nil {
			return err
		}
		docImp := &dst.ImportSpec{
			Name: &dst.Ident{
				Name: "_",
				Obj:  nil,
				Path: "",
				Decs: dst.IdentDecorations{},
			},
			Path: &dst.BasicLit{
				Kind:  9,
				Value: fmt.Sprintf("\"%s/docs\"", projectPkg),
				Decs:  dst.BasicLitDecorations{},
			},
			Decs: dst.ImportSpecDecorations{
				NodeDecs: dst.NodeDecs{
					Before: 0,
					Start:  nil,
					End: dst.Decorations{
						"// docs is generated by Swag CLI, you have to import it.",
					},
					After: 0,
				},
				Name: nil,
			},
		}
		dstFile.Imports = append(dstFile.Imports, docImp)
	}
	err = tm.genDstFileToFile(file, dstFile)
	if err != nil {
		return errors.New("FormatCodeToFile - write file error: " + err.Error())
	}
	return nil
}

func (tm *TemplateKipleSwagger) parseIrisParty(p interface{}) (resp TemplateIrisParty, err error) {
	switch pt := p.(type) {
	case *dst.Ident:
		if v, ok := tm.TemplateIris.Parties[pt.Name]; ok {
			return *v, nil
		} else {
			err = errors.New("not found party in map:" + pt.Name)
			return
		}
	case *dst.CallExpr:
		selec, ok := pt.Fun.(*dst.SelectorExpr)
		if !ok {
			return
		}
		if selec.Sel.Name != "Party" {
			return
		}
		resp.Url = strings.Trim(pt.Args[0].(*dst.BasicLit).Value, "\"")
		parent := selec.X.(*dst.Ident).Name
		if v, ok := tm.TemplateIris.Parties[parent]; ok {
			resp.ParentParty = v
			if v.Url[0] != '/' {
				v.Url = "/" + v.Url
			}
			resp.Url = v.Url + resp.Url
		}
	default:

	}
	return
}

func (tm *TemplateKipleSwagger) parseIrisHandle(callReal *dst.CallExpr) (cons []*TemplateIrisController, p *TemplateIrisParty, err error) {
	cons = []*TemplateIrisController{}
	switch fun := callReal.Fun.(type) {
	case *dst.SelectorExpr:
		switch call := fun.X.(type) {
		case *dst.CallExpr:
			cons, p, err = tm.parseIrisHandle(call)
			if err != nil {
				return
			}
		case *dst.Ident:
			if call.Name == "mvc" && fun.Sel.Name == "New" {
				tmp, e := tm.parseIrisParty(callReal.Args[0])
				if e != nil {
					err = e
					return
				}
				p = &tmp
				return
			}
		}
		for _, argExpr := range callReal.Args {
			switch arg := argExpr.(type) {
			case *dst.UnaryExpr:
				conName := arg.X.(*dst.CompositeLit).Type.(*dst.Ident).Name
				tmp := TemplateIrisController{
					Name: conName,
					Url:  "",
				}
				if p != nil {
					tmp.Url = p.Url
				}
				cons = append(cons, &tmp)
			case *dst.CallExpr:
				var tmps []*TemplateIrisController
				tmps, p, err = tm.parseIrisHandle(arg)
				if err != nil {
					return
				}
				if len(tmps) > 0 {
					cons = append(cons, tmps...)
				}
			}
		}
	case *dst.Ident:
		// TODO(illidan/2020/12/21):
	default:
		err = errors.New("Cannot parse Expr")
		return
	}

	return
}

func (tm *TemplateKipleSwagger) parseIrisMvcConfig(call *dst.CallExpr) (err error) {
	if len(call.Args) < 2 {
		return
	}
	p, err := tm.parseIrisParty(call.Args[0])
	if err != nil {
		return err
	}
	cFun, ok := call.Args[1].(*dst.FuncLit)
	if !ok {
		return
	}
	for _, stmt := range cFun.Body.List {
		exprStmt1, ok := stmt.(*dst.ExprStmt)
		if !ok {
			continue
		}
		callExpr, ok := exprStmt1.X.(*dst.CallExpr)
		if ok {
			tmps, _, err := tm.parseIrisHandle(callExpr)
			if err != nil {
				return err
			}
			for _, tmp := range tmps {
				tmp.Url = p.Url
				tm.TemplateIris.Controllers[tmp.Name] = *tmp
			}
		}
	}
	return
}

func (tm *TemplateKipleSwagger) parseIrisMvcFuncName(cName, name, comment string) error {
	k := 0
	for i := 1; i < len(name); i++ {
		if common.IsUpperLetter(rune(name[i])) {
			k = i
			break
		}
	}
	method := strings.ToLower(name[:k])
	if _, ok := irisMethodMap[method]; !ok {
		return nil
	}
	//name = common.ToLowerCamelCase(name[k:])
	url := common.ToLowerSnakeCase(name[k:])
	err := tm.parseFuncDef(cName, method, url, name, cName, comment)
	if err != nil {
		return err
	}
	return nil
}

func (tm *TemplateKipleSwagger) ParseControllerDir(dir string) error {
	if !common.IsExists(dir) {
		return errors.New(dir + " is not exist")
	}
	rd, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, file := range rd {
		if file.IsDir() {
			tm.ParseControllerDir(file.Name())
		} else {
			dstFile, err := tm.GetDstTree(filepath.Join(dir, file.Name()))
			if err != nil {
				return err
			}
			for _, decl := range dstFile.Decls {
				fc, ok := decl.(*dst.FuncDecl)
				if !ok {
					continue
				}
				if fc.Name.Name == "BeforeActivation" {
					for _, stmt := range fc.Body.List {
						expr, ok := stmt.(*dst.ExprStmt)
						if !ok {
							continue
						}
						call, ok := expr.X.(*dst.CallExpr)
						if !ok {
							continue
						}
						if len(call.Args) < 3 {
							break
						}
						method, ok := call.Args[0].(*dst.BasicLit)
						url, ok := call.Args[1].(*dst.BasicLit)
						name, ok := call.Args[2].(*dst.BasicLit)
						comment := tm.ParseDstCommentFromNode(expr.Decs.NodeDecs, true)
						controllerName := tm.ParseDstNodeType(fc.Recv.List[0].Type)
						err = tm.parseFuncDef(controllerName, method.Value, url.Value, name.Value, controllerName, comment)
						if err != nil {
							return err
						}
					}
				} else if fc.Name.Name == "RegisterGlobalModel" {
					for _, stmt := range fc.Body.List {
						switch realStmt := stmt.(type) {
						case *dst.ExprStmt:
							call, ok := realStmt.X.(*dst.CallExpr)
							if !ok {
								continue
							}
							fun, ok := call.Fun.(*dst.SelectorExpr)
							if !ok {
								continue
							}
							switch fun.Sel.Name {
							case "Configure":
								tm.parseIrisMvcConfig(call)
							default:
								tmps, _, err := tm.parseIrisHandle(call)
								if err != nil {
									return err
								}
								for _, tmp := range tmps {
									tm.TemplateIris.Controllers[tmp.Name] = *tmp
								}
								continue
							}

						case *dst.AssignStmt:
							for k, lh := range realStmt.Lhs {
								if call, ok := realStmt.Rhs[k].(*dst.CallExpr); ok {
									p, err := tm.parseIrisParty(call)
									if err != nil {
										return err
									}
									p.Name = lh.(*dst.Ident).Name
									tm.TemplateIris.Parties[p.Name] = &p
								}
							}
						default:

						}
					}
				} else if fc.Recv != nil && len(fc.Recv.List) > 0 {
					// parse functions of controller
					controllerName := tm.ParseDstNodeType(fc.Recv.List[0].Type)
					if fun, ok := tm.TemplateIris.Funcs[fc.Name.Name]; ok && fun.BelongController != nil && fun.BelongController.Name == controllerName {
						var oldDecs string
						for _, v := range fc.Decs.Start.All() {
							if !strings.Contains(v, "// @") {
								oldDecs += v + "\n"
							}
						}
						oldDecs = strings.Trim(oldDecs, "\n")
						fun.Description = oldDecs
						tm.TemplateIris.Funcs[fc.Name.Name] = fun
					} else {
						comment := tm.ParseDstCommentFromNode(fc.Decs.NodeDecs, true)
						err = tm.parseIrisMvcFuncName(controllerName, fc.Name.Name, comment)
						if err != nil {
							return err
						}
					}
				}
			}

			if tm.IsInit != 0 {
				err = tm.genDstFileToFile(filepath.Join(dir, file.Name()), dstFile)
				if err != nil {
					return errors.New("FormatCodeToFile - write file error: " + err.Error())
				}
			}
		}
	}

	return nil
}

func (tm *TemplateKipleSwagger) OverWriteControllerDir(dir string) error {
	if tm.IsInit == 0 {
		return nil
	}
	rd, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, file := range rd {
		if file.IsDir() {
			tm.OverWriteControllerDir(file.Name())
		} else {
			dstFile, err := tm.GetDstTree(filepath.Join(dir, file.Name()))
			if err != nil {
				return err
			}
			for i := 0; i < len(dstFile.Decls); i++ {
				decl := dstFile.Decls[i]
				fc, ok := decl.(*dst.FuncDecl)
				if !ok {
					continue
				}
				if fc.Recv != nil && len(fc.Recv.List) > 0 {
					tag := tm.ParseDstNodeType(fc.Recv.List[0].Type)
					for _, path := range tm.TemplateIris.Funcs {
						if tag == path.Tag && fc.Name.Name == path.FuncName {
							var oldDecs string
							descp := bytes.Buffer{}
							// 1 写入-跳过已存在；2 覆盖写入；3 全部清除；
							if tm.IsInit == 1 {
								oldDecs = strings.Join(fc.Decs.Start.All(), "\n")
							} else if tm.IsInit == 2 || tm.IsInit == 3 {
								decs := fc.Decs.Start.All()
								fc.Decs.Start.Clear()
								for _, s := range decs {
									if !strings.HasPrefix(s, "// @") {
										fc.Decs.Start.Prepend(s)
										descp.WriteString(strings.Trim(s, "//"))
										descp.WriteString(", ")
									}
								}
								if tm.IsInit == 3 {
									continue
								}
							}
							if !strings.Contains(oldDecs, "// @Tags") {
								fc.Decs.Start.Append("// @Tags " + tag)
							}
							if !strings.Contains(oldDecs, "// @Summary") {
								fc.Decs.Start.Append("// @Summary " + path.FuncName)
							}
							if !strings.Contains(oldDecs, "// @Description") {
								fc.Decs.Start.Append("// @Description " + descp.String())
							}
							if !strings.Contains(oldDecs, "// @Accept") {
								fc.Decs.Start.Append("// @Accept json")
							}
							if !strings.Contains(oldDecs, "// @Produce") {
								fc.Decs.Start.Append("// @Produce json")
							}
							if !strings.Contains(oldDecs, "// @Param") {
								for _, param := range path.Parameters {
									tp := tm.getSwagDefPropertyTagType(param.Package, param.Type)
									if param.Description == "" {
										param.Description = param.In
									}
									fc.Decs.Start.Append(fmt.Sprintf("// @Param %s %s %s %t \"%s\"", param.Name, param.In, tp, param.Required, strings.TrimPrefix(param.Description, "//")))
								}
							}
							if !strings.Contains(oldDecs, "// @Success") {
								for _, resp := range path.Responses {
									typeMap := tm.getSwagDefPropertyTagType(resp.Model.Package, resp.Model.Type)
									if !strings.HasPrefix(resp.Model.Type, "[]") {
										fc.Decs.Start.Append(fmt.Sprintf("// @Success %s {object} %s", resp.ResponseCode, typeMap))
									} else {
										fc.Decs.Start.Append(fmt.Sprintf("// @Success %s {array} %s", resp.ResponseCode, typeMap))
									}
								}
							}
							if !strings.Contains(oldDecs, "// @Router") {
								fc.Decs.Start.Append("// @Router " + path.Url + " [" + path.Method + "]")
							}
							break
						}
					}
				}
			}
			err = tm.genDstFileToFile(filepath.Join(dir, file.Name()), dstFile)
			if err != nil {
				return errors.New("FormatCodeToFile - write file error: " + err.Error())
			}
		}
	}
	return nil
}

// Set Swagger Paths with new url
func (tm *TemplateKipleSwagger) SetSwaggerPaths() error {
	for k, irisFunc := range tm.TemplateIris.Funcs {
		if con, ok := tm.TemplateIris.Controllers[irisFunc.Tag]; ok {
			irisFunc.BelongController = &con
		}
		tm.TemplateIris.Funcs[k] = irisFunc
	}

	tm.Swagger.Paths = map[string]map[string]SwaggerPath{}
	for k, path := range tm.TemplateIris.Funcs {
		if path.BelongController == nil {
			return errors.New("not found controller:" + path.Tag)
		}
		api := SwaggerPath{
			Consumes:    path.Consumes,
			Produces:    path.Produces,
			Tags:        []string{path.Tag},
			Summary:     path.Summary,
			Description: path.Description,
			Parameters:  []map[string]interface{}{},
			Responses:   map[string]map[string]interface{}{},
		}
		newPath := path.Url
		if p, ok := tm.TemplateIris.Controllers[path.Tag]; ok {
			u := strings.Trim(p.Url, "/")
			path.Url = strings.Trim(path.Url, "/")
			newPath = fmt.Sprintf("/%s/%s", u, path.Url)
		}
		path.Url = newPath
		tm.TemplateIris.Funcs[k] = path

		for _, param := range path.Parameters {
			typeMap := tm.getSwagReqParam(param)
			api.Parameters = append(api.Parameters, typeMap)
		}
		for _, resp := range path.Responses {
			api.Responses[resp.ResponseCode] = tm.getSwagResp(resp)
		}

		if p, ok := tm.Swagger.Paths[newPath]; ok {
			p[path.Method] = api
			tm.Swagger.Paths[newPath] = p
		} else {
			tm.Swagger.Paths[newPath] = map[string]SwaggerPath{
				path.Method: api,
			}
		}
	}
	return nil
}

func (tm *TemplateKipleSwagger) FormatToFiles(cmdDir string) error {
	content, err := json.MarshalIndent(tm.Swagger.SwaggerRoot, "", "  ")
	if err != nil {
		return err
	}
	yamlContent, err := yaml.JSONToYAML(content)
	//yamlContent, err := yaml.Marshal(cmdtp.Template.Swagger.SwaggerRoot)
	if err != nil {
		return err
	}

	err = tm.WriteToFile(filepath.Join(cmdDir, "docs/swagger.json"), content)
	if err != nil {
		return err
	}
	err = tm.WriteToFile(filepath.Join(cmdDir, "docs/swagger.yaml"), yamlContent)
	if err != nil {
		return err
	}

	docContent, err := json.MarshalIndent(tm.Swagger, "", "  ")
	if err != nil {
		return err
	}
	bt, err := tm.ParseTemplate(templateSwagTxt, "templateSwagTxt", map[string]string{
		"Docs":        fmt.Sprintf("`%s`", string(docContent)),
		"Description": tm.Swagger.Info.Description,
		"Title":       tm.Swagger.Info.Title,
		"Version":     tm.Swagger.Info.Version,
	})
	if err != nil {
		return err
	}
	err = tm.FormatCodeToFile(filepath.Join(cmdDir, "docs/docs.go"), bt)
	if err != nil {
		return err
	}
	return nil
}
