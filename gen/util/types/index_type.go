package types

import "strings"

//go:generate stringer -type=IndexType
type IndexType uint8

const (
	INDEX_TYPE__UNKNOWN      IndexType = iota // 未知
	INDEX_TYPE__PRIMARY                       // 主键
	INDEX_TYPE__UNIQUE_INDEX                  // 唯一索引
	INDEX_TYPE__INDEX                         // 普通索引
)

func (i IndexType) KeyString() string {
	arr := strings.Split(i.String(), "__")
	return arr[1]
}
func (i IndexType) KeyLowerString() string {
	arr := strings.Split(i.String(), "__")
	return strings.ToLower(arr[1])
}
