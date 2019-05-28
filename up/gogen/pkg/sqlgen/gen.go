package sqlgen

import (
	"fmt"
	"go/types"
	"strings"
	"text/template"

	cmstr "etop.vn/backend/pkg/common/str"
	"github.com/dustin/go-humanize/english"
)

var funcs = template.FuncMap{
	"join":      fnJoin,
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

func fnJoin(ss []string) string {
	return strings.Join(ss, ",")
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
	ts := gt.TypeString(typ)[1:]
	return fmt.Sprintf("(*%v)(nil).SQLTableName()", ts)
}

func fnListColsForType(typ types.Type) string {
	ts := gt.TypeString(typ)[1:]
	return fmt.Sprintf("(*%v)(nil).SQLListCols()", ts)
}

func fnNonZero(col *colDef) string {
	return genIfNotEqualToZero(col)
}

func fnUpdateArg(col *colDef) string {
	return genUpdateArg(col)
}

func fnTypeName(typ types.Type) string {
	name := gt.TypeString(typ)
	if name[0] == '*' {
		name = name[1:]
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
	s = cmstr.ToTitle(s)
	s = strings.Replace(s, "Id", "ID", -1)
	return s
}

func (g *gen) Generate(typs []types.Type) error {
	if err := g.validateTypes(); err != nil {
		return err
	}

	g.generateCommon()

	return g.genQueryFor(typs[0])
}

func (g *gen) generateCommon() {
	if g.init {
		return
	}
	g.init = true

	str := `
type SQLWriter = core.SQLWriter
`
	p := g.printer
	p.P(str)
}

func (g *gen) genQueryFor(typ types.Type) error {
	defer func() {
		g.nGen++
		if g.nGen == g.nAdd {
			g.genFilter.Generate()
		}
	}()

	p := g.printer
	g.Generating(typ)
	def := g.mapType[typ.String()]
	pStr := gt.TypeString(typ)
	Str := pStr[1:]
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
		extra = ", _ " + gt.TypeString(def.base)
	}
	var joinTypes, joinAs, joinConds []string
	if len(def.joins) != 0 {
		extra += ", as sq.AS"
		joinTypes = make([]string, len(def.joins))
		joinAs = make([]string, len(def.joins))
		joinConds = make([]string, len(def.joins))
		for i, join := range def.joins {
			extra += fmt.Sprintf(
				", t%v sq.JOIN_TYPE, _ %v, a%v sq.AS, c%v string",
				i, gt.TypeString(join.JoinType), i, i)
			joinTypes[i] = fmt.Sprintf("t%v", i)
			joinAs[i] = fmt.Sprintf(`a%v`, i)
			joinConds[i] = fmt.Sprintf("c%v", i)
		}
	}

	var ptrElems []pathElem
	for _, s := range def.structs {
		if s.ptr {
			ptrElems = append(ptrElems, s)
		}
	}

	vars := map[string]interface{}{
		"IsSimple":  len(def.joins) == 0,
		"IsJoin":    len(def.joins) != 0,
		"IsPreload": len(def.preloads) > 0,
		"IsAll":     def.all,
		"IsSelect":  def.selecT,
		"IsInsert":  def.insert,
		"IsUpdate":  def.update,
		"IsNow":     "",

		"DeriveFuncName": g.GetFuncName(typ),
		"FuncExtraArgs":  extra,

		"BaseType":  def.base,
		"TypeName":  Str,
		"TypeNames": Strs,
		"TableName": tableName,
		"Cols":      def.cols,
		"ColsList":  listColumns("", def.cols),
		"QueryArgs": listInsertArgs(def.cols),
		"NumCols":   len(def.cols),
		"NumJoins":  len(def.joins),
		"PtrElems":  ptrElems,
		"ScanArgs":  listScanArgs(def.cols),
		"TimeLevel": def.timeLevel,

		"Joins":     def.joins,
		"JoinTypes": joinTypes,
		"JoinAs":    joinAs,
		"JoinConds": joinConds,

		"Preloads": def.preloads,

		"_ListCols":  fmt.Sprintf("__sql%v_ListCols", Str),
		"_Table":     fmt.Sprintf("__sql%v_Table", Str),
		"_Insert":    fmt.Sprintf("__sql%v_Insert", Str),
		"_Select":    fmt.Sprintf("__sql%v_Select", Str),
		"_UpdateAll": fmt.Sprintf("__sql%v_UpdateAll", Str),
		"_JoinTypes": fmt.Sprintf("__sql%v_JoinTypes", Str),
		"_Join":      fmt.Sprintf("__sql%v_Join", Str),
		"_JoinConds": fmt.Sprintf("__sql%v_JoinConds", Str),
		"_As":        fmt.Sprintf("__sql%v_As", Str),
		"_JoinAs":    fmt.Sprintf("__sql%v_JoinAs", Str),
	}

	var b strings.Builder
	b.Grow(len(tplStr) * 3 / 2)
	if err := tpl.Execute(&b, vars); err != nil {
		return err
	}

	p.P(b.String())
	return nil
}
