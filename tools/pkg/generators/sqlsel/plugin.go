package sqlsel

import (
	"fmt"
	"go/types"
	"io"
	"reflect"
	"strings"
	"text/template"

	"golang.org/x/tools/go/packages"

	"o.o/backend/tools/pkg/generator"
	"o.o/backend/tools/pkg/generators/api/parse"
	"o.o/backend/tools/pkg/generators/sqlgen"
	"o.o/backend/tools/pkg/genutil"
)

const CmdPrefix = "sqlsel"

type plugin struct {
	generator.Qualifier
}

func New() generator.Plugin {
	return &plugin{
		Qualifier: genutil.Qualifier{},
	}
}

func (p *plugin) Name() string { return CmdPrefix }

func (p *plugin) Filter(ng generator.FilterEngine) error {
	sqlgen.CurrentInfo = parse.NewInfo(ng)
	return generator.FilterByCommand(CmdPrefix).FilterAll(ng)
}

func (p *plugin) Generate(ng generator.Engine) error {
	sqlgen.CurrentInfo.Init(ng)
	return ng.GenerateEachPackage(p.generateEachPackage)
}

func (p *plugin) generateEachPackage(ng generator.Engine, pkg *packages.Package, printer generator.Printer) error {
	if err := genutil.NoUnboundDirectives(ng, pkg, CmdPrefix); err != nil {
		return err
	}

	var typs []*types.Named
	for _, obj := range ng.GetObjectsByPackage(pkg) {
		typ, err := parseDirectives(ng, pkg, obj)
		if err != nil {
			return err
		}
		if typ != nil {
			typs = append(typs, typ)
		}
	}

	printer.Import("", "database/sql")
	printer.Import("core", "o.o/backend/pkg/common/sql/sq/core")
	for _, typ := range typs {
		if err := genStruct(printer, typ); err != nil {
			return err
		}
	}
	return nil
}

func parseDirectives(ng generator.Engine, pkg *packages.Package, obj types.Object) (*types.Named, error) {
	ds := ng.GetDirectives(obj).FilterBy(CmdPrefix)
	if len(ds) == 0 {
		return nil, nil
	}

	// must be struct type
	typeName, ok := obj.(*types.TypeName)
	if !ok {
		return nil, generator.Errorf(nil, "sqlsel must be used on struct declaration (%v is not)", obj)
	}
	named := typeName.Type().(*types.Named)
	if st := unwrapNamedStruct(named); st == nil {
		return nil, generator.Errorf(nil, "sqlsel must be used on struct declaration (%v is not)", obj)
	}
	return named, nil
}

func unwrapNamedStruct(named *types.Named) *types.Struct {
	underlying := genutil.UnwrapNamed(named)
	st, _ := underlying.(*types.Struct)
	return st
}

func genStruct(p generator.Printer, named *types.Named) error {
	s := unwrapNamedStruct(named)

	var selects strings.Builder
	args := make([]string, 0, s.NumFields())
	name := p.TypeString(named)
	for i, n := 0, s.NumFields(); i < n; i++ {
		tag := reflect.StructTag(s.Tag(i))
		selTag := tag.Get(`sel`)
		if selTag == "" || selTag == "-" {
			continue
		}

		if len(args) > 0 {
			selects.WriteString(", ")
		}
		selects.WriteString(cleanSpaces(selTag))

		field := s.Field(i)
		path := "m." + field.Name()
		typ := field.Type()

		arg := sqlgen.GenScanArg(p, path, typ)
		args = append(args, arg)
	}

	text := `
type {{.NamePlural}} []*{{.Name}}

func (m *{{.Name}}) SQLTableName() string { return "" }
func (m *{{.NamePlural}}) SQLTableName() string { return "" }

func (m *{{.Name}}) SQLScanArgs(opts core.Opts) []interface{} {
	return []interface{}{
{{- range .Args}}
		{{.}},
{{- end}}
	}
}

func (m *{{.Name}}) SQLScan(opts core.Opts, row *sql.Row) error {
	return row.Scan(m.SQLScanArgs(opts)...)
}

func (ms *{{.NamePlural}}) SQLScan(opts core.Opts, rows *sql.Rows) error {
	res := make({{.NamePlural}}, 0, 128)
	for rows.Next() {
		m := new({{.Name}})
		args := m.SQLScanArgs(opts)
		if err := rows.Scan(args...); err != nil {
			return err
		}
		res = append(res, m)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	*ms = res
	return nil
}

func (_ *{{.Name}}) SQLSelect(w core.SQLWriter) error {
	w.WriteRawString(` + "`" + `SELECT {{.Selects}}` + "`" + `)
	return nil
}

func (_ *{{.NamePlural}}) SQLSelect(w core.SQLWriter) error {
	w.WriteRawString(` + "`" + `SELECT {{.Selects}}` + "`" + `)
	return nil
}
`
	tpl := template.Must(template.New("tpl").Parse(text))
	return tpl.Execute(p, map[string]interface{}{
		"Name":       name,
		"NamePlural": genutil.Plural(name),
		"Args":       args,
		"Selects":    selects.String(),
	})
}

func cleanSpaces(s string) string {
	b := make([]byte, 1, len(s))
	b[0] = s[0]
	for i := 1; i < len(s); i++ {
		if s[i] == ' ' && s[i-1] == ' ' {
			// skip
		} else {
			b = append(b, s[i])
		}
	}
	return strings.TrimSpace(string(b))
}

func w(p io.Writer, format string, args ...interface{}) {
	_, err := fmt.Fprintf(p, format, args...)
	if err != nil {
		panic(err)
	}
}
