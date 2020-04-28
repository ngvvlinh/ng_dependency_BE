package projectpath

import (
	"fmt"
	"os"
	"path/filepath"
)

var mainPath string

func GetPath() string {
	return mainPath
}

func init() {
	var err error
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
