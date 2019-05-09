package gen

import (
	"errors"
	"fmt"
	"go/importer"
	"go/types"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

const ProjectImport = "etop.vn/backend"
const ProjectImport_ = ProjectImport + "/"

var (
	gopath      string
	projectPath string
	rootTree    string
)

func init() {
	var err error
	gopath, err = loadGoPath()
	if err != nil {
		Fatalf("can not load gopath: %v", err)
	}
	projectPath, err = loadProjectPath()
	if err != nil {
		Fatalf("can not load project path: %v", err)
	}
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

func loadProjectPath() (path string, err error) {
	dirEnv := os.Getenv("ETOPDIR")
	if dirEnv == "" {
		return "", errors.New("ETOPDIR not set")
	}

	rootTree, err = filepath.Abs(filepath.Join(dirEnv, ".."))
	if err != nil {
		return "", err
	}

	path = filepath.Join(dirEnv, "backend")
	path, err = filepath.Abs(path)
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

// GoPath ...
func GoPath() string {
	return gopath
}

// ProjectPath is absolute path to etop.vn/backend
func ProjectPath() string {
	return projectPath
}

// RootTree is absolute path to etop.vn/..
func RootTree() string {
	return rootTree
}

// WriteFile ...
func WriteFile(outputPath string, data []byte) {
	absPath := GetAbsPath(outputPath)
	err := os.MkdirAll(filepath.Dir(absPath), os.ModePerm)
	if err != nil {
		Fatalf("Unable to create output directory: %v", err)
	}

	err = ioutil.WriteFile(absPath, data, os.ModePerm)
	if err != nil {
		fmt.Printf("Unable to write file `%v`.\n  Error: %v\n", absPath, err)
		os.Exit(1)
	}
	fmt.Println("Generated file:", absPath)

	FormatFile(outputPath)
}

// FormatFile ...
func FormatFile(outputPath string) {
	absPath := GetAbsPath(outputPath)
	out, err := exec.Command("goimports", "-w", absPath).Output()
	if err != nil {
		Fatalf("Unable to run `gofmt -w %v`.\n  Error: %v\n", absPath, err)
	}
	fmt.Print(string(out))
}

// GetAbsPath ...
func GetAbsPath(inputPath string) string {
	if strings.HasPrefix(inputPath, "/") {
		return inputPath
	}
	inputPath = strings.TrimPrefix(inputPath, ProjectImport)
	return filepath.Join(projectPath, inputPath)
}

func GetModPath(pkgPath string) (string, error) {
	if pkgPath == ProjectImport {
		return projectPath, nil
	}
	if strings.HasPrefix(pkgPath, ProjectImport_) {
		return filepath.Join(projectPath, pkgPath[len(ProjectImport_):]), nil
	}
	cmd := exec.Command("go", "list", "-m", pkgPath)
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("can not get abs path: %v", err)
	}
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("can not get abs path: %v", err)
	}
	return strings.Replace(string(out), " ", "@", 1), nil
}

// NoError ...
func NoError(err error, msg string, args ...interface{}) {
	if err != nil {
		fmt.Printf(msg, args...)
		fmt.Printf("\n  Error: %v\n", err)
		os.Exit(1)
	}
}

// Fatalf ...
func Fatalf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
	fmt.Println()
	os.Exit(1)
}

type Importer struct {
	m types.ImporterFrom
}

func NewImporter() *Importer {
	return &Importer{importer.Default().(types.ImporterFrom)}
}

func (imp *Importer) Import(path string) (*types.Package, error) {
	pkg, err := imp.m.Import(path)
	return pkg, err
}

var reQuote = regexp.MustCompile(`"[^"]+"`)

func (imp *Importer) ImportFrom(path, srcDir string, mode types.ImportMode) (*types.Package, error) {
	pkg, err := imp.m.ImportFrom(path, srcDir, mode)
	if err != nil && TryFixingImportError(err) {
		// Import it again
		pkg, err = imp.m.ImportFrom(path, srcDir, mode)
	}
	return pkg, err
}

func TryFixingImportError(err error) bool {
	if !strings.Contains(err.Error(), "can't find import") {
		return false
	}
	s := reQuote.FindString(err.Error())
	pkgPath := s[1 : len(s)-1] // remove the quote

	// Build the imported package
	cmd := exec.Command("go", "install", pkgPath)
	output, _ := cmd.CombinedOutput()
	fmt.Printf("Build package: %v\n%s", pkgPath, output)
	return true
}
