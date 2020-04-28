package sqlgen

import (
	"go/types"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/dustin/go-humanize/english"

	"o.o/backend/tools/pkg/generator"
	"o.o/backend/tools/pkg/generators/sqlgen/filtergen"
	"o.o/backend/tools/pkg/generators/sqlgen/substruct"
	"o.o/backend/tools/pkg/genutil"
	"o.o/common/strs"
)

const sqlTag = "sq"
const sqlTypeTag = "sql_type"

var pr generator.Printer

type genImpl struct {
	ng generator.Engine
	generator.Printer

	init    bool
	bases   []types.Type
	mapBase map[string]bool
	mapType map[string]*typeDef

	nAdd      int
	nGen      int
	genFilter *filtergen.Gen
}

type typeDef struct {
	typ      types.Type
	base     types.Type
	alias    string
	cols     []*colDef
	joins    []*joinDef
	preloads []*preloadDef

	tableName string
	structs   pathElems

	all    bool
	selecT bool
	insert bool
	update bool

	timeLevel timeLevel
}

type colDef struct {
	ColumnName       string
	FieldName        string
	columnDBType     string
	columnDBTag      string
	columnEnumValues []string

	fieldType  types.Type
	fieldTag   string
	columnType string
	timeLevel  timeLevel
	fkey       string
	pathElems

	exclude     bool
	_nonNilPath string
}

func (c *colDef) GenNonNilPath() string {
	if c._nonNilPath == "" {
		c._nonNilPath = genNonNilPath("m", c.pathElems)
	}
	return c._nonNilPath
}

func genNonNilPath(prefix string, path pathElems) string {
	var v string
	for _, elem := range path.BasePath() {
		if elem.ptr {
			v += prefix + "." + elem.Path + ` != nil && `
		}
	}
	if v == "" {
		return ""
	}
	return v[:len(v)-4] // remove the last " && "
}

func (c *colDef) String() string {
	return c.FieldName
}

type pathElems []pathElem

func (p pathElems) String() string {
	return p.Path()
}

func (p pathElems) Path() string {
	if p == nil {
		return "<nil>"
	}
	return p[len(p)-1].Path
}

func (p pathElems) Last() pathElem {
	return p[len(p)-1]
}

func (p pathElems) BasePath() pathElems {
	if len(p) == 0 {
		return nil
	}
	return p[:len(p)-1]
}

type pathElem struct {
	Path string
	Name string
	ptr  bool
	typ  types.Type

	basePath string
	TypeName string
}

func (p pathElems) append(field *types.Var) pathElems {
	name := field.Name()
	typ := field.Type()
	pStr := pr.TypeString(typ)
	ptr := pStr[0] == '*'
	Str := pStr
	if ptr {
		Str = pStr[1:]
	}

	elem := pathElem{
		Name:     name,
		ptr:      ptr,
		typ:      typ,
		TypeName: Str,
	}

	if p == nil {
		elem.Path = name
		elem.basePath = ""
		return []pathElem{elem}
	}

	elem.Path = p.Path() + "." + name
	elem.basePath = p.Path()
	pdef := make([]pathElem, 0, len(p)+1)
	pdef = append(pdef, p...)
	pdef = append(pdef, elem)
	return pdef
}

type joinDef struct {
	JoinKeyword string
	JoinAlias   string
	JoinCond    string
	JoinType    types.Type
	BaseType    types.Type
}

type preloadDef struct {
	FieldType     types.Type
	FieldName     string
	TableName     string
	PluralTypeStr string
	BaseType      types.Type
	Fkey          string
}

func (g *genImpl) validateTypes() error {
	// for _, def := range g.mapType {
	// 	if def.base != nil {
	// 		if !g.mapBase[def.base.String()] {
	// 			return generator.Errorf(err,
	// 				"Type %v is based on %v but the latter is not defined as a table",
	// 				gt.TypeString(def.typ), gt.TypeString(def.base))
	// 		}
	// 	}
	// }

	// TODO: Validate join
	return nil
}

var (
	reTagColumnName = regexp.MustCompile(`'[0-9A-Za-z._-]+'`)
	reTagKeyword    = regexp.MustCompile(`\b[a-z]+\b`)
	reTagSpaces     = regexp.MustCompile(`^\s*$`)
	reTagPreload    = regexp.MustCompile(`^preload,fkey:'([0-9A-Za-z._-]+)'$`)
)

func genColumnDBType(typ types.Type) string {
	result := ""

	if t, ok := typ.Underlying().(*types.Slice); ok {
		result = "[]" + genColumnDBType(t.Elem())
	}
	if t, ok := typ.Underlying().(*types.Pointer); ok {
		result = "*" + genColumnDBType(t.Elem())
	}
	if _, ok := typ.Underlying().(*types.Struct); ok {
		result = "struct"
	}
	if CurrentInfo.IsEnum(typ) {
		result = "enum"
	}
	if result == "" {
		result = typ.Underlying().String()
	}

	return result
}

