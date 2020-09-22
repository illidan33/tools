package gen

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"myprojects/tools/common"
	"myprojects/tools/gen/util/types"
	"os"
	"sort"
	"strconv"
	"strings"
)

type GormTable struct {
	Name          string
	Engine        string
	Charset       string
	Collate       string
	Comment       string
	AutoIncrement uint64
	Fields        []GormField
	Indexs        []GormIndex
}

type GormField struct {
	Name            string
	Type            string
	IsUnsigned      bool
	IsNull          bool
	Default         string
	IsAutoIncrement bool
	Character       string
	Collate         string
	Comment         string
	IsKeyField      bool
	KeyName         string
	KeyType         types.IndexType
	ModelSort       int
}

type GormIndex struct {
	Name      string
	Fields    []GormIndexField
	Type      types.IndexType
	Using     string
	IndexSort int
}

type GormIndexField struct {
	GormField
	IndexFieldSort int
}

type GormFlags struct {
	HasGorm      bool
	IsSimpleGorm bool
	HasJson      bool
	HasDefault   bool
}

var FieldType = map[string]string{
	"INTEGER":    "int64",
	"INT":        "int64",
	"SMALLINT":   "int64",
	"TINYINT":    "int8",
	"MEDIUMINT":  "int64",
	"BIGINT":     "int64",
	"CHAR":       "string",
	"VARCHAR":    "string",
	"BINARY":     "string",
	"VARBINARY":  "string",
	"FLOAT":      "float32",
	"DOUBLE":     "float64",
	"TINYTEXT":   "string",
	"MEDIUMTEXT": "string",
	"TEXT":       "string",
	"LONGTEXT":   "string",
	"ENUM":       "int8",
	"SET":        "int8",
	"TINYBLOB":   "[]byte",
	"MEDIUMBLOB": "[]byte",
	"BLOB":       "[]byte",
	"LONGBLOB":   "[]byte",
	"DATE":       "time.Time",
	"TIME":       "time.Time",
	"DATETIME":   "time.Time",
	"TIMESTAMP":  "time.Time",
	"YEAR":       "time.Time",
}

func (gt *GormTable) isCreateTitle(s string) bool {
	if strings.Contains(strings.ToUpper(s), types.SQLCONTENTTYPE__CREATE_TABLE) {
		return true
	}
	return false
}

func (gt *GormTable) isTableKey(s string) bool {
	if strings.Contains(strings.ToUpper(s), types.SQLCONTENTTYPE__KEY) {
		return true
	}
	return false
}

func (gt *GormTable) isPrimaryKey(s string) bool {
	if strings.Contains(strings.ToUpper(s), types.SQLCONTENTTYPE__PRIMARY_KEY) {
		return true
	}
	return false
}

func (gt *GormTable) isUniqueKey(s string) bool {
	if strings.Contains(strings.ToUpper(s), types.SQLCONTENTTYPE__UNIQUE_KEY) {
		return true
	}
	return false
}

func (gt *GormTable) isEngineEnd(s string) bool {
	if strings.Contains(strings.ToUpper(s), types.SQLCONTENTTYPE__ENGINE) {
		return true
	}
	return false
}

func (gt *GormTable) isComent(s string) bool {
	return strings.ToUpper(s) == types.SQLCONTENTTYPE__COMMENT
}

func (gt *GormTable) parseLineToTokens(s string) (rs []string, e error) {
	rs = make([]string, 0)
	tmp := bytes.Buffer{}
	commentS := false
	keyS := false
	for i := 0; i < len(s); i++ {
		if s[i] == ' ' || s[i] == ',' {
			ts := tmp.String()
			if ts != "" {
				if commentS || keyS {
					tmp.WriteByte(s[i])
				} else {
					if gt.isComent(ts) {
						commentS = true
					}
					rs = append(rs, tmp.String())
					tmp = bytes.Buffer{}
				}
			}
		} else {
			if s[i] == '(' {
				keyS = true
			} else if s[i] == ')' {
				keyS = false
			}
			tmp.WriteByte(s[i])
		}
	}
	if tmp.String() != "" {
		rs = append(rs, tmp.String())
		tmp = bytes.Buffer{}
	}
	return
}

