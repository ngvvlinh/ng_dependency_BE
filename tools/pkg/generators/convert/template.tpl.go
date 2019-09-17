package convert

const tplConvertTypeText = `
func convert_{{.InStr}}_{{.OutStr}}(in *{{.InType}}, out *{{.OutType}}) {
	{{- range .Fields}}
		out.{{.|fieldName}} = {{.|fieldValue "in"}} {{lastComment -}}
  {{end}}
}

func convert_{{.InStr|plural}}_{{.OutStr|plural}}(ins []*{{.InType}})(outs []*{{.OutType}}) {
  tmps := make([]{{.OutType}}, len(ins))
  outs = make([]*{{.OutType}}, len(ins))
	for i := range tmps {
    out := &tmps[i]
		outs[i] = out
		convert_{{.InStr}}_{{.OutStr}}(ins[i], out)
  }
  return outs
}
`

const tplCreateText = `
func apply_{{.ArgStr}}(arg *{{.ArgType}}, out *{{.BaseType}}) {
  {{- range .Fields}}
		out.{{.|fieldName}} = {{.|fieldValue "arg"}} {{lastComment -}}
	{{end}}
}
`

const tplUpdateText = `
func apply_{{.ArgStr}}(in *{{.BaseType}}, arg *{{.ArgType}})(out *{{.BaseType}}) {
	if in == nil {
		return nil
	}
  {{- range .Fields}}
	in.{{.|fieldName}} = {{.|fieldApply "in"}} {{lastComment -}}
  {{end}}
	return in
}
`
