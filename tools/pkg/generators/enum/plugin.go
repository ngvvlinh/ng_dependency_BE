package enum

import (
	"fmt"
	"go/types"
	"sort"
	"strings"
	"text/template"

	"golang.org/x/tools/go/packages"

	"o.o/backend/tools/pkg/generator"
	"o.o/backend/tools/pkg/generators/api/defs"
	"o.o/backend/tools/pkg/generators/api/parse"
	"o.o/backend/tools/pkg/genutil"
)

type keyWithNull struct{}
type keyModelType struct{}
type keyModelZero struct{}

type plugin struct {
	generator.Qualifier
}

func New() generator.Plugin {
	return &plugin{Qualifier: genutil.Qualifier{}}
}

func (p *plugin) Name() string { return "enum" }

func (p *plugin) Filter(ng generator.FilterEngine) error {
	currentInfo = parse.NewInfo(ng)
	for _, pkg := range ng.ParsingPackages() {
		ds := pkg.InlineDirectives
		if _, ok := ds.Get("enum"); ok {
			pkg.Include()
		}
	}
	return nil
}

func (p *plugin) Generate(ng generator.Engine) error {
	currentInfo.Init(ng)
	return ng.GenerateEachPackage(p.generatePackage)
}

func (p *plugin) generatePackage(ng generator.Engine, pkg *packages.Package, printer generator.Printer) error {
	mapEnum, err := parse.ParseEnumInPackage(ng, pkg)
	if err != nil {
		return err
	}

	// parse directives
	enums := make([]*defs.Enum, 0, len(mapEnum))
	for _, enum := range mapEnum {
		if err := parseDirectives(ng, enum); err != nil {
			return err
		}
		enums = append(enums, enum)
	}
	sort.Slice(enums, func(i, j int) bool { return enums[i].Name < enums[j].Name })

	// parse NullEnum
	for _, obj := range ng.GetObjectsByPackage(pkg) {
		_, ok := obj.(*types.TypeName)
		if !ok {
			continue
		}
		if !strings.HasPrefix(obj.Name(), "Null") {
			continue
		}
		enumName := strings.TrimPrefix(obj.Name(), "Null")
		if enumName == "" {
			return generator.Errorf(nil, "invalid name (%v)", obj.Name())
		}
		if _, ok := ng.GetDirectives(obj).Get("enum"); ok {
			return generator.Errorf(nil, "%v must not have enum directive", obj.Name())
		}
		enum := mapEnum[enumName]
		if enum == nil {
			return generator.Errorf(nil, "enum for %v not found", obj.Name())
		}
		if !currentInfo.IsNullStruct(obj.Type(), "Enum") {
			return generator.Errorf(nil, "%v must be struct { Enum <Enum> ; Valid bool }", obj.Name())
		}
		st := obj.Type().Underlying().(*types.Struct)
		if st.Field(0).Type() != enum.Type {
			return generator.Errorf(nil, "%v must be struct { Enum <Enum> ; Valid bool }", obj.Name())
		}
		currentInfo.Set(enum, keyWithNull{}, true)
	}

	printer.Import("dot", "o.o/capi/dot")
	printer.Import("driver", "database/sql/driver")
	printer.Import("fmt", "fmt")
	printer.Import("mix", "o.o/capi/mix")
	vars := map[string]interface{}{
		"Enums": enums,
	}
	return tpl.Execute(printer, vars)
}

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
	"withNull":     fnWithNull,
}

func fnModelType(enum *defs.Enum) string {
	modelType := currentInfo.Get(enum, keyModelType{})
	if modelType != nil {
		return modelType.(string)
	}
	return ""
}

func fnQuote(s string) string {
	return fmt.Sprintf("%q", s)
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

func fnWithNull(enum *defs.Enum) bool {
	withNull := currentInfo.Get(enum, keyWithNull{})
	return withNull != nil
}