func parseColumnsFromType(path pathElems, root *types.Named, sTyp *types.Struct) ([]*colDef, []*colDef, error) {
	var cols, excols []*colDef
	for i, n := 0, sTyp.NumFields(); i < n; i++ {
		field := sTyp.Field(i)
		if !field.Exported() {
			continue
		}
		fieldPath := path.append(field)

		columnDBTag := ""
		tag := ""
		if rawTag := reflect.StructTag(sTyp.Tag(i)); rawTag != "" {
			t, ok := rawTag.Lookup(sqlTag)
			if !ok && t != "" {
				return nil, nil, generator.Errorf(nil,
					"Invalid tag at `%v`.%v",
					pr.TypeString(root), fieldPath)
			}
			tag = t

			if t, ok := rawTag.Lookup(sqlTypeTag); ok {
				columnDBTag = t
			}
		}
		if strings.HasPrefix(tag, "-") {
			// Skip the field
			continue
		}

		columnName := toSnake(field.Name())
		columnType := pr.TypeString(field.Type())
		columnDBType := genColumnDBType(field.Type())
		var columnEnumValues []string
		if columnDBType == "enum" {
			columnEnumValues = CurrentInfo.GetEnum(field.Type()).Names
		}

		inline, create, update := false, false, false
		var fkey string
		if tag != "" {
			ntag := tag
			if strings.HasPrefix(ntag, "preload") {
				parts := reTagPreload.FindStringSubmatch(ntag)
				if len(parts) == 0 {
					return nil, nil, generator.Errorf(nil, "`preload` tag must have format \"preload,fkey:'<column>'\" (Did you forget the single quote?)")
				}
				tag = "preload"
				fkey = parts[1]
				goto endparse
			}
			if s := reTagColumnName.FindString(ntag); s != "" {
				columnName = s[1 : len(s)-1]
				ntag = strings.Replace(ntag, s, "", -1)
			}
			keywords := reTagKeyword.FindAllString(ntag, -1)
			for _, keyword := range keywords {
				switch keyword {
				case "inline":
					inline = true
				case "create", "created":
					create = true
					if columnType != "time.Time" && columnType != "*time.Time" {
						return nil, nil, generator.Errorf(nil, "`create` flag can only be used on time.Time or *time.Time field")
					}
				case "update", "updated":
					update = true
					if columnType != "time.Time" && columnType != "*time.Time" {
						return nil, nil, generator.Errorf(nil, "`update` flag can only be used on time.Time or *time.Time field")
					}
				default:
					return nil, nil, generator.Errorf(nil,
						"Unregconized keyword `%v` at `%v`.%v",
						keyword, pr.TypeString(root), fieldPath)
				}
				ntag = strings.Replace(ntag, keyword, "", -1)
			}
			if !reTagSpaces.MatchString(ntag) {
				return nil, nil, generator.Errorf(nil,
					"Invalid tag at `%v`.%v (Did you forget the single quote?)",
					pr.TypeString(root), fieldPath)
			}
		}

		if countFlags(inline, create, update) > 1 {
			return nil, nil, generator.Errorf(nil,
				"`inline`, `create`, `update` flags can not be used together (at `%v`.%v)", pr.TypeString(root), fieldPath)
		}
		if inline {
			typ := field.Type()
			if t, ok := typ.Underlying().(*types.Pointer); ok {
				typ = t.Elem()
			}
			if t, ok := typ.Underlying().(*types.Struct); ok {
				inlineCols, inlineExCols, err := parseColumnsFromType(fieldPath, root, t)
				if err != nil {
					return nil, nil, err
				}
				cols = append(cols, inlineCols...)
				excols = append(excols, inlineExCols...)
				continue
			}
			return nil, nil, generator.Errorf(nil,
				"`inline` can only be used with struct or *struct (at `%v`.%v)", pr.TypeString(root), fieldPath)
		}

	endparse:
		col := &colDef{
			FieldName:        field.Name(),
			fieldType:        field.Type(),
			fieldTag:         tag,
			ColumnName:       columnName,
			columnDBType:     columnDBType,
			columnDBTag:      columnDBTag,
			columnEnumValues: columnEnumValues,
			columnType:       columnType,
			pathElems:        fieldPath,
			fkey:             fkey,
			exclude:          tag == "preload",
		}
		if create {
			col.timeLevel = timeCreate
		} else if update {
			col.timeLevel = timeUpdate
		}
		if col.exclude {
			excols = append(excols, col)
		} else {
			cols = append(cols, col)
		}
	}
	return cols, excols, nil
}

func getStructsFromCols(cols []*colDef) (res []pathElem) {
	cpath := ""
	for _, col := range cols {
		elem := col.pathElems.Last()
		if elem.basePath == "" {
			continue
		}
		if elem.basePath == cpath {
			continue
		}
		cpath = elem.basePath
		res = append(res, col.pathElems.BasePath().Last())
	}
	return res
}

func listColumns(prefix string, cols []*colDef) string {
	b := make([]byte, 0, 1024)
	for i, col := range cols {
		if i > 0 {
			b = append(b, `,`...)
		}
		if prefix != "" {
			b = append(b, prefix...)
			b = append(b, '.')
		}
		b = append(b, '"')
		b = append(b, col.ColumnName...)
		b = append(b, '"')
	}
	return string(b)
}

