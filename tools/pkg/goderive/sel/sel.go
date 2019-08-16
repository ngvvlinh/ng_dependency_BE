package sel

import (
	"fmt"
	"go/types"
	"reflect"
	"strings"
	"text/template"

	"github.com/awalterschulze/goderive/derive"

	"etop.vn/backend/tools/pkg/sqlgen"
)

// NewPlugin creates a new clone plugin.
// This function returns the plugin name, default prefix and a constructor for the clone code generator.
func NewPlugin() derive.Plugin {
	return derive.NewPlugin("sel", "sel", New)
}

// New is a constructor for the clone code generator.
// This generator should be reconstructed for each package.
func New(typesMap derive.TypesMap, p derive.Printer, deps map[string]derive.Dependency) derive.Generator {
	return &gen{
		TypesMap: typesMap,
		Printer:  p,
	}
}

type gen struct {
	derive.TypesMap
	derive.Printer
	init bool
}

func (g *gen) Add(name string, typs []types.Type) (string, error) {
	if len(typs) == 0 {
		return "", fmt.Errorf("%s must have one or more arguments", name)
	}
	return g.SetFuncName(name, typs...)
}

func (g *gen) Generate(typs []types.Type) error {
	if !g.init {
		g.init = true
		g.NewImport("", "database/sql")()
		g.NewImport("core", "etop.vn/common/sq/core")()

		g.P(`
type SQLWriter = core.SQLWriter`)
	}

	g.Generating(typs...)
	name := g.GetFuncName(typs...)
	g.P(`func %v(_ ...interface{}) bool {return true}`, name)

	for _, typ := range typs {
		if err := g.genStruct(typ); err != nil {
			return err
		}
	}
	return nil
}

func (g *gen) genStruct(typ types.Type) error {
	s, ok := pointerToStruct(typ)
	if !ok {
		return fmt.Errorf("Type %v must be pointer to struct", g.TypeString(typ))
	}

	var selects strings.Builder
	args := make([]string, 0, s.NumFields())
	name := g.TypeString(typ)[1:]
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

		arg := sqlgen.GenScanArg(path, typ)
		args = append(args, arg)
	}

	text := `
func (m *{{.Name}}) SQLTableName() string { return "" }

func (m *{{.Name}}) SQLScan(opts core.Opts, row *sql.Row) error {
	args := []interface{}{
{{- range .Args}}
		{{.}},
{{- end}}
	}
	return row.Scan(args...)
}

func (_ *{{.Name}}) SQLSelect(w SQLWriter) error {
	w.WriteRawString(` + "`" + `SELECT {{.Selects}}` + "`" + `)
	return nil
}
`
	tpl := template.Must(template.New("tpl").Parse(text))
	var buf strings.Builder
	if err := tpl.Execute(&buf, map[string]interface{}{
		"Name":    name,
		"Args":    args,
		"Selects": selects.String(),
	}); err != nil {
		return err
	}
	g.P(buf.String())
	return nil
}

func pointerToStruct(typ types.Type) (*types.Struct, bool) {
	pt, ok := typ.Underlying().(*types.Pointer)
	if !ok {
		return nil, false
	}
	st, ok := pt.Elem().Underlying().(*types.Struct)
	return st, ok
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
