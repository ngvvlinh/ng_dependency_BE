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

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: clean-imports DIR ...\n")
		flag.PrintDefaults()
	}
	flag.Parse()
	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(2)
	}

	for _, arg := range flag.Args() {
		absPath, err := filepath.Abs(arg)
		if err != nil {
			must(err, "invalid path")
		}
		f, err := os.Stat(absPath)
		if err != nil {
			must(err, "can not read directory")
		}
		if !f.IsDir() {
			panicf("%v must be directory", absPath)
		}
	}
	var files []string
	for _, arg := range flag.Args() {
		files = append(files, walk(arg)...)
	}
	if len(files) == 0 {
		return
	}
	var goimportsArgs []string
	goimportsArgs = append(goimportsArgs, "-local", "etop.vn", "-w")
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

func walk(dir string) (files []string) {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("process %v: %v", path, err)
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(info.Name(), ".go") {
			return nil
		}
		body, err := ioutil.ReadFile(path)
		must(err, "")
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
	return
}

var reImport = regexp.MustCompile(`import \(\n(\t"[^)]+")\n\)\n`)
var reNewline = regexp.MustCompile(`\n\n+`)
var newline = []byte("\n")
var quote = []byte(`"`)
var etopvn = []byte("etop.vn/")
var dot = []byte(".")

func cleanImport(data []byte) [][]byte {
	idx := reImport.FindSubmatchIndex(data)
	if len(idx) == 0 {
		return nil
	}
	importBytes := data[idx[2]:idx[3]]
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

func getMode(line []byte) int {
	switch {
	case bytes.Contains(line, etopvn): // "etop.vn/..."
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