func (gt *GormTable) parseIndexFieldString(gi *GormIndex, s string) error {
	arr := strings.Split(s, ",")
	fieldsMap := map[string]GormField{}
	for _, field := range gt.Fields {
		fieldsMap[field.Name] = field
	}
	for i, f := range arr {
		f = strings.Trim(strings.Trim(f, " "), "`")
		if v, ok := fieldsMap[f]; ok {
			v.IsKeyField = true
			v.KeyName = gi.Name
			v.KeyType = gi.Type
			tmp := GormIndexField{
				GormField:      v,
				IndexFieldSort: i,
			}
			gi.Fields = append(gi.Fields, tmp)

		} else {
			e := fmt.Errorf("Index field not map table field: %s", f)
			return e
		}
	}
	return nil
}

func (gt *GormTable) parseLineKey(s string, indexSortNum int) (isLineField bool, e error) {
	lineStrs, err := gt.parseLineToTokens(s)
	if err != nil {
		e = err
		return
	}

	// maybe line field
	if len(lineStrs) > 6 {
		isLineField = true
		return
	}

	using := ""
	for k, str := range lineStrs {
		if strings.ToUpper(str) == types.SQLCONTENTTYPE__USING {
			using = lineStrs[k+1]
			break
		}
	}
	gormIndex := GormIndex{
		Name:      "",
		Type:      types.INDEX_TYPE__INDEX,
		Using:     using,
		Fields:    []GormIndexField{},
		IndexSort: indexSortNum,
	}
	if gt.isPrimaryKey(s) {
		gormIndex.Type = types.INDEX_TYPE__PRIMARY
		fd := gt.getDataBetweenString(lineStrs[2], "(", ")")

		gt.parseIndexFieldString(&gormIndex, fd)
		gt.Indexs = append(gt.Indexs, gormIndex)
		return
	}
	gormIndex.Name = strings.Trim(lineStrs[1], "`")

	keyNameInx := 0
	for i, str := range lineStrs {
		strU := strings.ToUpper(str)
		if strU == types.SQLCONTENTTYPE__KEY {
			gormIndex.Name = strings.Trim(lineStrs[i+1], "`")
			keyNameInx = i + 1
		} else if strU == types.SQLCONTENTTYPE__USING {
			gormIndex.Using = lineStrs[i+1]
		}
	}

	if strings.Contains(strings.ToUpper(s), types.SQLCONTENTTYPE__UNIQUE_KEY) {
		gormIndex.Type = types.INDEX_TYPE__UNIQUE_INDEX
	}

	if keyNameInx == 0 {
		e = fmt.Errorf("Can not find name of index: %s", lineStrs)
		return
	}
	fieldStr := gt.getDataBetweenString(lineStrs[keyNameInx+1], "(", ")")
	if fieldStr == "" {
		e = fmt.Errorf("Field string is empty")
		return
	}

	gt.parseIndexFieldString(&gormIndex, fieldStr)
	if gt.Indexs == nil {
		gt.Indexs = make([]GormIndex, 0)
	}
	gt.Indexs = append(gt.Indexs, gormIndex)

	return
}

