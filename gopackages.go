package gopackages

// ke: {"package": {"complete": true}}

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"os"

	"github.com/davelondon/kerr"
)

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
	return "", NotFoundError{Struct: kerr.New("SUTCWEVRXS", "%s not found", path)}
}

type NotFoundError struct {
	kerr.Struct
}

func GetPackageFromDir(gopath string, dir string) (string, error) {
	gopaths := filepath.SplitList(gopath)
	var savedError error
	for _, gopath := range gopaths {
		if strings.HasPrefix(dir, gopath) {
			gosrc := fmt.Sprintf("%s/src", gopath)
			relpath, err := filepath.Rel(gosrc, dir)
			if err != nil {
				// ke: {"block": {"notest": true}}
				// I don't *think* we can trigger this error if dir starts with gopath
				savedError = err
				continue
			}
			if relpath == "" {
				// ke: {"block": {"notest": true}}
				// I don't *think* we can trigger this either
				continue
			}
			// Remember we're returning a package path which uses forward slashes even on windows
			return filepath.ToSlash(relpath), nil
		}
	}
	if savedError != nil {
		// ke: {"block": {"notest": true}}
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
