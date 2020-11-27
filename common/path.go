package common

import (
	"errors"
	"fmt"
	"github.com/dave/kerr"
	"go/build"
	"os"
	"os/exec"
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
		path.PackageName, _ = GetImportPackageName(path.CmdDir)
		if path.PackageName == "" {
			path.PackageName = filepath.Base(path.CmdDir)
		}
	}
	if path.Sys == "" {
		path.Sys = runtime.GOOS
	}
	if path.SysArch == "" {
		path.SysArch = runtime.GOARCH
	}
	return
}

func GetCurrentDirectory() (string, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "", err
	}
	return dir, nil
}

func GetImportPath(path string) (string, error) {
	pkg, err := build.ImportDir(path, build.FindOnly)
	if err != nil {
		return "", err
	}
	return pkg.ImportPath, nil
}

func GetImportByPackage(pkg string) (string, error) {
	p, err := build.ImportDir(filepath.Join(GetGoPath(), "src", pkg), build.FindOnly)
	if err != nil {
		return "", err
	}
	return p.Dir, nil
}

func GetImportPackageName(path string) (string, error) {
	if strings.HasSuffix(path, ".go") {
		path = filepath.Dir(path)
	}
	pkg, err := build.ImportDir(path, build.ImportComment)
	if err != nil {
		return "", err
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
