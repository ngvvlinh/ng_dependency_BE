package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

var flV = flag.Bool("v", false, "verbose")
var flLocal = flag.String("local", "o.o", "local namespace for grouping imports")

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: clean-imports DIR ...\n")
		flag.PrintDefaults()
	}
	flRegexp := flag.String("name", "", "regular expression to filter file name")
	flCheckAlias := flag.Bool("check-alias", false, "also check that all imports from non-stdlib have alias")
	flag.Parse()
	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(2)
	}
	if *flLocal == "" || !strings.Contains(*flLocal, ".") {
		fmt.Fprintf(flag.CommandLine.Output(), "invalid local (should be format example.com): `%v`", *flLocal)
		os.Exit(1)
	}
	localPath = []byte(strings.TrimSuffix(*flLocal, "/") + "/")

	var reFilter *regexp.Regexp
	if *flRegexp != "" {
		re, err := regexp.Compile(*flRegexp)
		if err != nil {
			fmt.Fprintf(flag.CommandLine.Output(), "invalid regexp (%v)", *flRegexp)
			os.Exit(1)
		}
		reFilter = re
	}

	args := flag.Args()
	for i, arg := range args {
		path, err := filepath.EvalSymlinks(arg)
		must(err, "invalid path")
		absPath, err := filepath.Abs(path)
		must(err, "invalid path")
		f, err := os.Stat(absPath)
		must(err, "can not read directory")
		if !f.IsDir() {
			panicf("%v must be directory", absPath)
		}
		args[i] = absPath
	}
	var files []string
	for _, arg := range args {
		files = append(files, walk(arg, reFilter, *flCheckAlias)...)
	}
	if len(files) == 0 {
		return
	}
	var goimportsArgs []string
	goimportsArgs = append(goimportsArgs, "-local", *flLocal, "-w")
	goimportsArgs = append(goimportsArgs, files...)
	output, err := exec.Command("goimports", goimportsArgs...).CombinedOutput()
	must(err, "can not run goimports")
	fmt.Printf("%s", output)

	fmt.Printf("Cleaned %v files:\n", len(files))
	for _, file := range files {
		fmt.Printf("  %v\n", file)
	}
}

func must(err error, msg string, args ...interface{}) {
	if err == nil {
		return
	}
	if msg == "" {
		panicf("%v", err)
	}
	panicf(msg+": %v", append(args, err)...)
}

func panicf(msg string, args ...interface{}) {
	panic(fmt.Sprintf(msg, args...))
}

func walk(dir string, filter *regexp.Regexp, requireAlias bool) (files []string) {
	msgAlias := ""
	if requireAlias {
		msgAlias = " (with alias checking)"
	}
	fmt.Printf("Clean imports%v: %v\n", msgAlias, dir)

	ok := true
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("process %v: %v", path, err)
		}
		if info.IsDir() {
			return nil
		}
		if !includeFile(filter, info.Name()) {
			return nil
		}
		if *flV {
			fmt.Printf("[DEBUG] %v\n", path)
		}
		body, err := ioutil.ReadFile(path)
		must(err, "")
		if requireAlias {
			if err := checkAlias(body); err != nil {
				fmt.Fprintf(os.Stdout, "%v: %v\n", path, err)
				ok = false
			}
		}
		changedBody := cleanImport(body)
		if len(changedBody) == 0 {
			return nil
		}
		files = append(files, path)
		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, info.Mode())
		must(err, "")
		for _, part := range changedBody {
			_, err = f.Write(part)
			must(err, "can not write to file %v", path)
		}
		err = f.Close()
		must(err, "")
		return nil
	})
	must(err, "can not process %v", dir)
	if !ok {
		os.Exit(1)
	}
	return
}

func includeFile(filter *regexp.Regexp, name string) bool {
	if filter == nil {
		return strings.HasSuffix(name, ".go")
	}
	return filter.MatchString(name)
}

var reImport = regexp.MustCompile(`import \(\n(\t"[^)]+")\n\)\n`)
var reAliasLine = regexp.MustCompile(`^\t([A-z0-9_]+)|\. "`)
var reNewline = regexp.MustCompile(`\n\n+`)
var newline = []byte("\n")
var quote = []byte(`"`)
var dot = []byte(".")
var localPath []byte

func extractImportGroup(data []byte) ([]int, []byte) {
	idx := reImport.FindSubmatchIndex(data)
	if len(idx) == 0 {
		return nil, nil
	}
	return idx, data[idx[2]:idx[3]]
}

func cleanImport(data []byte) [][]byte {
	idx, importBytes := extractImportGroup(data)
	lines := bytes.Split(importBytes, newline)
	mode, space := 0, false
	for _, line := range lines {
		m := getMode(line)
		if m != 0 && m <= mode && space {
			return simplifyImports(data, idx)
		}
		if m != 0 {
			mode = m
		}
		space = m == 0
	}
	return nil
}

func checkAlias(data []byte) error {
	_, importBytes := extractImportGroup(data)
	lines := bytes.Split(importBytes, newline)
	for _, line := range lines {
		m := getMode(line)
		if m > 1 && !reAliasLine.Match(line) {
			return fmt.Errorf("import should have alias (%s)", bytes.TrimSpace(line))
		}
	}
	return nil
}

func getMode(line []byte) int {
	switch {
	case bytes.Contains(line, localPath): // "o.o/..."
		return 3

	case bytes.Contains(line, dot): // "example.com/..."
		return 2

	case bytes.Contains(line, quote): // "fmt"
		return 1

	default:
		return 0
	}
}

func simplifyImports(data []byte, idx []int) [][]byte {
	importBytes := data[idx[0]:idx[1]]
	simplifiedBytes := reNewline.ReplaceAll(importBytes, newline)
	return [][]byte{
		data[:idx[0]],
		simplifiedBytes,
		data[idx[1]:],
	}
}
