package filtergen

import (
	"fmt"
	"go/types"
	"reflect"
	"strings"
	"text/template"

	"etop.vn/backend/tools/pkg/generator"
	"etop.vn/backend/tools/pkg/typedesc"
)

type ColumnDef struct {
	ColumnName string
	FieldName  string
	FieldType  types.Type
	TypeDesc   *typedesc.TypeDesc

	origPkg bool
}

type PathDef struct {
	Ptr  bool
	Path string
}

type TableDef struct {
	TableName  string
	StructName string
	Cols       []*ColumnDef
}

type JoinTableDef struct {
	StructName string
	SubStructs []string
}

type Gen struct {
	PkgPath    string
	tables     []*TableDef
	joinTables []*JoinTableDef

	importOrigPkg bool
}

func NewGen(pkgPath string) *Gen {
	return &Gen{
		PkgPath: pkgPath,
	}
}

func (g *Gen) AddTable(tableName string, structName string, cols []*ColumnDef) {
	if g == nil {
		return
	}
	g.tables = append(g.tables, &TableDef{
		TableName:  tableName,
		StructName: structName,
		Cols:       cols,
	})

	for _, col := range cols {
		ch := col.TypeDesc.TypeString[0]
		if ch == '*' {
			ch = col.TypeDesc.TypeString[1]
		}
		if ch >= 'A' && ch <= 'Z' {
			col.origPkg = true
			g.importOrigPkg = true
		}
	}
}

func (g *Gen) AddJoinTable(structName string, substructs []string) {
	// g.joinTables = append(g.joinTables, &JoinTableDef{
	// 	StructName: structName,
	// 	SubStructs: substructs,
	// })
}

var printer generator.Printer
var funcMap = template.FuncMap{
	"baseName":   fnBaseName,
	"filterType": fnFilterType,
	"generate":   fnGenerate,
	"genIsZero":  fnGenIsZero,
	"isPtr":      fnIsPtr,
	"ptrType":    fnPtrType,
	"type":       fnType,
}

func fnBaseName(s string) string {
	parts := strings.Split(s, ".")
	return parts[len(parts)-1]
}

func fnGenerate(col *ColumnDef) bool {
	return col.TypeDesc.IsBasic() || col.TypeDesc.IsTime()
}

func fnType(col *ColumnDef) string {
	s := printer.TypeString(col.FieldType)
	return s
}

func fnPtrType(col *ColumnDef) string {
	s := printer.TypeString(col.FieldType)
	if col.TypeDesc.Ptr {
		return s
	}
	return "*" + s
}

func fnFilterType(col *ColumnDef) string {
	if col.TypeDesc.Ptr {
		return "sq.ColumnFilterPtr"
	}
	return "sq.ColumnFilter"
}

func fnIsPtr(col *ColumnDef) bool {
	return col.TypeDesc.Ptr
}

func fnGenIsZero(ptr bool, col *ColumnDef) string {
	// overwrite pointer type
	desc := *col.TypeDesc
	desc.Ptr = false

	path := col.FieldName
	if ptr {
		path = "(*" + path + ")"
	}
	res := genIsZero(path, &desc, col.FieldType)
	if desc.Elem == reflect.Bool && desc.Container == 0 {
		res = "bool(" + res + ")"
	}
	return res
}

func genIsZero(path string, desc *typedesc.TypeDesc, typ types.Type) string {
	switch {
	case desc.IsBareTime():
		return path + ".IsZero()"
	case desc.IsNillable():
		return path + " == nil"
	case desc.IsNullType(typ):
		return "!" + path + ".Valid"
	case desc.IsNumber():
		return path + " == 0"
	case desc.IsKind(reflect.Bool):
		return "!" + path
	case desc.IsKind(reflect.String):
		return path + ` == ""`
	case desc.IsKind(reflect.Struct):
		return "false"
	}

	panic("unsupported type: " + desc.TypeString)
}

var tpl = template.Must(template.New("tpl").Funcs(funcMap).Parse(tplStr))

func (g *Gen) Generate(p generator.Printer) {
	if g == nil {
		return
	}
	vars := map[string]interface{}{
		"Tables":     g.tables,
		"JoinTables": g.joinTables,
	}
	p.Import("sq", "etop.vn/backend/pkg/common/sql/sq")
	p.Import("etopmodel", "etop.vn/backend/pkg/etop/model")
	p.Import("orderingtypes", "etop.vn/api/main/ordering/types") // TODO: remove this

	printer = p
	err := tpl.Execute(p, vars)
	if err != nil {
		panic(fmt.Sprintf("Unable to generate filters: %v", err))
	}
}
