package common

import (
	"bytes"
	"errors"
	"go/build"
	"os"
	"strings"
)

func ParseFilePath(isDebug bool) (path CmdFilePath, err error) {
	if isDebug {
		// for test
		os.Setenv("GOFILE", "main.go")
		os.Setenv("GOPACKAGE", "main")
	}

	path.CmdFileName = os.Getenv("GOFILE")
	if path.CmdFileName == "" {
		err = errors.New("fileName empty")
		return
	}
	path.PackageName = os.Getenv("GOPACKAGE")

	if path.PackageName == "" {
		err = errors.New("packageName empty")
		return
	}
	return
}

func GetImportPathName(path string) string {
	pkg, err := build.ImportDir(path, build.FindOnly)
	if err != nil {
		panic(err)
	}
	return pkg.ImportPath
}

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
		} else if letter == '_' || letter == '-' {
			flag = true
		} else if flag {
			flag = false
			dst.WriteString(TransLetterToUpper(letter))
		} else {
			dst.WriteString(string(letter))
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
		} else if letter == '_' || letter == '-' {
			flag = true
		} else if flag {
			flag = false
			dst.WriteString(TransLetterToUpper(letter))
		} else {
			dst.WriteString(string(letter))
		}
	}

	return dst.String()
}

// like transform "ToLowerSnakeCase" to "TO_LOWER_SNAKE_CASE"
func ToUpperSnakeCase(s string) string {
	var dst bytes.Buffer
	for i, letter := range s {
		if IsUpperLetter(letter) {
			if i == 0 {
				dst.WriteString(string(letter))
			} else {
				dst.WriteString("_")
				dst.WriteString(string(letter))
			}
		} else {
			dst.WriteString(TransLetterToUpper(letter))
		}
	}
	return strings.ToUpper(dst.String())
}

// like transform "ToLowerSnakeCase" to "to_lower_snake_case"
func ToLowerSnakeCase(s string) string {
	var dst bytes.Buffer
	for i, letter := range s {
		if IsUpperLetter(letter) {
			if i == 0 {
				dst.WriteString(string(letter))
			} else {
				dst.WriteString("_")
				dst.WriteString(string(letter))
			}
		} else {
			dst.WriteString(string(letter))
		}
	}
	return strings.ToLower(dst.String())
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
