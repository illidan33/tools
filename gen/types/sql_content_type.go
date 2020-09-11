package types

import "strings"

//go:generate stringer -type=SqlContentType
type SqlContentType uint8

const (
	SQLCONTENTTYPE__UNKNOWN SqlContentType = iota
	SQLCONTENTTYPE__CREATE_TABLE
	SQLCONTENTTYPE__ENGINE
	SQLCONTENTTYPE__KEY
	SQLCONTENTTYPE__PRIMARY_KEY
	SQLCONTENTTYPE__UNIQUE_KEY
	SQLCONTENTTYPE__USING
	SQLCONTENTTYPE__DEFAULT
	SQLCONTENTTYPE__CHARSET
	SQLCONTENTTYPE__COLLATE
	SQLCONTENTTYPE__AUTO__INCREMENT
	SQLCONTENTTYPE__COMMENT
	SQLCONTENTTYPE__NOT_NULL
	SQLCONTENTTYPE__NULL
	SQLCONTENTTYPE__NOT
	SQLCONTENTTYPE__CHARACTER_SET
	SQLCONTENTTYPE__CHARACTER
	SQLCONTENTTYPE__SET
	SQLCONTENTTYPE__UNSIGNED
)

func (i SqlContentType) KeyString() string {
	str := i.String()
	s := strings.Index(str, "__")
	if strings.Contains(str[s+2:], "__") {
		str = strings.ReplaceAll(str[s+2:], "__", "_")
	} else {
		str = strings.ReplaceAll(str[s+2:], "_", " ")
	}
	return str
}

func (i SqlContentType) KeyLowerString() string {
	return strings.ToLower(i.KeyString())
}
