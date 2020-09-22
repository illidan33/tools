package types

import "strings"

type SqlContentType string

const (
	SQLCONTENTTYPE__UNKNOWN         SqlContentType = ""
	SQLCONTENTTYPE__CREATE_TABLE                   = "CREATE TABLE"
	SQLCONTENTTYPE__ENGINE                         = "ENGINE"
	SQLCONTENTTYPE__KEY                            = "KEY"
	SQLCONTENTTYPE__PRIMARY_KEY                    = "PRIMARY KEY"
	SQLCONTENTTYPE__UNIQUE_KEY                     = "UNIQUE KEY"
	SQLCONTENTTYPE__USING                          = "USING"
	SQLCONTENTTYPE__DEFAULT                        = "DEFAULT"
	SQLCONTENTTYPE__CHARSET                        = "CHARSET"
	SQLCONTENTTYPE__COLLATE                        = "COLLATE"
	SQLCONTENTTYPE__AUTO__INCREMENT                = "AUTO_INCREMENT"
	SQLCONTENTTYPE__COMMENT                        = "COMMENT"
	SQLCONTENTTYPE__NOT_NULL                       = "NOT NULL"
	SQLCONTENTTYPE__NULL                           = "NULL"
	SQLCONTENTTYPE__NOT                            = "NOT"
	SQLCONTENTTYPE__CHARACTER_SET                  = "CHARACTER SET"
	SQLCONTENTTYPE__CHARACTER                      = "CHARACTER"
	SQLCONTENTTYPE__SET                            = "SET"
	SQLCONTENTTYPE__UNSIGNED                       = "UNSIGNED"
	SQLCONTENTTYPE__DEF                            = "def"
)

func (i SqlContentType) String() string {
	return string(i)
}

func (i SqlContentType) LowerString() string {
	return strings.ToLower(i.String())
}