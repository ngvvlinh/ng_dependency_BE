package substruct

import (
	"fmt"
	"go/types"
	"io"
	"regexp"
	"strings"

	"github.com/dustin/go-humanize/english"

	"etop.vn/backend/tools/pkg/generator"
)

func Generate(p generator.Printer, name string, out, in types.Type) error {
	sIn, ok := namedStruct(in)
	if !ok {
		return generator.Errorf(nil, "Expected pointer to struct (got %v)", in.String())
	}
	sOut, ok := namedStruct(out)
	if !ok {
		return generator.Errorf(nil, "Expected pointer to struct (got %v)", out.String())
	}

	inStr := p.TypeString(in)
	outStr := p.TypeString(out)
	sInStr := inStr[1:]
	sOutStr := outStr[1:]

	inMap := make(map[string]types.Type)
	for i, n := 0, sIn.NumFields(); i < n; i++ {
		v := sIn.Field(i)
		inMap[v.Name()] = v.Type()
	}

	outMap := make(map[string]types.Type)
	for i, n := 0, sOut.NumFields(); i < n; i++ {
		v := sOut.Field(i)
		outMap[v.Name()] = v.Type()

		name := v.Name()
		vInType := inMap[name]
		if vInType == nil {
			return generator.Errorf(nil, "Field (%v).%v does not exist in (%v)", outStr, name, inStr)
		}

		vStr := p.TypeString(v.Type())
		vInStr := p.TypeString(vInType)
		if vStr != vInStr {
			return generator.Errorf(nil,
				"Field (%v).%v has different type with (%v).%v: Expect `%v`, got `%v`",
				inStr, name, outStr, name, vInStr, vStr)
		}
	}

	// If the name does not start with "substruct", we replace the prefix.
	if !strings.HasPrefix(name, "substruct") {
		re := regexp.MustCompile(`^[a-z]+`)
		prefix := re.FindString(name)
		name = "substruct" + name[len(prefix):]
	}

	w(p, "")
	w(p, "// %v is a substruct of %v", outStr, inStr)
	w(p, "func %v(_ %v, _ %v) bool { return true }", name, outStr, inStr)

	w(p, "")
	w(p, "func %vFrom%v(ps []%v) []%v {", plural(capitalize(sOutStr)), plural(capitalize(sInStr)), inStr, outStr)
	w(p, "\tss := make([]%v, len(ps))", outStr)
	w(p, "\tfor i, p := range ps {")
	w(p, "\t\tss[i] = New%vFrom%v(p)", capitalize(sOutStr), capitalize(sInStr))
	w(p, "\t}")
	w(p, "\treturn ss")
	w(p, "}")

	w(p, "")
	w(p, "func %vTo%v(ss []%v) []%v {", plural(capitalize(sOutStr)), plural(capitalize(sInStr)), outStr, inStr)
	w(p, "\tps := make([]%v, len(ss))", inStr)
	w(p, "\tfor i, s := range ss {")
	w(p, "\t\tps[i] = s.To%v()", sInStr)
	w(p, "\t}")
	w(p, "\treturn ps")
	w(p, "}")

	w(p, "")
	w(p, "func New%vFrom%v(sp %v) %v {", capitalize(sOutStr), capitalize(sInStr), inStr, outStr)
	w(p, "\tif sp == nil {")
	w(p, `\t\treturn nil`)
	w(p, "\t}")
	w(p, "\ts := new(%v)", sOutStr)
	w(p, "\ts.CopyFrom(sp)")
	w(p, "\treturn s")
	w(p, "}")

	w(p, "")
	w(p, "func (s %v) To%v() %v {", outStr, capitalize(sInStr), inStr)
	w(p, "\tif s == nil {")
	w(p, "\t\treturn nil")
	w(p, "\t}")
	w(p, "\tsp := new(%v)", sInStr)
	w(p, "\ts.AssignTo(sp)")
	w(p, "\treturn sp")
	w(p, "}")

	w(p, "")
	w(p, "func (s %v) CopyFrom(sp %v) {", outStr, inStr)
	for i, n := 0, sOut.NumFields(); i < n; i++ {
		v := sOut.Field(i)
		name := v.Name()
		w(p, "\ts.%v = sp.%v", name, name)
	}
	w(p, "}")

	w(p, "")
	w(p, "func (s %v) AssignTo(sp %v) {", outStr, inStr)
	for i, n := 0, sOut.NumFields(); i < n; i++ {
		v := sOut.Field(i)
		name := v.Name()
		w(p, "\tsp.%v = s.%v", name, name)
	}
	w(p, "}")

	return nil
}

func namedStruct(typ types.Type) (*types.Struct, bool) {
	st, ok := typ.Underlying().(*types.Struct)
	return st, ok
}

func capitalize(s string) string {
	return strings.ToUpper(s[0:1]) + s[1:]
}

func plural(s string) string {
	return english.PluralWord(2, s, "")
}

func w(p io.Writer, format string, args ...interface{}) {
	_, err := fmt.Fprintf(p, format, args...)
	if err != nil {
		panic(err)
	}
}
