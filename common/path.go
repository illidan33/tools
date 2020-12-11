package common

import (
	"errors"
	"fmt"
	"github.com/dave/kerr"
	"go/build"
	"os"
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
		path.PackageName, _ = GetPackageNameFromPath(path.CmdDir)
	}
	if path.Sys == "" {
		path.Sys = runtime.GOOS
	}
	if path.SysArch == "" {
		path.SysArch = runtime.GOARCH
	}
	return
}

func GetBuildPackageFromDir(path string) (*build.Package, error) {
	pkg, err := build.ImportDir(path, build.FindOnly)
	if err != nil {
		return nil, err
	}
	return pkg, nil
}

func GetDirFromPackage(pkg string) (string, error) {
	if pkg[0] != '/' {
		pkg = filepath.Join(GetGoPath(), "src", pkg)
	}
	p, err := build.ImportDir(pkg, build.FindOnly)
	if err != nil {
		return "", err
	}
	return p.Dir, nil
}

func GetPackageNameFromPath(path string) (string, error) {
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

//func GetDirFromPackage(environ []string, packagePath string) (string, error) {
//	exe := exec.Command("go", "list", "-f", "{{.Dir}}", packagePath)
//	exe.Env = environ
//	out, err := exe.CombinedOutput()
//	if err == nil {
//		return strings.TrimSpace(string(out)), nil
//	}
//
//	dir, err := GetDirFromEmptyPackage(packagePath)
//	if err != nil {
//		return "", kerr.Wrap("GXTUPMHETV", err)
//	}
//	return dir, nil
//
//}

func GetDirFromEmptyPackage(path string) (string, error) {
	gopath := GetGoPath()
	gopaths := filepath.SplitList(gopath)
	for _, gopath := range gopaths {
		dir := filepath.Join(gopath, "src", path)
		if s, err := os.Stat(dir); err == nil && s.IsDir() {
			return dir, nil
		}
	}
	return "", errors.New("not found")
}

func GetPackageFromDir(dir string) (string, error) {
	gopath := GetGoPath()
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

//func GetPackagePathFromDir(currentDir string) string {
//	gopath := GetGoPath()
//	gopaths := filepath.SplitList(gopath)
//	for _, gopath := range gopaths {
//		if strings.HasPrefix(currentDir, gopath) {
//			return gopath
//		}
//	}
//	return gopaths[0]
//}
