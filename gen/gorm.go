package gen

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"tools/common"
	"tools/gen/util/types"
	"io"
	"io/ioutil"
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
	TypeLength      int64
	IsUnsigned      bool
	IsNull          types.GormFieldType
	Default         string
	IsAutoIncrement bool
	Character       string
	Collate         string
	Comment         string
	IsKeyField      bool
	KeyName         string
	KeyType         types.IndexType
	IndexFieldSort  int
	ModelSort       int
}

type GormIndex struct {
	Name      string
	Fields    []*GormField
	Type      types.IndexType
	Using     string
	IndexSort int
}

type GormFlags struct {
	HasGorm      bool
	IsSimpleGorm bool
	HasJson      bool
	HasDefault   bool
}

type GormTableList []GormTable

var FieldType = map[string]string{
	"INTEGER":    "int64",
	"INT":        "int64",
	"SMALLINT":   "int64",
	"TINYINT":    "int8",
	"MEDIUMINT":  "int64",
	"BIGINT":     "int64",
	"bool":       "int8",
	"boolean":    "int8",
	"CHAR":       "string",
	"VARCHAR":    "string",
	"BINARY":     "string",
	"VARBINARY":  "string",
	"FLOAT":      "float32",
	"DOUBLE":     "float64",
	"DECIMAL":    "float64",
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

var NumberMap = map[byte]bool{
	'0': true,
	'1': true,
	'2': true,
	'3': true,
	'4': true,
	'5': true,
	'6': true,
	'7': true,
	'8': true,
	'9': true,
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

func (gt *GormTable) isForeignKey(s string) bool {
	if strings.Contains(strings.ToUpper(s), types.SQLCONTENTTYPE__FOREIGN_KEY) {
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
	fieldsMap := map[string]*GormField{}
	for i := 0; i < len(gt.Fields); i++ {
		fieldsMap[gt.Fields[i].Name] = &gt.Fields[i]
	}
	for i, f := range arr {
		f = gt.getDataBetweenString(f, "`", "`")
		if v, ok := fieldsMap[f]; ok {
			v.IsKeyField = true
			v.KeyName = gi.Name
			v.KeyType = gi.Type
			v.IndexFieldSort = i
			gi.Fields = append(gi.Fields, v)

		} else {
			e := fmt.Errorf("Index field not map table field: %s", f)
			return e
		}
	}
	return nil
}

func (gt *GormTable) parseLineKey(s string, indexSortNum int) (e error) {
	s = gt.TrimField(s)
	lineStrs, err := gt.parseLineToTokens(s)
	if err != nil {
		e = err
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
		Fields:    []*GormField{},
		IndexSort: indexSortNum,
	}
	if gt.isPrimaryKey(s) {
		gormIndex.Type = types.INDEX_TYPE__PRIMARY
		fd := gt.getDataBetweenString(lineStrs[2], "(", ")")

		gt.parseIndexFieldString(&gormIndex, fd)
		gt.Indexs = append(gt.Indexs, gormIndex)
		return
	} else if gt.isForeignKey(s) {
		gormIndex.Name = gt.TrimField(lineStrs[1])
		gormIndex.Type = types.INDEX_TYPE__FOREIGN_INDEX
		fd := gt.getDataBetweenString(lineStrs[4], "(", ")")
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
	if !strings.Contains(lineStrs[0], "`") {
		return fmt.Errorf("Parse tokens to field error, first string is not field name: %#v", lineStrs)
	}
	gormField := GormField{
		ModelSort: sort,
	}
	gormField.Name = gt.TrimField(lineStrs[0])

	for i := 1; i < len(lineStrs); i++ {
		if i == 1 {
			gormField.Type = lineStrs[i]
			s := strings.Index(gormField.Type, "(")
			if s != -1 {
				e := strings.LastIndex(gormField.Type, ")")
				if e != -1 {
					length := gormField.Type[s+1 : e]
					quoteIndex := strings.Index(length, ",")
					if quoteIndex != -1 {
						length = length[:quoteIndex]
					}
					gormField.Type = gormField.Type[:s]
					gormField.TypeLength, err = strconv.ParseInt(length, 10, 64)
					if err != nil {
						return err
					}
				}
			}
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
					i += 2
				} else {
					return fmt.Errorf("Parse character of field error: %v", lineStrs)
				}
			case types.SQLCONTENTTYPE__NOT:
				if lineStrs[i+1] == types.SQLCONTENTTYPE__NULL {
					gormField.IsNull = types.GORM_FIELD_TYPE__NOT_NULL
					i++
				}
			case types.SQLCONTENTTYPE__NULL:
				gormField.IsNull = types.GORM_FIELD_TYPE__NULL
			case types.SQLCONTENTTYPE__DEFAULT:
				gormField.Default = strings.Trim(lineStrs[i+1], "'")
				i++
			case types.SQLCONTENTTYPE__COMMENT:
				gormField.Comment = strings.Trim(strings.Trim(lineStrs[i+1], "'"), ",")
				i++
			case types.SQLCONTENTTYPE__COLLATE:
				gormField.Collate = lineStrs[i+1]
				i++
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
	s = gt.TrimLine(s)
	lineStrs, err := gt.parseLineToTokens(s)
	if err != nil {
		return err
	}
	if len(lineStrs) < 3 {
		return fmt.Errorf("parseLineToTokens error, number not enough: %v", lineStrs)
	}
	gt.Name = gt.TrimField(lineStrs[2])
	return nil
}

func (gt *GormTable) parseEngineEnd(s string) error {
	s = gt.TrimLine(s)
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
				gt.Charset = gt.TrimField(str[s+1:])
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

func (gt *GormTable) TrimLine(source string) string {
	source = strings.Trim(source, "\n")
	source = strings.Trim(source, "\t")
	source = strings.Trim(source, " ")
	return source
}

func (gt *GormTable) TrimField(source string) string {
	if source[0] == '(' {
		s := strings.Index(source, "(")
		e := strings.LastIndex(source, ")")
		if s != -1 && e != -1 {
			source = source[s+1 : e]
		}
	}
	source = strings.Trim(source, "`")
	return source
}

func (gt *GormTable) ParseStringToFieldOrKey(tx string, keyStart bool, fieldSort, indexSort int) (bool, int, int, error) {
	if gt.isPrimaryKey(tx) || keyStart {
		err := gt.parseLineKey(tx, indexSort)
		if err != nil {
			return false, 0, 0, err
		}
		indexSort++
		keyStart = true
	} else {
		err := gt.parseLineField(tx, fieldSort)
		if err != nil {
			return false, 0, 0, err
		}
		fieldSort++
	}
	return keyStart, fieldSort, indexSort, nil
}

func (gt *GormTable) parseStringToTable(s string) error {
	fieldStart := strings.Index(s, "(")
	fieldEnd := strings.LastIndex(s, ")")
	if fieldStart == -1 || fieldEnd == -1 || fieldStart == fieldEnd {
		return errors.New("Parse field content string error")
	}

	// parse title
	err := gt.parseLineTableTitle(s[:fieldStart])
	if err != nil {
		return err
	}
	err = gt.parseEngineEnd(s[fieldEnd+1:])
	if err != nil {
		return err
	}
	strStart := fieldStart + 1
	strEnd := strStart
	indexSort := 0
	fieldSort := 0
	keyStart := false
	brackets := 0     // like ()
	singleQuotes := 0 // like ''
	for {
		if strEnd >= fieldEnd {
			if strStart != strEnd {
				tx := s[strStart:strEnd]
				tx = gt.TrimLine(tx)
				keyStart, fieldSort, indexSort, err = gt.ParseStringToFieldOrKey(tx, keyStart, fieldSort, indexSort)
				if err != nil {
					return err
				}
			}
			break
		}
		if s[strEnd] != ',' {
			if s[strEnd] == '\'' {
				singleQuotes++
			} else if s[strEnd] == '(' {
				brackets++
			} else if s[strEnd] == ')' {
				brackets--
			}
			strEnd++
		} else if brackets == 0 && singleQuotes%2 == 0 {
			tx := s[strStart:strEnd]
			tx = gt.TrimLine(tx)
			keyStart, fieldSort, indexSort, err = gt.ParseStringToFieldOrKey(tx, keyStart, fieldSort, indexSort)
			if err != nil {
				return err
			}

			strEnd++
			strStart = strEnd
			brackets, singleQuotes = 0, 0
		} else {
			strEnd++
		}
	}

	return nil
}

func (gt *GormTableList) Parse(path string, gormFlags GormFlags) ([]TemplateModel, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	bf := bufio.NewReader(bytes.NewReader(content))
	gormList := make([]GormTable, 0)
	flag := false
	tmContent := bytes.Buffer{}
	for {
		line, _, err := bf.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		var tmp GormTable
		if tmp.isCreateTitle(string(line)) {
			flag = true
			tmContent.Write(line)
			tmContent.WriteByte('\n')
		} else if tmp.isEngineEnd(string(line)) {
			flag = false
			tmContent.Write(line)
			tmContent.WriteByte('\n')

			gorm := GormTable{}
			err = gorm.parseStringToTable(tmContent.String())
			if err != nil {
				return nil, err
			}
			gormList = append(gormList, gorm)
			tmContent = bytes.Buffer{}
		} else if flag {
			tmContent.Write(line)
			tmContent.WriteByte('\n')
		}
	}

	tms := make([]TemplateModel, 0)
	for _, gorm := range gormList {
		tm := TemplateModel{}
		gorm.transformGormToModel(&tm, gormFlags)
		tms = append(tms, tm)
	}

	return tms, nil
}

func (gt *GormTable) filterFieldName(s string) string {
	flag := 0
	for i := 0; i < len(s); i++ {
		if _, ok := NumberMap[s[i]]; !ok {
			flag = i
			break
		}
	}
	return s[flag:]
}

func (gt *GormTable) transformGormToModel(tm *TemplateModel, gormFlags GormFlags) (err error) {
	tm.ModelName = gt.Name

	// fields
	for _, field := range gt.Fields {
		noNumName := common.ToUpperCamelCase(gt.filterFieldName(field.Name))
		tgsf := TemplateModelField{
			Name:     noNumName,
			GormName: field.Name,
			JsonName: common.ToLowerSnakeCase(noNumName),
			Default:  field.Default,
			Type:     "",
			Tag:      "",
			Comment:  field.Comment,
		}

		tgsf.Type = field.Type

		if v, ok := FieldType[strings.ToUpper(tgsf.Type)]; ok {
			tgsf.Type = v
			if field.IsUnsigned {
				tgsf.Type = "u" + tgsf.Type
			}
		} else {
			err = fmt.Errorf("Field type string not in map: %s", tgsf.Type)
			return
		}

		tgsf.Tag = "`"
		null := ""
		if field.IsNull == types.GORM_FIELD_TYPE__NOT_NULL {
			null = types.SQLCONTENTTYPE__NOT_NULL + ";"
		} else if field.IsNull == types.GORM_FIELD_TYPE__NULL {
			null = types.SQLCONTENTTYPE__NULL + ";"
		}
		if gormFlags.HasGorm {
			if gormFlags.IsSimpleGorm {
				tgsf.Tag = fmt.Sprintf("%s%s:\"column:%s\"", tgsf.Tag, types.MODEL_TAG_TYPE__GORM, tgsf.GormName)
			} else {
				tgsf.Tag = fmt.Sprintf("%s%s:\"column:%s;type:%s;%sdefault:%s\"", tgsf.Tag, types.MODEL_TAG_TYPE__GORM, tgsf.GormName, field.Type, null, field.Default)
			}
		}
		if gormFlags.HasJson {
			tgsf.Tag = fmt.Sprintf("%s %s:\"%s\"", tgsf.Tag, types.MODEL_TAG_TYPE__JSON, tgsf.JsonName)
		}
		if gormFlags.HasDefault {
			tgsf.Tag = fmt.Sprintf("%s %s:\"%s\"", tgsf.Tag, types.MODEL_TAG_TYPE__DEFAULT, tgsf.Default)
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
			for i := 0; i < len(index.Fields); i++ {
				field := index.Fields[i]
				name := common.ToUpperCamelCase(gt.filterFieldName(field.Name))
				if i == 0 {
					if index.Type == types.INDEX_TYPE__PRIMARY {
						str = fmt.Sprintf("\n// @%s %s %s", types.SQLCONTENTTYPE__DEF, index.Type.KeyLowerString(), name)
					} else if index.Type == types.INDEX_TYPE__UNIQUE_INDEX || index.Type == types.INDEX_TYPE__INDEX || index.Type == types.INDEX_TYPE__FOREIGN_INDEX {
						str = fmt.Sprintf("\n// @%s %s:%s %s", types.SQLCONTENTTYPE__DEF, index.Type.KeyLowerString(), index.Name, name)
					}
				} else {
					str = fmt.Sprintf("%s %s", str, name)
				}
			}
			tm.ModelComment += str
		}
	}

	return
}
