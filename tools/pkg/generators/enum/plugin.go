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
	currentInfo = parse.NewInfo(ng)
	return ng.GenerateEachPackage(p.generatePackage)
}

func (p *plugin) generatePackage(ng generator.Engine, pkg *packages.Package, printer generator.Printer) error {
	mapEnum, err := parse.ParseEnumInPackage(ng, pkg)
	if err != nil {
		return err
	}
	enums := make([]*defs.Enum, 0, len(mapEnum))
	for _, enum := range mapEnum {
		if err := parseDirectives(ng, enum); err != nil {
			return err
		}
		enums = append(enums, enum)
	}
	sort.Slice(enums, func(i, j int) bool { return enums[i].Name < enums[j].Name })

	printer.Import("dot", "etop.vn/capi/dot")
	printer.Import("driver", "database/sql/driver")
	printer.Import("fmt", "fmt")
	printer.Import("mix", "etop.vn/capi/mix")
	vars := map[string]interface{}{
		"Enums": enums,
	}
	return tpl.Execute(printer, vars)
}

type keyModelType struct{}
type keyModelZero struct{}

func parseDirectives(ng generator.Engine, enum *defs.Enum) error {
	obj := enum.Type.Obj()
	ds := ng.GetDirectives(obj)
	{
		modelType := ds.GetArg("enum:sql")
		switch modelType {
		case "":
			// no-op
		case "int", "uint64":
			currentInfo.Set(enum, keyModelType{}, modelType)
		default:
			return generator.Errorf(nil, "invalid enum:sql for %v.%v", obj.Pkg().Path(), obj.Name())
		}
	}
	{
		zeroType := ds.GetArg("enum:zero")
		switch zeroType {
		case "":
		// no-op
		case "null":
			currentInfo.Set(enum, keyModelZero{}, true)
		default:
			return generator.Errorf(nil, "invalid enum:zero for %v.%v", obj.Pkg().Path(), obj.Name())
		}
	}
	return nil
}

var currentInfo *parse.Info
var tpl = template.Must(template.New("template").Funcs(funcs).Parse(tplText))

var funcs = map[string]interface{}{
	"modelType":    fnModelType,
	"quote":        fnQuote,
	"zeroAsNull":   fnZeroAsNull,
	"valueType":    fnValueType,
	"valueTypeCap": fnValueTypeCap,
}

func fnModelType(enum *defs.Enum) string {
	modelType := currentInfo.Get(enum, keyModelType{})
	if modelType != nil {
		return modelType.(string)
	}
	return ""
}

func fnZeroAsNull(enum *defs.Enum) bool {
	zero := currentInfo.Get(enum, keyModelZero{})
	return zero != nil
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
