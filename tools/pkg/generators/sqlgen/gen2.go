package sqlgen

import (
	"fmt"
	"go/types"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/dustin/go-humanize/english"

	"o.o/backend/tools/pkg/gen"
	"o.o/backend/tools/pkg/generator"
	"o.o/backend/tools/pkg/genutil"
	"o.o/common/strs"
)

var funcs = template.FuncMap{
	"concat":    fnConcat,
	"go":        fnGo,
	"quote":     fnQuote,
	"nonzero":   fnNonZero,
	"updateArg": fnUpdateArg,
	"plural":    fnPlural,
	"toTitle":   fnToTitle,
	"typeName":  fnTypeName,
	"baseName":  fnBaseName,

	"tableForType":    fnTableForType,
	"listColsForType": fnListColsForType,
}

var tpl = template.Must(template.New("tpl").Funcs(funcs).Parse(tplStr))

func fnConcat(ss ...string) string {
	return strings.Join(ss, "")
}

func fnGo(v interface{}) string {
	switch vv := v.(type) {
	case []byte:
		v = string(vv)
	}
	return fmt.Sprintf("%#v", v)
}

func fnQuote(v interface{}) string {
	return strings.Replace(fmt.Sprintf("%#v", v), `"`, `\"`, -1)
}

func fnTableForType(typ types.Type) string {
	ts := pr.TypeString(typ)
	return fmt.Sprintf("(*%v)(nil).SQLTableName()", ts)
}

func fnListColsForType(typ types.Type) string {
	ts := pr.TypeString(typ)
	return fmt.Sprintf("(*%v)(nil).SQLListCols()", ts)
}

func fnNonZero(p generator.Printer, col *colDef) string {
	return genIfNotEqualToZero(p, col)
}

func fnUpdateArg(p generator.Printer, col *colDef) string {
	return genUpdateArg(p, col)
}

func fnTypeName(typ types.Type) string {
	name := pr.TypeString(typ)
	if name[0] == '*' {
		name = name
	}
	return name
}

func fnBaseName(s string) string {
	parts := strings.Split(s, ".")
	return parts[len(parts)-1]
}

func fnPlural(n int, word string) string {
	return english.Plural(n, word, "")
}

func fnToTitle(s string) string {
	s = strs.ToTitle(s)
	s = strings.Replace(s, "Id", "ID", -1)
	return s
}

func (g *genImpl) Generate(typ types.Type) error {
	if err := g.validateTypes(); err != nil {
		return err
	}

	g.generateCommon()
	return g.genQueryFor(typ)
}

func (g *genImpl) generateCommon() {
	if g.init {
		return
	}
	g.init = true
	p := g.Printer

	str := `
var __sqlModels []interface{ SQLVerifySchema(db *cmsql.Database) }
var __sqlonce sync.Once

func SQLVerifySchema(db *cmsql.Database) {
	__sqlonce.Do(func() {
		for _, m := range __sqlModels {
			m.SQLVerifySchema(db)
		}
	})
}

type SQLWriter = core.SQLWriter
`
	w(p, str)
}

