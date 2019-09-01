package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

var flV = flag.Bool("v", false, "verbose")
var flPb = flag.String("ext-pb", ".pb.go", "extension of generated protobuf Go files")
var flOut = flag.String("format-out", ".d.go", "extension of generated declaration Go files")

func usage() {
	const text = `
split-pbgo-def finds all type definitions in .pb.go files and move them to
.def.go for easily scan

Usage: split-pbgo-def [OPTION] FILE ...

Options:
`

	fmt.Print(text[1:])
	flag.PrintDefaults()
}

func init() {
	flag.Usage = usage
}

func debugf(format string, args ...interface{}) {
	if *flV {
		fmt.Printf(format, args...)
	}
}

func fatalf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
	os.Exit(1)
}

func must(err error) {
	if err != nil {
		fatalf("%v\n", err)
	}
}

func p(w io.Writer, format string, args ...interface{}) {
	_, err := fmt.Fprintf(w, format, args...)
	must(err)
}

func main() {
	flag.Parse()
	extPb := *flPb
	extOut := *flOut
	if extPb == "" || extOut == "" {
		fatalf("invalid extension")
	}

	files := flag.Args()
	if len(files) == 0 {
		flag.Usage()
		os.Exit(2)
	}

	var changedFiles []string
	for _, pbFile := range files {
		if !strings.HasSuffix(pbFile, extPb) {
			fmt.Printf("ignored %v\n", pbFile)
			continue
		}
		debugf("processing %v\n", pbFile)

		body, err := ioutil.ReadFile(pbFile)
		if err != nil {
			fatalf("%v\n", err)
		}

		pbBody, defBody, err := process(body)
		if err != nil {
			fatalf("error while processing %v: %v\n", pbFile, err)
		}

		defFile := strings.TrimSuffix(pbFile, extPb) + extOut
		must(ioutil.WriteFile(defFile, defBody, 0666))
		must(ioutil.WriteFile(pbFile, pbBody, 666))

		changedFiles = append(changedFiles, pbFile, defFile)
		fmt.Printf("generated %v\n", defFile)
		fmt.Printf("  updated %v\n", pbFile)
	}
	execGoimport(changedFiles)
}

var (
	rePackage     = regexp.MustCompile(`package [_0-9A-z]+`)
	reImport      = regexp.MustCompile(`^import .*$`)
	reImportGroup = regexp.MustCompile(`import \(([^)]+\n)+\)\n`)
	reTypeDef     = regexp.MustCompile(`(\\\\.*\n)*type [^\n]*([^}]*\n)*}\n`)
	reOneOf       = regexp.MustCompile(`type ([0-9A-z]+)_[_0-9A-z]+ struct {\n.+protobuf:"bytes,([-0-9]+),opt,name=([_0-9A-z]+),oneof"`)
)

func process(input []byte) (pbBody []byte, defBody []byte, err error) {
	pkg := rePackage.Find(input)
	if pkg == nil {
		return nil, nil, errors.New("can not find package definition")
	}
	imports, err := findImports(input)
	if err != nil {
		return nil, nil, err
	}

	pbBuf := &bytes.Buffer{}
	defBuf := &bytes.Buffer{}

	p(defBuf, "// Code generated by split-pbgo-def. DO NOT EDIT.\n\n")
	p(defBuf, "%s\n\n", pkg)

	imports = cleanImports(imports)
	if len(imports) > 0 {
		p(defBuf, "import (\n")
		for _, imp := range imports {
			p(defBuf, "\t%s\n", imp)
		}
		p(defBuf, ")\n\n")
	}

	oneOfs := findOneOfs(input)
	if len(oneOfs) > 0 {
		lastGroup := ""
		for _, item := range oneOfs {
			if item.Group != lastGroup {
				lastGroup = item.Group
				p(defBuf, "type %vEnum int32\n", item.Group)
			}
		}

		p(defBuf, "\nconst (\n")
		for _, item := range oneOfs {
			p(defBuf, "\t%v %vEnum = %v\n", item.Name, item.Group, item.Tag)
		}
		p(defBuf, ")\n\n")

		for _, item := range oneOfs {
			p(defBuf, "func (_ %v_%v) GetEnumTag() %vEnum { return %v }\n",
				item.Group, item.Name, item.Group, item.Tag)
		}
	}

	indexes := findTypesIndex(input)
	lastIdx := 0
	for _, idx := range indexes {
		pbBuf.Write(input[lastIdx:idx[0]])
		defBuf.Write(input[idx[0]:idx[1]])
		defBuf.Write([]byte("\n"))

		lastIdx = idx[1]
	}
	pbBuf.Write(input[lastIdx:])
	return pbBuf.Bytes(), defBuf.Bytes(), nil
}

func findImports(input []byte) (result [][]byte, err error) {
	imports := reImport.FindAll(input, -1)
	for _, imp := range imports {
		result = append(result, imp[len("import "):])
	}

	groups := reImportGroup.FindAllSubmatch(input, -1)
	for _, group := range groups {
		imports := bytes.Split(group[1], []byte("\n"))
		for _, imp := range imports {
			imp = bytes.TrimSpace(imp)
			if len(imp) == 0 {
				continue
			}
			result = append(result, imp)
		}
	}

	if len(result) == 0 {
		err = errors.New("can not find import definition")
	}
	return
}

func cleanImports(imports [][]byte) [][]byte {
	result := make([][]byte, 0, len(imports))
	for _, imp := range imports {
		if bytes.HasPrefix(imp, []byte("_ ")) {
			continue
		}
		result = append(result, imp)
	}
	return result
}

var indexVar = []byte("\nvar ")

func findTypesIndex(input []byte) [][]int {
	indexes := reTypeDef.FindAllIndex(input, -1)

	// handle special case with enum
	for _, idx := range indexes {
		block := input[idx[0]:idx[1]]
		varIdx := bytes.Index(block, indexVar)
		if varIdx >= 0 {
			idx[1] = idx[0] + varIdx
		}
	}
	return indexes
}

type OneOfItem struct {
	Group string
	Name  string
	Tag   int
}

func findOneOfs(input []byte) []OneOfItem {
	matches := reOneOf.FindAllSubmatch(input, -1)
	result := make([]OneOfItem, len(matches))
	for i, match := range matches {
		item := OneOfItem{
			Group: string(match[1]),
			Name:  string(match[3]),
			Tag:   mustParseInt(string(match[2])),
		}
		result[i] = item
	}
	return result
}

func mustParseInt(s string) int {
	i, err := strconv.Atoi(s)
	must(err)
	return i
}

func execGoimport(files []string) {
	args := []string{"-w"}
	args = append(args, files...)
	cmd := exec.Command("goimports", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fatalf("%s\n\n%s\n", err, out)
	}
	debugf("%s", out)
}