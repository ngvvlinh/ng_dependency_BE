package enum

var tplText = `
var __jsonNull = []byte("null")

{{range $enum := .Enums}}
var enum{{.Name}}Name = map[{{$enum|valueType}}]string {
{{range $v := .Values -}}
	{{$v}}: {{index $enum.MapName $v|quote}},
{{end -}}
}

var enum{{.Name}}Value = map[string]int {
{{range $name := .Names -}}
	{{$name|quote}}: {{index $enum.MapValue $name}},
{{end -}}
}

func Parse{{.Name}}(s string) ({{.Name}}, bool) {
	val, ok := enum{{.Name}}Value[s]
	return {{.Name}}(val), ok
}

func (e {{.Name}}) Enum() {{$enum|valueType}} {
	return {{$enum|valueType}}(e)
}

func (e {{.Name}}) Wrap() Null{{.Name}} {
	return Wrap{{.Name}}(e)
}

func (e {{.Name}}) Name() string {
	return enum{{.Name}}Name[e.Enum()]
}

func (e {{.Name}}) String() string {
	s, ok := enum{{.Name}}Name[e.Enum()]
	if ok {
		return s
	}
	return fmt.Sprintf("{{.Name}}(%v)", e.Enum())
}

func (e {{.Name}}) MarshalJSON() ([]byte, error) {
	return []byte("\"" + enum{{.Name}}Name[e.Enum()] + "\""), nil
}

func (e *{{.Name}}) UnmarshalJSON(data []byte) error {
	value, err := mix.UnmarshalJSONEnum{{$enum|valueTypeCap}}(enum{{.Name}}Value, data, {{.Name|quote}})
	if err != nil {
		return err
	}
	*e = {{.Name}}(value)
	return nil
}

{{if $enum|modelType }}
func (e {{.Name}}) Value() (interface{}, error) {
	return {{$enum|modelType}}(e), nil
}
{{else}}
func (e {{.Name}}) Value() (interface{}, error) {
	return e.String(), nil
}
{{end}}

func (e *{{.Name}}) Scan(src interface{}) error {
	value, err := mix.ScanEnum{{$enum|valueTypeCap}}(enum{{.Name}}Value, src, {{.Name|quote}})
	*e = ({{.Name}})(value)
	return err
}

type Null{{.Name}} struct {
	Enum {{.Name}}
	Valid bool
}

func Wrap{{.Name}}(enum {{.Name}}) Null{{.Name}} {
	return Null{{.Name}}{ Enum: enum, Valid: true }
}

func (n Null{{.Name}}) Apply(s {{.Name}}) {{.Name}} {
	if n.Valid {
		return n.Enum
	}
	return s
}

func (n Null{{.Name}}) Value() (interface{}, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Enum.Value()
}

func (n *Null{{.Name}}) Scan(src interface{}) error {
	if src == nil {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.Scan(src)
}

func (n Null{{.Name}}) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return n.Enum.MarshalJSON()
	}
	return __jsonNull, nil
}

func (n *Null{{.Name}}) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.UnmarshalJSON(data)
}
{{end}}
`
