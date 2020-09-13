package gen

import (
	"bytes"
	"fmt"
	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"html/template"
	"io/ioutil"
	"myprojects/tools/gen/types"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type CmdGen interface {
	CmdHandle(cwd string)
}

type TemplateGenModle struct {
	TemplateGenStruct
	ModleStructFields map[string]TemplateGenStructField
}

type TemplateGenStruct struct {
	PackageName   string
	PackageList   map[string]string
	StructName    string
	StructComment string
	TemplateFuncs map[string]interface{}
}

type TemplateGenStructField struct {
	Name    string
	Type    string
	TypeLen uint64
	Tag     string
	Comment string
	Sort    int
}

type TemplateGenModleFunc struct {
	Name           string
	BelongToStruct string
	Args           []TemplateGenStructField
	Returns        []string
}

func (tfm *TemplateGenModle) Registe(data map[string]interface{}) {
	for name, fc := range data {
		tfm.TemplateFuncs[name] = fc
	}
}

func (tfm *TemplateGenModle) Init() {
	d := map[string]interface{}{
		"var":   func(s string) string { return ToLowerCamelCase(s) },
		"type":  func(s string) string { return ToUpperCamelCase(s) },
		"snake": func(s string) string { return ToLowerSnakeCase(s) },
		"printf": func(s string, args ...interface{}) string {
			return fmt.Sprintf(s, args...)
		},
		"html": func(s string) template.HTML {
			return template.HTML(s)
		},
		"sort": func(fields map[string]TemplateGenStructField) []TemplateGenStructField {
			ts := make([]TemplateGenStructField, len(fields))
			i := 0
			for _, field := range fields {
				ts[i] = field
				i++
			}
			sort.Slice(ts, func(i, j int) bool {
				if ts[i].Sort < ts[j].Sort {
					return true
				}
				return false
			})
			return ts
		},
	}
	tfm.Registe(d)
}

func (tfm *TemplateGenModle) ParseTemplate(templateTxt string) (templateSource *bytes.Buffer, e error) {
	templateSource = bytes.NewBuffer([]byte(""))

	tp := template.New("sql_to_modle")
	tp.Funcs(tfm.TemplateFuncs)
	if tp, e = tp.Parse(templateTxt); e != nil {
		return
	}
	e = tp.Execute(templateSource, *tfm)
	return
}

func (tfm *TemplateGenModle) FormatCodeToFile(filePath string, templateSource *bytes.Buffer) (err error) {
	var file *os.File
	filePath, _ = filepath.Abs(filePath)
	file, err = os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	f, e := decorator.Parse(templateSource.String())
	if e != nil {
		err = e
		return
	}
	err = decorator.Fprint(file, f)
	return
}

func (tfm *TemplateGenModle) GetDstTree(filePath string) (*dst.File, error) {
	var file *os.File
	var err error
	if IsExists(filePath) {
		file, err = os.Open(filePath)
		if err != nil {
			return nil, err
		}
	} else {
		file, err = os.Create(filePath)
		if err != nil {
			return nil, err
		}
	}
	defer file.Close()
	templateSource, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	f, err := decorator.Parse(templateSource)
	if err != nil {
		panic(err)
	}
	return f, nil
}

func (tfm *TemplateGenModle) TransformGormToModle(gormTable GormTable, hasGorm bool, isSimpleGorm bool, hasJson bool, hasDefault bool) error {
	tfm.StructName = gormTable.Name

	// fields
	for _, field := range gormTable.Fields {
		tgsf := TemplateGenStructField{
			Name:    ToUpperCamelCase(field.Name),
			Type:    "",
			TypeLen: 0,
			Tag:     "",
			Comment: field.Comment,
			Sort:    field.ModleSort,
		}

		tgsf.Type = field.Type
		fs := strings.Index(field.Type, "(")
		if fs != -1 {
			tgsf.Type = field.Type[:fs]
			fe := GetDataBetweenFlag(field.Type, "(", ")")
			ftLen, err := strconv.ParseUint(fe, 10, 64)
			if err != nil {
				return fmt.Errorf("Parse field type length error: %s", err.Error())
			}
			tgsf.TypeLen = ftLen
		}

		if v, ok := FieldType[strings.ToUpper(tgsf.Type)]; ok {
			tgsf.Type = v
			if field.IsUnsigned {
				tgsf.Type = "u" + tgsf.Type
			}
			if tgsf.Type == "time.Time" {
				tfm.PackageList["time"] = "time"
			}
		} else {
			return fmt.Errorf("Field type string not in map: %s", tgsf.Type)
		}

		tgsf.Tag = "`"
		null := "NOT NULL"
		if field.IsNull {
			null = "NULL"
		}
		if hasGorm {
			if isSimpleGorm {
				tgsf.Tag = fmt.Sprintf("%sgorm:\"column:%s\"", tgsf.Tag, field.Name)
			} else {
				tgsf.Tag = fmt.Sprintf("%sgorm:\"column:%s;type:%s;%s;default:%s\"", tgsf.Tag, field.Name, field.Type, null, field.Default)
			}
		}
		if hasJson {
			tgsf.Tag = fmt.Sprintf("%s json:\"%s\"", tgsf.Tag, ToLowerCamelCase(field.Name))
		}
		if hasDefault {
			tgsf.Tag = fmt.Sprintf("%s default:\"%s\"", tgsf.Tag, field.Default)
		}

		tgsf.Tag += "`"
		tfm.ModleStructFields[tgsf.Name] = tgsf
	}

	// indexs
	indexs := make([]GormIndex, len(gormTable.Indexs))
	i := 0
	for _, index := range gormTable.Indexs {
		indexs[i] = index
		i++
	}
	sort.Slice(indexs, func(i, j int) bool {
		if indexs[i].IndexSort < indexs[j].IndexSort {
			return true
		}
		return false
	})
	if len(gormTable.Indexs) > 0 {
		for _, index := range indexs {
			str := ""
			for kk, field := range index.Fields {
				if kk == 0 {
					if index.Type == types.INDEXTYPE__PRIMARY {
						str = fmt.Sprintf("\n// @def %s %s", index.Type.KeyLowerString(), ToUpperCamelCase(field.Name))
					} else {
						str = fmt.Sprintf("\n// @def %s:%s %s", index.Type.KeyLowerString(), index.Name, ToUpperCamelCase(field.Name))
					}
				} else {
					str = fmt.Sprintf("%s %s", str, ToUpperCamelCase(field.Name))
				}
			}
			tfm.StructComment += str
		}
	}

	return nil
}

func (tfm *TemplateGenModle) GenFuncName(fields []TemplateGenStructField) string {
	str := ""
	for i, f := range fields {
		if i == 0 {
			str = f.Name
		} else {
			str += "_And" + f.Name
		}
	}
	return ToUpperCamelCase(str)
}
