package enum

import (
	"fmt"
	"go/types"
	"sort"
	"text/template"

	"golang.org/x/tools/go/packages"

	"etop.vn/backend/tools/pkg/generator"
	"etop.vn/backend/tools/pkg/generators/api/defs"
	"etop.vn/backend/tools/pkg/generators/api/parse"
	"etop.vn/backend/tools/pkg/genutil"
)

type plugin struct {
	generator.Qualifier
}

func New() generator.Plugin {
	return &plugin{Qualifier: genutil.Qualifier{}}
}

func (p *plugin) Name() string { return "enum" }

func (p *plugin) Filter(ng generator.FilterEngine) error {
	for _, pkg := range ng.ParsingPackages() {
		ds := pkg.InlineDirectives
		if _, ok := ds.Get("enum"); ok {
			pkg.Include()
		}
	}
	return nil
}

func (p *plugin) Generate(ng generator.Engine) error {
	return ng.GenerateEachPackage(p.generatePackage)
}

func (p *plugin) generatePackage(ng generator.Engine, pkg *packages.Package, printer generator.Printer) error {
	mapEnum, err := parse.ParseEnumInPackage(ng, pkg)
	if err != nil {
		return err
	}
	enums := make([]*defs.Enum, 0, len(mapEnum))
	for _, enum := range mapEnum {
		enums = append(enums, enum)
	}
	sort.Slice(enums, func(i, j int) bool { return enums[i].Name < enums[j].Name })

	printer.Import("fmt", "fmt")
	printer.Import("encode", "etop.vn/capi/encode")
	vars := map[string]interface{}{
		"Enums": enums,
	}
	return tpl.Execute(printer, vars)
}

var tpl = template.Must(template.New("template").Funcs(funcs).Parse(tplText))

var funcs = map[string]interface{}{
	"quote":        fnQuote,
	"valueType":    fnValueType,
	"valueTypeCap": fnValueTypeCap,
}

func fnValueType(enum *defs.Enum) string {
	return enum.Basic.Name()
}

var mapCap = map[types.BasicKind]string{
	types.Int:    "Int",
	types.Uint64: "Uint64",
}

func fnValueTypeCap(enum *defs.Enum) string {
	return mapCap[enum.Basic.Kind()]
}

func fnQuote(s string) string {
	return fmt.Sprintf("%q", s)
}
