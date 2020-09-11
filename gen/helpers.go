package gen

import (
	"bytes"
	"os"
	"strings"
)

// file/dir is exists or not
func IsExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// path is dir or not
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// get cmd exe path
func GetExeFilePath() (exeFilePath string, packageName string, e error) {
	exeFilePath, e = os.Getwd()
	if e != nil {
		return
	}

	//exeFilePath = "/data/golang/go/src/gotest/gorm_test"
	ss := strings.LastIndex(exeFilePath, "/")
	packageName = exeFilePath[ss+1:]
	return
}

// all letter of string are upper or not
func IsUpperLetterString(dstString string) bool {
	for _, letter := range dstString {
		if !IsUpperLetter(letter) {
			return false
		}
	}

	return true
}

// letter is upper or not
func IsUpperLetter(letter rune) bool {
	if letter >= 'A' && letter <= 'Z' {
		return true
	} else {
		return false
	}
}

func IsLowerLetter(letter rune) bool {
	if letter >= 'a' && letter <= 'z' {
		return true
	} else {
		return false
	}
}

func TransLetterToUpper(letter rune) string {
	if IsLowerLetter(letter) {
		letter -= 'a' - 'A'
	}
	return string(letter)
}

func TransLetterToLower(letter rune) string {
	if IsUpperLetter(letter) {
		letter += 'a' - 'A'
	}
	return string(letter)
}

// like transform "to_lower_snake_case" to "ToLowerSnakeCase"
func ToUpperCamelCase(s string) string {
	var dst bytes.Buffer
	var flag bool
	for index, letter := range s {
		if index == 0 {
			dst.WriteString(TransLetterToUpper(letter))
		} else if letter == '_' {
			flag = true
		} else if flag {
			flag = false
			dst.WriteString(TransLetterToUpper(letter))
		} else {
			dst.WriteString(TransLetterToLower(letter))
		}
	}

	return dst.String()
}

// like transform "to_lower_snake_case" to "toLowerSnakeCase"
func ToLowerCamelCase(s string) string {
	var dst bytes.Buffer
	var flag bool
	for index, letter := range s {
		if index == 0 {
			dst.WriteString(TransLetterToLower(letter))
		} else if letter == '_' {
			flag = true
		} else if flag {
			flag = false
			dst.WriteString(TransLetterToUpper(letter))
		} else {
			dst.WriteString(TransLetterToLower(letter))
		}
	}

	return dst.String()
}

// like transform "ToLowerSnakeCase" to "TO_LOWER_SNAKE_CASE"
func ToUpperSnakeCase(s string) string {
	var dst bytes.Buffer
	for _, letter := range s {
		if IsUpperLetter(letter) {
			dst.WriteString("_")
			dst.WriteString(TransLetterToUpper(letter))
		} else {
			dst.WriteString(TransLetterToUpper(letter))
		}
	}
	return dst.String()
}

// like transform "ToLowerSnakeCase" to "to_lower_snake_case"
func ToLowerSnakeCase(s string) string {
	var dst bytes.Buffer
	for _, letter := range s {
		if IsUpperLetter(letter) {
			dst.WriteString("_")
			dst.WriteString(TransLetterToLower(letter))
		} else {
			dst.WriteString(TransLetterToLower(letter))
		}
	}
	return dst.String()
}

func GetDataBetweenFlag(source string, flag1 string, flag2 string) string {
	s := strings.Index(source, flag1)
	if s == -1 {
		return ""
	}
	e := strings.LastIndex(source, flag2)
	if e == -1 {
		return ""
	}
	if s == e {
		return ""
	}
	return source[s+1 : e]
}
