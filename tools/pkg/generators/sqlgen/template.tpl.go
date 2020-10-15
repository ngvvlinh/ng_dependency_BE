package sqlgen

var tplStr = `
type {{.TypeNames}} []*{{.TypeName}}

{{if .IsSimple}}
	const {{._Table}} = {{.TableName | go}}
	const {{._ListCols}} = {{.ColsList | go}}
	const {{._ListColsOnConflict}} = {{.ColsListUpdateOnConflict | go}}
	const {{._Insert}} = "INSERT INTO {{.TableName | quote}} (" + {{._ListCols}} + ") VALUES"
	const {{._Select}} = "SELECT " + {{._ListCols}} + " FROM {{.TableName | quote}}"
	const {{._Select}}_history = "SELECT " + {{._ListCols}} + " FROM history.{{.TableName | quote}}"
	const {{._UpdateAll}} = "UPDATE {{.TableName | quote}} SET (" + {{._ListCols}} + ")"
	const {{._UpdateOnConflict}} = " ON CONFLICT ON CONSTRAINT {{.TableName}}_pkey DO UPDATE SET"
{{end}}

func (m *{{.TypeName}}) SQLTableName() string { return {{.TableName | go}} }
func (m *{{.TypeNames}}) SQLTableName() string { return {{.TableName | go}} }
{{if .IsSimple -}}
	func (m *{{.TypeName}}) SQLListCols() string { return {{._ListCols}} }

	func (m *{{.TypeName}}) SQLVerifySchema(db *cmsql.Database) {
		query := "SELECT " + {{._ListCols}} + " FROM \"{{.TableName}}\" WHERE false"
		if _, err := db.SQL(query).Exec(); err != nil {
			db.RecordError(err)
		}	
	}

	func (m *{{.TypeName}}) Migration(db *cmsql.Database) {
		var mDBColumnNameAndType map[string]string 
		if val, err := migration.GetColumnNamesAndTypes(db, "{{.TableName}}"); err != nil {
			db.RecordError(err)
			return
		} else {
			mDBColumnNameAndType = val 
		}
		mModelColumnNameAndType := {{.ColNamesAndTypes}}
		if err := migration.Compare(db, "{{.TableName}}", mModelColumnNameAndType, mDBColumnNameAndType); err != nil {
			db.RecordError(err)
		}
	}

	func init() {
		__sqlModels = append(__sqlModels, (*{{.TypeName}})(nil))
	}
{{end}}

{{if or .IsAll .IsInsert .IsUpdate}}
	func (m *{{.TypeName}}) SQLArgs(opts core.Opts, create bool) []interface{} {
    {{if gt .TimeLevel 0 -}}
		now := time.Now()
    {{end -}}
	return []interface{}{
    {{range .QueryArgs -}}
        {{.}},
    {{end -}}
	}
	}
{{end}}

{{if or .IsAll .IsSelect}}
	func (m *{{.TypeName}}) SQLScanArgs(opts core.Opts) []interface{} {
    {{range .PtrElems -}}
		m.{{.Name}} = new({{.TypeName}})
    {{end -}}
	return []interface{}{
    {{range .ScanArgs -}}
        {{.}},
    {{end -}}
	}
	}
{{end}}

{{if or .IsAll .IsSelect .IsJoin}}
	func (m *{{.TypeName}}) SQLScan(opts core.Opts, row *sql.Row) error {
	return row.Scan(m.SQLScanArgs(opts)...)
	}

	func (ms *{{.TypeNames}}) SQLScan(opts core.Opts, rows *sql.Rows) error {
	res := make({{.TypeNames}}, 0, 128)
	for rows.Next() {
	m := new({{.TypeName}})
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
{{end}}

{{if or .IsAll .IsSelect}}
	func (_ *{{.TypeName}}) SQLSelect(w SQLWriter) error {
	w.WriteQueryString({{._Select}})
	return nil
	}

	func (_ *{{.TypeNames}}) SQLSelect(w SQLWriter) error {
	w.WriteQueryString({{._Select}})
	return nil
	}
{{end}}

{{if or .IsAll .IsInsert}}
	func (m *{{.TypeName}}) SQLInsert(w SQLWriter) error {
	w.WriteQueryString({{._Insert}})
	w.WriteRawString(" (")
	w.WriteMarkers({{.NumCols}})
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), true))
	return nil
	}

	func (ms {{.TypeNames}}) SQLInsert(w SQLWriter) error {
	w.WriteQueryString({{._Insert}})
	w.WriteRawString(" (")
	for i := 0; i < len(ms); i++ {
	w.WriteMarkers({{.NumCols}})
	w.WriteArgs(ms[i].SQLArgs(w.Opts(), true))
	w.WriteRawString("),(")
	}
	w.TrimLast(2)
	return nil
	}
{{end}}

{{if or .IsAll .IsInsert}}
	func (m *{{.TypeName}}) SQLUpsert(w SQLWriter) error {
	m.SQLInsert(w)
	w.WriteQueryString({{._UpdateOnConflict}})
	w.WriteQueryString(" ")
	w.WriteQueryString({{._ListColsOnConflict}})
	return nil
	}

	func (ms {{.TypeNames}}) SQLUpsert(w SQLWriter) error {
	ms.SQLInsert(w)
	w.WriteQueryString({{._UpdateOnConflict}})
	w.WriteQueryString(" ")
	w.WriteQueryString({{._ListColsOnConflict}})
	return nil
	}
{{end}}

{{if or .IsAll .IsUpdate}}
	func (m *{{.TypeName}}) SQLUpdate(w SQLWriter) error {
	now, opts := time.Now(), w.Opts()
	_, _ = now, opts // suppress unuse error
	var flag bool
	w.WriteRawString("UPDATE ")
	w.WriteName({{.TableName | go}})
	w.WriteRawString(" SET ")
    {{range .Cols -}}
	  {{if isUpdatedAt $.p . | not -}}
		if {{nonzero $.p .}} {
      {{else -}}
        if true { // always update time
      {{end -}}
		flag = true
		w.WriteName({{.ColumnName | go}})
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg({{updateArg $.p .}})
		}
    {{end -}}
	if !flag {
	return core.ErrNoColumn
	}
	w.TrimLast(1)
	return nil
	}

	func (m *{{.TypeName}}) SQLUpdateAll(w SQLWriter) error {
	w.WriteQueryString({{._UpdateAll}})
	w.WriteRawString(" = (")
	w.WriteMarkers({{.NumCols}})
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), false))
	return nil
	}
{{end}}

{{if .IsJoin}}
	func (m *{{.TypeName}}) SQLSelect(w SQLWriter) error {
	(*{{.TypeName}})(nil).__sqlSelect(w)
	w.WriteByte(' ')
	(*{{.TypeName}})(nil).__sqlJoin(w)
	return nil
	}

	func (m *{{.TypeNames}}) SQLSelect(w SQLWriter) error {
	return (*{{.TypeName}})(nil).SQLSelect(w)
	}

	func (m *{{.TypeName}}) SQLJoin(w SQLWriter) error {
	m.__sqlJoin(w)
	return nil
	}

	func (m *{{.TypeNames}}) SQLJoin(w SQLWriter) error {
	return (*{{.TypeName}})(nil).SQLJoin(w)
	}

	func (m *{{.TypeName}}) __sqlSelect(w SQLWriter) {
	w.WriteRawString("SELECT ")
	core.WriteCols(w, {{.As | go}}, {{.BaseType | listColsForType}})
    {{range $i, $join := .Joins -}}
		w.WriteByte(',')
		core.WriteCols(w, {{.JoinAlias | go}}, {{$join.JoinType | listColsForType}})
    {{end -}}
	}

	func (m *{{.TypeName}}) __sqlJoin(w SQLWriter) {
	w.WriteRawString("FROM ")
	w.WriteName({{.TableName | go}})
	w.WriteRawString(" AS ")
	w.WriteName({{.As | go}})
    {{range $i, $join := .Joins -}}
		w.WriteRawString({{concat " " .JoinKeyword " " | go}})
		w.WriteName({{$join.JoinType | tableForType}})
		w.WriteRawString({{concat " AS " .JoinAlias " ON" | go}})
		w.WriteQueryString({{concat " " .JoinCond | go}})
    {{end -}}
	}

	func (m *{{.TypeName}}) SQLScanArgs(opts core.Opts) []interface{} {
	args := make([]interface{}, 0, 64) // TODO: pre-calculate length
	m.{{.BaseType | typeName | baseName}} = new({{.BaseType | typeName}})
	args = append(args, m.{{.BaseType | typeName | baseName}}.SQLScanArgs(opts)...)
    {{range $i, $join := .Joins -}}
		m.{{$join.JoinType | typeName | baseName}} = new({{$join.JoinType | typeName}})
		args = append(args, m.{{$join.JoinType | typeName | baseName}}.SQLScanArgs(opts)...)
    {{end -}}
	return args
	}
{{end}}

{{if .IsPreload}}
	func (m *{{.TypeName}}) SQLPreload(table string) *core.PreloadDesc {
	switch table {
    {{range .Preloads -}}
		case {{.TableName | go}}:
		var items {{.PluralTypeStr}}
		return &core.PreloadDesc{
		Fkey: {{.Fkey | go}},
		IDs: []interface{}{m.ID},
		Items: &items,
		}
    {{end -}}
	default:
	return nil
	}
	}

	func (m {{.TypeNames}}) SQLPreload(table string) *core.PreloadDesc {
	switch table {
    {{range .Preloads -}}
		case {{.TableName | go}}:
		ids := make([]interface{}, len(m))
		for i, item := range m {
		ids[i] = item.ID
		}
		var items {{.PluralTypeStr}}
		return &core.PreloadDesc{
		Fkey: {{.Fkey | go}},
		IDs: ids,
		Items: &items,
		}
    {{end -}}
	default:
	return nil
	}
	}
{{end}}

{{if .IsSimple}}
	type {{.TypeName}}History map[string]interface{}
	type {{.TypeName}}Histories []map[string]interface{}

	func (m *{{.TypeName}}History) SQLTableName() string { return "history.\"{{.TableName}}\"" }
	func (m {{.TypeName}}Histories) SQLTableName() string { return "history.\"{{.TableName}}\"" }

	func (m *{{.TypeName}}History) SQLSelect(w SQLWriter) error {
	w.WriteQueryString({{._Select}}_history)
	return nil
	}

	func (m {{.TypeName}}Histories) SQLSelect(w SQLWriter) error {
	w.WriteQueryString({{._Select}}_history)
	return nil
	}

    {{range .Cols -}}
		func (m {{$.TypeName}}History) {{.FieldName}}() core.Interface { return core.Interface{m[{{.ColumnName | go}}]} }
    {{end}}

	func (m *{{.TypeName}}History) SQLScan(opts core.Opts, row *sql.Row) error {
	data := make([]interface{}, {{.Cols | len}})
	args := make([]interface{}, {{.Cols | len}})
	for i := 0; i < {{.Cols | len}}; i++ {
	args[i] = &data[i]
	}
	if err := row.Scan(args...); err != nil {
	return err
	}
	res := make({{.TypeName}}History, {{.Cols | len}})
    {{range $i, $col := .Cols -}}
		res[{{$col.ColumnName | go}}] = data[{{$i}}]
    {{end -}}
	*m = res
	return nil
	}

	func (ms *{{.TypeName}}Histories) SQLScan(opts core.Opts, rows *sql.Rows) error {
	data := make([]interface{}, {{.Cols | len}})
	args := make([]interface{}, {{.Cols | len}})
	for i := 0; i < {{.Cols | len}}; i++ {
	args[i] = &data[i]
	}
	res := make({{.TypeName}}Histories, 0, 128)
	for rows.Next() {
	if err := rows.Scan(args...); err != nil {
	return err
	}
	m := make({{.TypeName}}History)
    {{range $i, $col := .Cols -}}
		m[{{$col.ColumnName | go}}] = data[{{$i}}]
    {{end -}}
	res = append(res, m)
	}
	if err := rows.Err(); err != nil {
	return err
	}
	*ms = res
	return nil
	}
{{end}}
`
