package projectpath

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var (
	mainPath string
	goPath   string
)

func GetPath() string {
	return mainPath
}

func GetGoPath() string {
	return goPath
}

func init() {
	var err error
	goPath, err = loadGoPath()
	if err != nil {
		panic(fmt.Sprintf("can not load gopath: %v", err))
	}
	mainPath, err = loadPath()
	if err != nil {
		panic(fmt.Sprintf("can not load path: %v", err))
	}
}

func loadPath() (string, error) {
	dir := os.Getenv("PROJECT_DIR")
	if dir == "" {
		return os.Getwd()
	}
	return loadProjectPath(dir)
}

func loadGoPath() (path string, err error) {
	gopathEnv := os.Getenv("GOPATH")
	if gopathEnv == "" {
		homeEnv := os.Getenv("HOME")
		if homeEnv == "" {
			return "", errors.New("HOME not set")
		}
		path = filepath.Join(homeEnv, "go")

	} else {
		ps := filepath.SplitList(gopathEnv)
		path = ps[0]
	}

	path, err = filepath.Abs(path)
	if err != nil {
		return "", err
	}
	return path, validDir(path)
}

func loadProjectPath(dir string) (string, error) {
	path := filepath.Join(dir, "backend")
	path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	return path, validDir(path)
}

func validDir(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	mod := info.Mode()
	if mod&os.ModeDir == 0 && mod&os.ModeSymlink == 0 {
		return fmt.Errorf(`"%v" is not a directory`, path)
	}
	return nil
}
