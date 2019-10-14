package apix

const tplText = `
{{range $s := .Services}}
type {{.Name}}Impl struct {
// TODO
}

{{range $m := .Methods}}
func (s *{{$s.Name}}Impl) {{.Name}}(ctx context.Context, req *{{.Request.Type|type}}) (*{{.Response.Type|type}}, error) {
    panic("TODO")
}
{{end}}
{{end}}
`
