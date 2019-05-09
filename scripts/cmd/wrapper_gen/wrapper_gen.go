package main

import (
	"bytes"
	"flag"
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
	wrapper_gen [FLAGS] FILENAME ...

FILENAME : Protobuf file
FLAGS	 :
    -p   : Package prefix (for conflicting package name: etop/shop and ext/shop)
	-s	 : Strip prefix
	-o	 : Output directory
`

var (
	flPkgPrefix   = flag.String("p", "", "Package prefix (for conflicting package name: etop/shop and ext/shop)")
	flStripPrefix = flag.String("s", "", "Strip prefix")
	flOutputDir   = flag.String("o", "wrapper", "Output directory")
)

func main() {
	flag.Usage = func() {
		fmt.Println(usage)
	}

	flag.Parse()
	filenames := flag.Args()
	if len(filenames) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	for _, filename := range filenames {
		generateWrapper(filename)
	}
}

func generateWrapper(filename string) {
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
		s.PkgPrefix = *flPkgPrefix
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

	filename = strings.TrimPrefix(filename, *flStripPrefix)
	genFilename := filepath.Join(
		*flOutputDir,
		filepath.Dir(filename),
		strings.Split(filepath.Base(filename), ".")[0]+".gen.go")
	gen.WriteFile(genFilename, buf.Bytes())
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