func (g *genImpl) genQueryFor(typ types.Type) (_err error) {
	defer func() {
		g.nGen++
		if g.nGen == g.nAdd {
			if g.genFilter == nil {
				return
			}
			pkgSortPath := strings.TrimPrefix(g.genFilter.PkgPath, "o.o/backend/")
			dir := filepath.Join(gen.ProjectPath(), pkgSortPath, "../sqlstore")
			sqlstorePath, err := filepath.Abs(dir)
			if err != nil {
				panic(err)
			}
			fi, err := os.Stat(sqlstorePath)
			if err != nil {
				fmt.Printf("can not find package sqlstore: %v", err)
				return
			}
			if !fi.IsDir() {
				panic(fmt.Sprintf("%v is not a directory", sqlstorePath))
			}
			filePath := filepath.Join(sqlstorePath, "filters.gen.go")
			p, err := g.ng.GenerateFile("sqlstore", filePath)
			if err != nil {
				_err = err
				return
			}
			g.genFilter.Generate(p)
			_err = p.Close()
		}
	}()

	p := g.Printer
	def := g.mapType[typ.String()]
	pStr := pr.TypeString(typ)
	Str := pStr
	Strs := plural(Str)
	tableName := def.tableName

	// generate convert methods
	if def.base != nil && len(def.joins) == 0 {
		if err := g.genConvertMethodsFor(def.typ, def.base); err != nil {
			return err
		}
	}

	extra := ""
	if def.base != nil {
		extra = ", _ " + pr.TypeString(def.base)
	}
	var joinTypes, joinAs, joinConds []string
	if len(def.joins) != 0 {
		extra += ", as sq.AS"
		joinTypes = make([]string, len(def.joins))
		joinAs = make([]string, len(def.joins))
		joinConds = make([]string, len(def.joins))
		for i, join := range def.joins {
			joinTypes[i] = fmt.Sprintf("t%v", i)
			joinAs[i] = join.JoinAlias
			joinConds[i] = join.JoinCond
		}
	}

	var ptrElems []pathElem
	for _, s := range def.structs {
		if s.ptr {
			ptrElems = append(ptrElems, s)
		}
	}

	vars := map[string]interface{}{
		"p":         p,
		"IsSimple":  len(def.joins) == 0,
		"IsJoin":    len(def.joins) != 0,
		"IsPreload": len(def.preloads) > 0,
		"IsAll":     def.all,
		"IsSelect":  def.selecT,
		"IsInsert":  def.insert,
		"IsUpdate":  def.update,
		"IsNow":     "",

		"DeriveFuncName": "sqlgen" + genutil.ExtractNamed(typ).Obj().Name(),
		"FuncExtraArgs":  extra,

		"BaseType":                 def.base,
		"TypeName":                 Str,
		"TypeNames":                Strs,
		"TableName":                tableName,
		"Cols":                     def.cols,
		"ColNamesAndTypes":         getMapColTypes(def.cols),
		"ColsList":                 listColumns("", def.cols),
		"ColsListUpdateOnConflict": listUpdateOnConflictColumns("", def.cols),
		"QueryArgs":                listInsertArgs(p, def.cols),
		"NumCols":                  len(def.cols),
		"NumJoins":                 len(def.joins),
		"PtrElems":                 ptrElems,
		"ScanArgs":                 listScanArgs(p, def.cols),
		"TimeLevel":                def.timeLevel,

		"As":    def.alias,
		"Joins": def.joins,

		"Preloads": def.preloads,

		"_ListCols":           fmt.Sprintf("__sql%v_ListCols", Str),
		"_ListColsOnConflict": fmt.Sprintf("__sql%v_ListColsOnConflict", Str),
		"_Table":              fmt.Sprintf("__sql%v_Table", Str),
		"_Insert":             fmt.Sprintf("__sql%v_Insert", Str),
		"_Select":             fmt.Sprintf("__sql%v_Select", Str),
		"_UpdateAll":          fmt.Sprintf("__sql%v_UpdateAll", Str),
		"_UpdateOnConflict":   fmt.Sprintf("__sql%v_UpdateOnConflict", Str),
	}

	var b strings.Builder
	b.Grow(len(tplStr) * 3 / 2)
	if err := tpl.Execute(&b, vars); err != nil {
		return err
	}

	w(p, b.String())
	return nil
}

func w(p io.Writer, format string, args ...interface{}) {
	_, err := fmt.Fprintf(p, format, args...)
	if err != nil {
		panic(err)
	}
}

func getMapColTypes(cols []*colDef) string {
	result := ""
	for _, col := range cols {
		enumValues := "[]string{"
		if len(col.columnEnumValues) > 0 {
			for _, val := range col.columnEnumValues {
				enumValues += fmt.Sprintf("%q,", val)
			}
			enumValues = enumValues[:len(enumValues)-1]
		}
		enumValues += "}"
		result += fmt.Sprintf("%q : {\nColumnName: %q,\nColumnType: %q,\nColumnDBType: %q,\nColumnTag: %q,\nColumnEnumValues: %s,\n},\n", col.ColumnName, col.ColumnName, col.columnType, col.columnDBType, col.columnDBTag, enumValues)
	}

	return fmt.Sprintf("map[string]migration.ColumnDef{\n%s}", result)
}
