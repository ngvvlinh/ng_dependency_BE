package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"etop.vn/backend/up/gogen/pkg/gen"
	"etop.vn/backend/up/gogen/pkg/grpcgen"
)

const usage = `
Usage:
    wrapper_gen FLAGS FILENAME ... [FLAGS FILENAME] ...

FILENAME : Protobuf file
FLAGS    :
    -p   : Package prefix (for conflicting package name: etop/shop and ext/shop)
    -s   : Strip prefix
    -o   : Output directory

Example:
    wrapper_gen -p shop shop/one.go shop/two.go -p main main/a.go
`

type Config struct {
	PkgPrefix   string
	StripPrefix string
	OutputDir   string
}

type InputGroup struct {
	Config
	Filenames []string
}

type InputGroups []InputGroup

func main() {
	groups := parseFlags()

	var outputFilenames []string
	for _, group := range groups {
		for _, filename := range group.Filenames {
			genFilename := generateWrapper(group.Config, filename)
			outputFilenames = append(outputFilenames, genFilename)
		}
	}
	gen.FormatFiles(outputFilenames...)
}

func parseFlags() (groups InputGroups) {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println(usage)
		os.Exit(2)
	}

	var current InputGroup
	var lastFile bool

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "" {
			continue
		}
		if arg[0] == '-' && lastFile {
			// start a new group
			groups = append(groups, current)
			current = InputGroup{}
		}
		lastFile = arg[0] != '-'
		if arg[0] == '-' {
			switch arg {
			case "-p":
				i++
				mustReadParamForFlag(&current.PkgPrefix, arg, args, i)

			case "-s":
				i++
				mustReadParamForFlag(&current.StripPrefix, arg, args, i)

			case "-o":
				i++
				mustReadParamForFlag(&current.OutputDir, arg, args, i)

			default:
				fatalf("unknown flag %v", arg)
			}
		} else {
			current.Filenames = append(current.Filenames, arg)
		}
	}
	groups = append(groups, current)
	return groups
}

func fatalf(msg string, args ...interface{}) {
	fmt.Printf(msg+"\n", args...)
	os.Exit(2)
}

func mustReadParamForFlag(out *string, flag string, args []string, i int) {
	if *out != "" {
		fatalf("duplicated flag %v", flag)
	}
	result, err := readParamForFlag(flag, args, i)
	if err != nil {
		fatalf(err.Error())
	}
	*out = result
}

func readParamForFlag(flag string, args []string, i int) (string, error) {
	if i >= len(args) {
		return "", fmt.Errorf("no argument after flag %v", flag)
	}
	arg := args[i]
	if arg == "" || arg[0] == '-' {
		return "", fmt.Errorf("no argument after flag %v", flag)
	}
	return args[i], nil
}

func generateWrapper(cfg Config, filename string) (outputFilename string) {
	if cfg.OutputDir == "" {
		fatalf("no output directory")
	}

	result := grpcgen.ParseServiceFile(filename, grpcgen.Options{
		ImportCurrentPackage: true,
		IncludeInterface: func(name string) bool {
			return name != "HTTPClient" && name != "TwirpServer"
		},
	})
	for _, s := range result.Services {
		if !strings.HasSuffix(s.Name, "Service") {
			gen.Fatalf("Service %v must have suffix `Service` (expected %vService)", s.Name, s.Name)
		}
		s.PkgPrefix = cfg.PkgPrefix
	}

	var buf bytes.Buffer
	tpl := template.Must(template.New("tpl").Funcs(funcs).Parse(tplStr))
	err := tpl.Execute(&buf, map[string]interface{}{
		"PackageName": result.Package,
		"PackagePath": result.PackagePath,
		"ServiceName": toPascalCase(result.Package),
		"Imports":     result.ExtraImports(),
		"Services":    result.Services,
		"HasSecret":   hasSecret(result.Package),
	})
	gen.NoError(err, "Unable to execute template")

	filename, err = filepath.Rel(gen.ProjectPath(), filename)
	gen.NoError(err, "Unable to get relative path")

	filename = strings.TrimPrefix(filename, cfg.StripPrefix)
	genFilename := filepath.Join(
		cfg.OutputDir,
		filepath.Dir(filename),
		strings.Split(filepath.Base(filename), ".")[0]+".gen.go")
	gen.WriteFile(genFilename, buf.Bytes())
	return genFilename
}

func toPascalCase(name string) string {
	return strings.ToUpper(name[:1]) + name[1:]
}

func hasSecret(pkgName string) bool {
	prefix := pkgName + "."
	for k, perm := range ACL {
		if strings.HasPrefix(k, prefix) && perm.Type == Secret {
			return true
		}
	}
	return false
}
