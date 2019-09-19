package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func usage() {
	fmt.Print(`Usage:
    twirp-modifier FILE ...

`)
	flag.PrintDefaults()
}

func fatalf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
	os.Exit(1)
}

func main() {
	flag.Usage = usage
	flag.Parse()
	files := flag.Args()
	if len(files) == 0 {
		usage()
		os.Exit(2)
	}

	for _, file := range files {
		if !strings.HasSuffix(file, ".twirp.go") {
			continue
		}
		body, err := ioutil.ReadFile(file)
		if err != nil {
			fatalf("%v\n", err)
		}
		f, err := os.Create(file)
		if err != nil {
			fatalf("%v\n", err)
		}
		if err = cleanImports(f, body); err != nil {
			fatalf("processing %v: %v\n", file, err)
		}
	}
}

var txtImport0 = []byte("\nimport ")
var txtImport = []byte("import ")

func cleanImports(w io.Writer, body []byte) error {
	idx0 := bytes.Index(body, txtImport0)
	if idx0 < 0 {
		return fmt.Errorf("import not found (is it a twirp file?)")
	}
	idx1 := bytes.LastIndex(body, txtImport0)
	idx2 := bytes.IndexByte(body[idx1+1:], '\n')
	if idx1 < 0 || idx2 < 0 {
		panic("unexpected")
	}
	idx3 := idx1 + 1 + idx2

	importBytes := body[idx0:idx3]
	lines := bytes.Split(importBytes, []byte("\n"))
	result := make([]byte, 0, len(importBytes))
	result = append(result, "import (\n"...)
	for _, line := range lines {
		if bytes.HasPrefix(line, txtImport) {
			result = append(result, line[len(txtImport):]...)
			result = append(result, '\n')
		}
	}
	result = append(result, ")\n"...)

	if _, err := w.Write(body[:idx0]); err != nil {
		return err
	}
	if _, err := w.Write(result); err != nil {
		return err
	}
	if _, err := w.Write(body[idx3:]); err != nil {
		return err
	}
	return nil
}