func (gt *GormTable) parseLineField(s string, sort int) error {
	lineStrs, err := gt.parseLineToTokens(s)
	if err != nil {
		return err
	}
	if len(lineStrs) == 0 {
		return errors.New("Line string array is empty")
	}
	gormField := GormField{
		ModelSort: sort,
	}
	gormField.Name = strings.Trim(lineStrs[0], "`")

	for i := 1; i < len(lineStrs); i++ {
		if i == 1 {
			gormField.Type = lineStrs[i]
		} else {
			fc := strings.ToUpper(lineStrs[i])
			switch fc {
			case types.SQLCONTENTTYPE__UNSIGNED:
				gormField.IsUnsigned = true
			case types.SQLCONTENTTYPE__AUTO__INCREMENT:
				gormField.IsAutoIncrement = true
			case types.SQLCONTENTTYPE__CHARACTER:
				if lineStrs[i+1] == types.SQLCONTENTTYPE__SET {
					gormField.Character = lineStrs[i+2]
				} else {
					return fmt.Errorf("Parse character of field error: %v", lineStrs)
				}
			case types.SQLCONTENTTYPE__NULL:
				if lineStrs[i-1] != types.SQLCONTENTTYPE__NOT {
					gormField.IsNull = true
				}
			case types.SQLCONTENTTYPE__DEFAULT:
				if lineStrs[i+1] == types.SQLCONTENTTYPE__NULL {
					gormField.IsNull = true
				} else {
					gormField.Default = strings.Trim(lineStrs[i+1], "'")
				}
			case types.SQLCONTENTTYPE__COMMENT:
				gormField.Comment = strings.Trim(strings.Trim(lineStrs[i+1], "'"), ",")
			default:
			}
		}
	}
	if gt.Fields == nil {
		gt.Fields = make([]GormField, 0)
	}
	gt.Fields = append(gt.Fields, gormField)
	return nil
}

func (gt *GormTable) parseLineTableTitle(s string) error {
	lineStrs, err := gt.parseLineToTokens(s)
	if err != nil {
		return err
	}
	if len(lineStrs) < 3 {
		return fmt.Errorf("parseLineToTokens error, number not enough: %v", lineStrs)
	}
	gt.Name = gt.trimField(lineStrs[2], "`")
	return nil
}

func (gt *GormTable) parseEngineEnd(s string) error {
	lineStrs, err := gt.parseLineToTokens(s)
	if err != nil {
		return err
	}
	for _, str := range lineStrs {
		s := strings.Index(str, "=")
		if s != -1 {
			switch strings.ToUpper(str[:s]) {
			case "ENGINE":
				gt.Engine = str[s+1:]
			case "CHARSET":
				gt.Charset = str[s+1:]
			case "COMMENT":
				gt.Comment = str[s+1:]
			case "COLLATE":
				gt.Collate = str[s+1:]
			case "AUTO_INCREMENT":
				ai, err := strconv.ParseUint(str[s+1:], 10, 64)
				if err != nil {
					return err
				}
				gt.AutoIncrement = ai
			default:
			}
		}
	}
	return nil
}

func (gt *GormTable) getDataBetweenString(source string, flag1 string, flag2 string) string {
	s := strings.Index(source, flag1)
	if s == -1 {
		return ""
	}
	e := strings.LastIndex(source, flag2)
	if e == -1 {
		return ""
	}
	return source[s+1 : e]
}

func (gt *GormTable) trimField(source string, flag string) string {
	return strings.Trim(source, flag)
}

func (gt *GormTable) Parse(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	bf := bufio.NewReader(file)
	fieldSortNum := 0
	indexSortNum := 0
	for {
		line, _, err := bf.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		ls := strings.Trim(string(line), " ")
		if len(ls) == 0 {
			continue
		}
		if gt.isCreateTitle(ls) {
			err = gt.parseLineTableTitle(ls)
			if err != nil {
				return err
			}
		} else if gt.isEngineEnd(ls) {
			err = gt.parseEngineEnd(ls)
			if err != nil {
				return err
			}
		} else if gt.isTableKey(ls) {
			islineField, err := gt.parseLineKey(ls, indexSortNum)
			if err != nil {
				return err
			}
			if islineField {
				fieldSortNum++
				err = gt.parseLineField(ls, fieldSortNum)
				if err != nil {
					return err
				}
			} else {
				indexSortNum++
			}
		} else {
			fieldSortNum++
			err = gt.parseLineField(ls, fieldSortNum)
			if err != nil {
				return err
			}
		}
	}
	sort.Slice(gt.Fields, func(i, j int) bool {
		if gt.Fields[i].ModelSort < gt.Fields[j].ModelSort {
			return true
		}
		return false
	})

	return nil
}

