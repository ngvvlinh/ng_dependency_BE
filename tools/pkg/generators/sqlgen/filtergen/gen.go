package filtergen

import (
	"bytes"
	"fmt"
	"go/types"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"

	"etop.vn/backend/tools/pkg/gen"
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
	pkgPath    string
	tables     []*TableDef
	joinTables []*JoinTableDef

	importOrigPkg bool
}

func NewGen(pkgPath string) *Gen {
	return &Gen{
		pkgPath: pkgPath,
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
	s := col.TypeDesc.TypeString
	if strings.HasPrefix(s, "etop_vn_backend_pkg_etop_model") {
		s = strings.TrimPrefix(s, "etop_vn_backend_pkg_etop_")
	}
	if col.origPkg {
		s = addPrefix(s, "m.")
	}
	return s
}

func fnPtrType(col *ColumnDef) string {
	s := col.TypeDesc.TypeString
	if col.origPkg {
		s = addPrefix(s, "m.")
	}
	if col.TypeDesc.Ptr {
		return s
	}
	return "*" + s
}

func addPrefix(s string, prefix string) string {
	if s[0] == '*' {
		return "*" + prefix + s[1:]
	}
	return "m." + s
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

func (g *Gen) Generate() {
	if g == nil {
		return
	}
	var b bytes.Buffer
	vars := map[string]interface{}{
		"Tables":      g.tables,
		"JoinTables":  g.joinTables,
		"Imports":     "",
		"OrigPackage": "",
	}
	fmt.Println("--", g.pkgPath)
	if g.importOrigPkg {
		vars["OrigPackage"] = `m "` + g.pkgPath + `"`
	}
	if g.pkgPath != "etop.vn/backend/pkg/etop/model" {
		vars["Imports"] = `"etop.vn/backend/pkg/etop/model"`
	}

	err := tpl.Execute(&b, vars)
	if err != nil {
		panic(fmt.Sprintf("Unable to generate filters: %v", err))
	}

	pkgSortPath := strings.TrimPrefix(g.pkgPath, "etop.vn/backend/")
	dir := filepath.Join(gen.ProjectPath(), pkgSortPath, "../sqlstore")
	sqlstorePath, err := filepath.Abs(dir)
	must(err)
	fi, err := os.Stat(sqlstorePath)
	if err != nil {
		fmt.Printf("can not find package sqlstore: %v", err)
		return
	}
	if !fi.IsDir() {
		panic(fmt.Sprintf("%v is not a directory", sqlstorePath))
	}

	filePath := filepath.Join(sqlstorePath, "filters.gen.go")
	gen.WriteFile(filePath, b.Bytes())
	gen.FormatFiles(filePath)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
