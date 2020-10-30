package common

import (
	"errors"
	"fmt"
	"github.com/dave/kerr"
	"go/build"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func GetGenEnvironmentValues(isDebug bool) (path CmdFilePath, err error) {
	if isDebug {
		// for test
		os.Setenv("GOFILE", "main.go")
		os.Setenv("GOPACKAGE", "main")
		os.Setenv("GOARCH", "amd64")
		os.Setenv("GOOS", "macos")
		os.Setenv("GOLINE", "5")
	}
	path = CmdFilePath{
		CmdFileName: os.Getenv("GOFILE"),
		PackageName: os.Getenv("GOPACKAGE"),
		SysArch:     os.Getenv("GOARCH"),
		Sys:         os.Getenv("GOOS"),
		CmdLine:     os.Getenv("GOLINE"),
		CmdDir:      "",
	}
	path.CmdDir, err = os.Getwd()
	fmt.Printf("%#v\n", path)
	return
}

func GetCurrentDirectory() (string, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "", err
	}
	return dir, nil
}

func GetImportPathName(path string) string {
	pkg, err := build.ImportDir(path, build.FindOnly)
	if err != nil {
		panic(err)
	}
	return pkg.ImportPath
}

func GetGoPath() string {
	var gopath string
	gopath = os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}
	return gopath
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

func GetDirFromPackage(environ []string, gopath string, packagePath string) (string, error) {
	exe := exec.Command("go", "list", "-f", "{{.Dir}}", packagePath)
	exe.Env = environ
	out, err := exe.CombinedOutput()
	if err == nil {
		return strings.TrimSpace(string(out)), nil
	}

	dir, err := GetDirFromEmptyPackage(gopath, packagePath)
	if err != nil {
		return "", kerr.Wrap("GXTUPMHETV", err)
	}
	return dir, nil

}

func GetDirFromEmptyPackage(gopathEnv string, path string) (string, error) {
	gopaths := filepath.SplitList(gopathEnv)
	for _, gopath := range gopaths {
		dir := filepath.Join(gopath, "src", path)
		if s, err := os.Stat(dir); err == nil && s.IsDir() {
			return dir, nil
		}
	}
	return "", errors.New("not found")
}

func GetPackageFromDir(gopath string, dir string) (string, error) {
	gopaths := filepath.SplitList(gopath)
	var savedError error
	for _, gopath := range gopaths {
		if strings.HasPrefix(dir, gopath) {
			gosrc := fmt.Sprintf("%s/src", gopath)
			relpath, err := filepath.Rel(gosrc, dir)
			if err != nil {
				// notest
				// I don't *think* we can trigger this error if dir starts with gopath
				savedError = err
				continue
			}
			if relpath == "" {
				// notest
				// I don't *think* we can trigger this either
				continue
			}
			// Remember we're returning a package path which uses forward slashes even on windows
			return filepath.ToSlash(relpath), nil
		}
	}
	if savedError != nil {
		// notest
		return "", savedError
	}
	return "", kerr.New("CXOETFPTGM", "Package not found for %s", dir)
}

func GetCurrentGopath(gopath string, currentDir string) string {
	gopaths := filepath.SplitList(gopath)
	for _, gopath := range gopaths {
		if strings.HasPrefix(currentDir, gopath) {
			return gopath
		}
	}
	return gopaths[0]
}