func (gt *GormTable) TransformGormToModel(tm *TemplateModel, gormFlags GormFlags) (packageList []string, err error) {
	tm.ModelName = gt.Name

	// fields
	for _, field := range gt.Fields {
		tgsf := TemplateModelField{
			Name:    common.ToUpperCamelCase(field.Name),
			Type:    "",
			Tag:     "",
			Comment: field.Comment,
		}

		tgsf.Type = field.Type
		fs := strings.Index(field.Type, "(")
		if fs != -1 {
			tgsf.Type = field.Type[:fs]
		}

		if v, ok := FieldType[strings.ToUpper(tgsf.Type)]; ok {
			tgsf.Type = v
			if field.IsUnsigned {
				tgsf.Type = "u" + tgsf.Type
			}
			if tgsf.Type == "time.Time" {
				packageList = []string{"time"}
			}
		} else {
			err = fmt.Errorf("Field type string not in map: %s", tgsf.Type)
			return
		}

		tgsf.Tag = "`"
		null := types.SQLCONTENTTYPE__NOT_NULL
		if field.IsNull {
			null = types.SQLCONTENTTYPE__NULL
		}
		if gormFlags.HasGorm {
			if gormFlags.IsSimpleGorm {
				tgsf.Tag = fmt.Sprintf("%s%s:\"column:%s\"", tgsf.Tag, types.MODEL_TAG_TYPE__GORM, field.Name)
			} else {
				tgsf.Tag = fmt.Sprintf("%s%s:\"column:%s;type:%s;%s;default:%s\"", tgsf.Tag, types.MODEL_TAG_TYPE__GORM, field.Name, field.Type, null, field.Default)
			}
		}
		if gormFlags.HasJson {
			tgsf.Tag = fmt.Sprintf("%s %s:\"%s\"", tgsf.Tag, types.MODEL_TAG_TYPE__JSON, common.ToLowerCamelCase(field.Name))
		}
		if gormFlags.HasDefault {
			tgsf.Tag = fmt.Sprintf("%s %s:\"%s\"", tgsf.Tag, types.MODEL_TAG_TYPE__DEFAULT, field.Default)
		}

		tgsf.Tag += "`"
		if tm.TemplateModelFields == nil {
			tm.TemplateModelFields = make([]TemplateModelField, 0)
		}
		tm.TemplateModelFields = append(tm.TemplateModelFields, tgsf)
	}

	// indexs
	indexs := make([]GormIndex, len(gt.Indexs))
	i := 0
	for _, index := range gt.Indexs {
		indexs[i] = index
		i++
	}
	sort.Slice(indexs, func(i, j int) bool {
		if indexs[i].IndexSort < indexs[j].IndexSort {
			return true
		}
		return false
	})
	if len(gt.Indexs) > 0 {
		for _, index := range indexs {
			str := ""
			for kk, field := range index.Fields {
				if kk == 0 {
					if index.Type == types.INDEX_TYPE__PRIMARY {
						str = fmt.Sprintf("\n// @%s %s %s", types.SQLCONTENTTYPE__DEF, index.Type.KeyLowerString(), common.ToUpperCamelCase(field.Name))
					} else {
						str = fmt.Sprintf("\n// @def %s:%s %s", index.Type.KeyLowerString(), index.Name, common.ToUpperCamelCase(field.Name))
					}
				} else {
					str = fmt.Sprintf("%s %s", str, common.ToUpperCamelCase(field.Name))
				}
			}
			tm.ModelComment += str
		}
	}

	return
}
