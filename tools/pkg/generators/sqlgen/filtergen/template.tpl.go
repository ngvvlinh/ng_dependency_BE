package filtergen

var tplStr = `
{{range $table := .JoinTables}}
type {{.StructName}}Filters struct {
    {{range $table.SubStructs -}}
    {{. | baseName}} {{.}}Filters
    {{end -}}
}
{{end}}

{{range $table := .Tables}}
type {{.StructName}}Filters struct{ prefix string }

func New{{$table.StructName}}Filters(prefix string) {{$table.StructName}}Filters {
    return {{$table.StructName}}Filters{prefix}
}

func (ft *{{$table.StructName}}Filters) Filter(pred string, args ...interface{}) sq.WriterTo {
    return sq.Filter(&ft.prefix, pred, args...)
}

func (ft {{$table.StructName}}Filters) Prefix() string {
    return ft.prefix
}

{{range $col := .Cols}}
{{if $col | generate}}
{{if $col | isPtr}}
func (ft *{{$table.StructName}}Filters) By{{.FieldName}}Ptr({{.FieldName}} {{$col | type}}) *sq.ColumnFilterPtr {
    return &sq.ColumnFilterPtr{
        Prefix: &ft.prefix,
        Column: "{{.ColumnName}}",
        Value: {{.FieldName}},
        IsNil: {{.FieldName}} == nil,
        IsZero: {{.FieldName}} != nil && {{$col | genIsZero true}},
    }
}
{{else}}
func (ft *{{$table.StructName}}Filters) By{{.FieldName}}({{.FieldName}} {{$col | type}}) *sq.ColumnFilter {
    return &sq.ColumnFilter{
        Prefix: &ft.prefix,
        Column: "{{.ColumnName}}",
        Value: {{.FieldName}},
        IsNil: {{$col | genIsZero false}},
    }
}

func (ft *{{$table.StructName}}Filters) By{{.FieldName}}Ptr({{.FieldName}} {{$col | ptrType}}) *sq.ColumnFilterPtr {
    return &sq.ColumnFilterPtr{
        Prefix: &ft.prefix,
        Column: "{{.ColumnName}}",
        Value: {{.FieldName}},
        IsNil: {{.FieldName}} == nil,
        IsZero: {{.FieldName}} != nil && {{$col | genIsZero true}},
    }
}
{{end}}{{end}}{{end}}{{end}}
`
