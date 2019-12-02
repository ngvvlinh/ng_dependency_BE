package api

const tplText = `
{{range $s := .Services -}}
type {{$s|busName}} struct{ bus capi.Bus }
{{end}}

{{range $s := .Services -}}
func New{{$s|busName}}(bus capi.Bus) {{$s|busName}} { return {{$s|busName}}{ bus } }
{{end}}

{{range $s := .Services -}}
func (b {{$s|busName}}) Dispatch(ctx context.Context, msg interface{ {{$s|interfaceMethod}}() }) error { return b.bus.Dispatch(ctx, msg) }
{{end}}

{{range $s := .Services}}
{{range $m := .Methods}}
type {{$m|messageName}} struct {
	{{$m|generateStruct}}

	{{$m|generateResult}}
}

{{$m|generateHandle}}
{{end}}
{{end}}

// implement interfaces
{{range $s := .Services -}}
{{range $m := .Methods -}}
func (q *{{$m|messageName}}) {{$s|interfaceMethod}}() {}
{{end}}
{{end}}

// implement conversion
{{range $s := .Services}}
{{range $m := .Methods}}
{{$m|generateGetArgs}}
{{$m|generateSetArgs}}
{{end}}
{{end}}

// implement dispatching
{{range $s := .Services}}
type {{.FullName}}Handler struct {
	inner {{.FullName}}
}

func New{{.FullName}}Handler(service {{.FullName}}) {{.FullName}}Handler { return {{.FullName}}Handler{service} }

func (h {{.FullName}}Handler) RegisterHandlers(b interface{
	capi.Bus
	AddHandler(handler interface{})
}) {{$s|busName}} {
	{{range $m := .Methods -}}
	b.AddHandler(h.Handle{{.Name}})
	{{end -}}
	return {{$s|busName}}{b}
}
{{end}}
`
