// deps is a simple tool for visualizing dependencies of a go package
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"golang.org/x/tools/go/packages"
)

var (
	flStdlib               = flag.Bool("stdlib", false, "include stdlib")
	flExternal             = flag.Bool("external", false, "include external libs")
	flFlat                 = flag.Bool("flat", false, "flatten dependencies")
	flParent               = flag.Bool("parent", false, "include parent on the right")
	flPrintFiles           = flag.Bool("print-files", false, "print list of source files")
	flCopyFiles            = flag.String("copy-files", "", "copy source code to directory")
	flTags                 = flag.String("tags", "release", "Go build tags (comma separated)")
	flRemoveComments       = flag.Bool("remove-comments", false, "Remove comments")
	flRenameGeneratedFiles = flag.Bool("rename-generated-fields", false, "Rename generated files")
)

func usage() {
	const text = `
Usage: deps [OPTION] PACKAGE_PATTERN
`
	fmt.Print(text[1:])
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()
	patterns := flag.Args()
	if len(patterns) != 1 {
		flag.Usage()
		os.Exit(2)
	}

	cfg := &packages.Config{
		Mode: packages.NeedImports | packages.NeedDeps |
			packages.NeedFiles | packages.NeedCompiledGoFiles,
	}
	if len(*flTags) != 0 {
		cfg.BuildFlags = []string{"-tags", *flTags}
	}
	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		panic(err)
	}
	switch len(pkgs) {
	case 0:
		fmt.Println("no package")
		os.Exit(1)
	case 1:
		// continue
	default:
		fmt.Println("only single package is supported")
		os.Exit(1)
	}

	outs := make([]string, 0, 1024)
	printDeps(&outs, ".", pkgs[0].Imports, "")

	if *flFlat {
		sort.Strings(outs)
	}
	for _, line := range outs {
		fmt.Println(line)
	}

	if *flPrintFiles || *flCopyFiles != "" {
		copySource(*flCopyFiles, pkgs[0], printedPkg)
	}
}

var printedPkg = map[string]*packages.Package{}

func printDeps(outs *[]string, parentPath string, pkgs map[string]*packages.Package, level string) {
	list := make([]string, 0, len(pkgs))
	for path := range pkgs {
		list = append(list, path)
	}
	sort.Strings(list)
	for _, path := range list {
		// ignore stdlib
		if !*flStdlib && !strings.Contains(path, ".") {
			continue
		}
		// ignore external libs
		if !*flExternal && !strings.HasPrefix(path, "o.o") {
			// beware if external libs import o.o/
			// if !printedPkg[path] {
			// 	printDeps(pkgs[path].Imports, level+"  .")
			// }
			continue
		}
		if printedPkg[path] != nil {
			if !*flFlat {
				output(outs, "[+] %v%v", level, path)
			}
			continue
		}
		printedPkg[path] = pkgs[path]
		if *flFlat {
			output(outs, "")
		} else {
			output(outs, "    ")
		}
		outputA(outs, "%v%v", level, path)
		if *flParent {
			outputA(outs, "\t<- %v", parentPath)
		}

		if *flFlat {
			printDeps(outs, path, pkgs[path].Imports, level)
		} else {
			printDeps(outs, path, pkgs[path].Imports, level+" . ")
		}
	}
}

func output(outs *[]string, format string, args ...interface{}) {
	*outs = append(*outs, fmt.Sprintf(format, args...))
}

func outputA(outs *[]string, format string, args ...interface{}) {
	_outs := *outs
	last := _outs[len(_outs)-1]
	last += fmt.Sprintf(format, args...)
	_outs[len(_outs)-1] = last
}

func copySource(dstPath string, rootPkg *packages.Package, pkgs map[string]*packages.Package) {
	if dstPath != "" {
		var err error
		dstPath, err = filepath.Abs(dstPath)
		if err != nil {
			panic(err)
		}
		stat, err := os.Stat(dstPath)
		if err != nil {
			panic(err)
		}
		if !stat.IsDir() {
			panic(fmt.Sprintf("%v is not a directory", dstPath))
		}
	}

	list := make([]string, 0, 1024)
	for path := range pkgs {
		if !strings.HasPrefix(path, "o.o") {
			continue
		}
		list = append(list, path)
	}
	sort.Strings(list)

	fmt.Print("\n--- files ---\n\n")
	copyPkg(dstPath, rootPkg)
	for _, path := range list {
		pkg := pkgs[path]
		copyPkg(dstPath, pkg)
	}
}

func copyPkg(dstPath string, pkg *packages.Package) {
	for _, file := range pkg.CompiledGoFiles {
		if *flPrintFiles {
			fmt.Println(file)
		}
		if dstPath != "" {
			idx := strings.LastIndex(file, "/o.o/")
			dstFile := filepath.Join(dstPath, file[idx+1:])
			if err := copyFile(file, dstFile); err != nil {
				panic(fmt.Sprintf("copy %v to %v: %v", file, dstFile, err))
			}
		}
	}
}

func copyFile(srcPath, dstPath string) error {
	body, err := ioutil.ReadFile(srcPath)
	if err != nil {
		return err
	}

	dstDir := filepath.Dir(dstPath)
	dstName := filepath.Base(dstPath)
	err = os.MkdirAll(dstDir, 0755)
	if err != nil {
		return err
	}
	body = cleanup(body)

	return ioutil.WriteFile(filepath.Join(dstDir, rename(dstName)), body, 0755)
}

func rename(name string) string {
	if !*flRenameGeneratedFiles {
		return name
	}
	name = strings.Replace(name, "zz_generated.", "g", 1)
	name = strings.Replace(name, "zz_release.", "r", 1)
	name = strings.Replace(name, "_gen.go", "g.go", 1)
	name = strings.Replace(name, ".gen.go", "g.go", 1)
	return name
}

var reCmt1 = regexp.MustCompile(`(^|\s)/\*([^*]|(\*[^/])|\n)*\*/`)
var reCmt2 = regexp.MustCompile(`(^|\s)//.*`)
var reNL = regexp.MustCompile(`\s*\n\s*\n(\s*\n)+`)
var rePR = regexp.MustCompile(`{\s*\n(\s*\n)+`)

// remove comments
func cleanup(body []byte) []byte {
	if !*flRemoveComments {
		return body
	}
	body = reCmt1.ReplaceAll(body, []byte{'\n'})
	body = reCmt2.ReplaceAll(body, []byte{'\n'})
	body = reNL.ReplaceAll(body, []byte{'\n', '\n'})
	body = rePR.ReplaceAll(body, []byte{'{', '\n'})
	return body
}
