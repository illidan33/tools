package common

import (
	"errors"
	"fmt"
	"go/build"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

func GetGenEnvironmentValues() (path CmdFilePath, err error) {
	path = CmdFilePath{
		CmdFileName: os.Getenv("GOFILE"),
		PackageName: os.Getenv("GOPACKAGE"),
		SysArch:     os.Getenv("GOARCH"),
		Sys:         os.Getenv("GOOS"),
		CmdLine:     os.Getenv("GOLINE"),
		CmdDir:      "",
	}
	path.CmdDir, err = os.Getwd()
	if err != nil {
		err = fmt.Errorf("Getwd error: %s", err.Error())
		return
	}
	if path.PackageName == "" {
		path.PackageName, _ = GetPackageNameFromDir(path.CmdDir)
	}
	if path.Sys == "" {
		path.Sys = runtime.GOOS
	}
	if path.SysArch == "" {
		path.SysArch = runtime.GOARCH
	}
	return
}

func GetBuildPackageFromDir(dir string) (*build.Package, error) {
	pkg, err := build.ImportDir(dir, build.FindOnly)
	if err != nil {
		return nil, err
	}
	return pkg, nil
}

// response total import path
func GetImportPathFromDir(dir string) string {
	return strings.TrimPrefix(dir, GetGoPath()+"/")
}

func GetImportPathFromFile(filePath string) string {
	return strings.TrimPrefix(path.Dir(filePath), GetGoPath()+"/")
}

// such as tools/xxx, parse to github.com/illidan33/tools/xxx
func GetTotalImportPathFromImport(imPath string) (string, error) {
	dir, err := GetDirFromImport(imPath)
	if err != nil {
		return "", err
	}
	return GetImportPathFromDir(dir), nil
}

// response abs dir path
func GetDirFromImport(importPath string) (string, error) {
	i := strings.Index(importPath, string(filepath.Separator))
	firstDirName := importPath[:i]
	lastDirPath := importPath[i+1:]
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	var file string
	dir := pwd
	for {
		tmp := path.Join(dir, "/", importPath)
		// 当前项目的go mod的module名称跟目录名称不一致的情况
		tmpp := path.Join(dir, "/", lastDirPath)
		if file == firstDirName {
			fmt.Println(tmp)
			return tmp, nil
		} else if IsExists(tmp) {
			fmt.Println(tmp)
			return tmp, nil
		} else if IsExists(tmpp) {
			fmt.Println(tmpp)
			return tmpp, nil
		}
		if dir == "" || dir == "/" {
			break
		}
		file = path.Base(dir)
		dir = path.Dir(dir)
	}
	return "", errors.New("not found dir by import path:" + importPath)
}

func GetPackageNameFromDir(path string) (string, error) {
	if !IsDir(path) {
		path = filepath.Dir(path)
	}
	pkg, err := build.ImportDir(path, build.IgnoreVendor)
	if err != nil {
		return "", err
	}
	if pkg.Name == "" {
		return filepath.Base(path), nil
	}
	return pkg.Name, nil
}

func GetGoPath() string {
	var gopath string
	gopath = os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}
	return gopath + "/src"
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
