package convert

import (
	"fmt"
	"go/types"
	"reflect"
	"strings"
	"text/template"

	"etop.vn/backend/tools/pkg/generator"
	"etop.vn/backend/tools/pkg/genutil"
	"etop.vn/capi/dot"
)

var tplConvertType, tplUpdate, tplCreate *template.Template
var currentPrinter generator.Printer
var capiPkgPath = reflect.TypeOf(dot.NullString{}).PkgPath()

func init() {
	funcMap := map[string]interface{}{
		"fieldName":   renderFieldName,
		"fieldValue":  renderFieldValue,
		"fieldApply":  renderFieldApply,
		"lastComment": renderLastComment,
		"plural":      renderPlural,
	}
	parse := func(text string) *template.Template {
		return template.Must(template.New("convert_type").Funcs(funcMap).Parse(text))
	}

	tplConvertType = parse(tplConvertTypeText)
	tplCreate = parse(tplCreateText)
	tplUpdate = parse(tplUpdateText)
}

func renderPlural(s string) string {
	return genutil.Plural(s)
}

var lastComment string

func renderLastComment() string {
	return lastComment
}

func renderFieldName(field fieldConvert) string {
	return field.Out.Name()
}

func renderFieldValue(prefix string, field fieldConvert) string {
	in, out := field.Arg, field.Out
	if in == nil {
		lastComment = "// zero value"
		return renderZero(out.Type())
	}
	if out.Type() == in.Type() {
		lastComment = "// simple assign"
		return prefix + "." + in.Name()
	}
	if result := renderSimpleConversion(in, out, prefix); result != "" {
		return result
	}
	lastComment = "// types do not match"
	return renderZero(out.Type())
}

func renderFieldApply(prefix string, field fieldConvert) string {
	arg, out := field.Arg, field.Out
	if field.IsIdentifier {
		lastComment = "// identifier"
		return "out." + out.Name()
	}
	if arg == nil {
		lastComment = "// no change"
		return "out." + out.Name()
	}
	// render NullString, NullInt, ...Apply()
	if argType, ok := arg.Type().(*types.Named); ok {
		typObj := argType.Obj()
		if typObj.Pkg().Path() == capiPkgPath &&
			strings.HasPrefix(typObj.Name(), "Null") {
			lastComment = "// apply change"
			return prefix + "." + arg.Name() + ".Apply(out." + out.Name() + ")"
		}
	}
	if result := renderSimpleConversion(arg, out, prefix); result != "" {
		return result
	}
	lastComment = "// types do not match"
	return "out." + out.Name()
}

func renderSimpleConversion(in, out *types.Var, prefix string) string {
	// convert basic types
	inBasic := checkBasicType(in.Type())
	outBasic := checkBasicType(out.Type())
	if inBasic != nil && outBasic != nil {
		outStr := currentPrinter.TypeString(out.Type())
		if inBasic.Kind() == outBasic.Kind() {
			lastComment = "// simple conversion"
			return outStr + "(" + prefix + "." + in.Name() + ")"
		}
		if inBasic.Info()&types.IsNumeric > 0 && outBasic.Info()&types.IsNumeric > 0 {
			lastComment = "// simple conversion"
			return outStr + "(" + prefix + "." + in.Name() + ")"
		}
	}
	return ""
}

func renderZero(typ types.Type) string {
	t := typ
	for ok := true; ok; _, ok = t.(*types.Named) {
		t = t.Underlying()
	}
	switch t := t.(type) {
	case *types.Basic:
		info := t.Info()
		switch {
		case info&types.IsBoolean > 0:
			return "false"
		case info&types.IsNumeric > 0:
			return "0"
		case info&types.IsString > 0:
			return `""`
		default:
			return "0"
		}

	case *types.Struct:
		if t == typ {
			if t.NumFields() == 0 {
				return "struct{}{}"
			}
			panic(fmt.Sprintf("struct must have a name (%v)", t))
		}
		return currentPrinter.TypeString(typ) + "{}"

	default:
		return "nil"
	}
}

func checkBasicType(typ types.Type) *types.Basic {
	for ok := true; ok; _, ok = typ.(*types.Named) {
		typ = typ.Underlying()
	}
	basic, _ := typ.(*types.Basic)
	return basic
}