func listUpdateOnConflictColumns(prefix string, cols []*colDef) string {
	b := make([]byte, 0, 1024)
	for i, col := range cols {
		if i > 0 {
			b = append(b, `,`...)
		}
		if prefix != "" {
			b = append(b, prefix...)
			b = append(b, '.')
		}
		b = append(b, '"')
		b = append(b, col.ColumnName...)
		b = append(b, '"')
		b = appends(b, " = EXCLUDED.")
		b = append(b, '"')
		b = append(b, col.ColumnName...)
		b = append(b, '"')
	}
	return string(b)
}

func listInsertArgs(p generator.Printer, cols []*colDef) []string {
	res := make([]string, len(cols))
	for i, col := range cols {
		res[i] = genInsertArg(p, col)
	}
	return res
}

func listScanArgs(p generator.Printer, cols []*colDef) []string {
	res := make([]string, len(cols))
	for i, col := range cols {
		res[i] = genScanArg(p, col)
	}
	return res
}

func (g *genImpl) genConvertMethodsFor(typ, base types.Type) error {
	funcName := "substruct" + genutil.ExtractNamed(typ).Obj().Name()
	return substruct.Generate(g.Printer, funcName, typ, base)
}

const helpJoin = `
    JOIN must have syntax: JoinType BaseType Condition

		JoinType  : One of JOIN, FULL_JOIN, LEFT_JOIN, RIGHT_JOIN,
                    NATUAL_JOIN, SELF_JOIN, CROSS_JOIN
		BaseType  : Must be a selectable struct.
		Condition : The join condition.
                    Use $L and $R as placeholders for table name.

	Example:
        sqlgenUserFullInfo(
            &UserFullInfo{}, &User{}, sq.AS("u"),
            sq.FULL_JOIN, &UserInfo{}, sq.AS("ui"), "$L.id = $R.user_id",
        )
        type UserFullInfo struct {
            User     *User
            UserInfo *UserInfo
        }

`

func (g *genImpl) parseJoin(typs []types.Type) (joins []*joinDef, err error) {
	if len(typs)%4 != 0 {
		return nil, generator.Errorf(err, "Invalid join definition")
	}
	for i := 0; i < len(typs); i = i + 4 {
		join, err := g.parseJoinLine(typs[i:])
		if err != nil {
			return nil, err
		}
		joins = append(joins, join)
	}
	return joins, nil
}

func (g *genImpl) parseJoinLine(typs []types.Type) (*joinDef, error) {
	if pr.TypeString(typs[0]) != "core.JoinType" {
		return nil, generator.Errorf(nil, "Invalid JoinType: must be one of predefined constants (got %v)", pr.TypeString(typs[0]))
	}

	base := typs[1]
	if _, ok := pointerToStruct(base); !ok {
		return nil, generator.Errorf(nil,
			"Invalid base type for join: must be pointer to struct (got %v)",
			pr.TypeString(base))
	}

	as := typs[2]
	if pr.TypeString(as) != "sq.AS" && pr.TypeString(as) != "string" {
		return nil, generator.Errorf(nil,
			"Invalid AS: must be sq.AS (got %v)", g.TypeString(as))
	}

	cond := typs[3]
	if pr.TypeString(cond) != "string" {
		return nil, generator.Errorf(nil,
			"Invalid condition for join: must be string (got %v)",
			pr.TypeString(cond))
	}

	return &joinDef{
		JoinType: base,
	}, nil
}

func pointerToNamedStruct(typ types.Type) (*types.Named, bool) {
	pt, ok := typ.Underlying().(*types.Pointer)
	if !ok {
		return nil, false
	}
	named, ok := pt.Elem().(*types.Named)
	if !ok {
		return nil, false
	}
	st := unwrapNamedStruct(named)
	return named, st != nil
}

func pointerToStruct(typ types.Type) (*types.Struct, bool) {
	pt, ok := typ.Underlying().(*types.Pointer)
	if !ok {
		return nil, false
	}
	st, ok := pt.Elem().Underlying().(*types.Struct)
	return st, ok
}

func plural(s string) string {
	return english.PluralWord(2, s, "")
}

func appends(b []byte, args ...interface{}) []byte {
	for _, arg := range args {
		switch arg := arg.(type) {
		case byte:
			b = append(b, arg)
		case rune:
			b = append(b, byte(arg))
		case string:
			b = append(b, arg...)
		case []byte:
			b = append(b, arg...)
		case int:
			b = strconv.AppendInt(b, int64(arg), 10)
		case int64:
			b = strconv.AppendInt(b, arg, 10)
		default:
			panic("Unsupport arg type: " + reflect.TypeOf(arg).Name())
		}
	}
	return b
}

func toSnake(s string) string {
	return strs.ToSnake(s)
}

func tableNameFromType(s string) string {
	return toSnake(fnBaseName(s))
}

func countFlags(args ...bool) int {
	c := 0
	for _, arg := range args {
		if arg {
			c++
		}
	}
	return c
}
