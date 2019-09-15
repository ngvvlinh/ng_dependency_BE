package convert

const tplConvertTypeText = `
func convert_{{.InStr}}_{{.OutStr}}(in *{{.InType}})(out *{{.OutType}}) {
	if in == nil {
		return nil
	}
	return &{{.OutType}}{
	{{- range .Fields}}
		{{.|fieldName}}: {{.|fieldValue}}, {{lastComment -}}
  {{end}}
	}
}

func convert_{{.InStr|plural}}_{{.OutStr|plural}}(ins []*{{.InType}})(outs []*{{.OutType}}) {
	outs = make([]*{{.OutType}}, len(ins))
	for i := range outs {
    outs[i] = convert_{{.InStr}}_{{.OutStr}}(ins[i])
  }
  return outs
}
`

const tplConvertApplyText = `
func apply_{{.ArgStr}}(in *{{.BaseType}}, arg *{{.ArgType}})(out *{{.BaseType}}) {
	if in == nil {
		return nil
	}
	return &{{.BaseType}}{
  {{range .Fields}}
		{{.|fieldName}}: {{.|fieldApply}}, {{lastComment}}
  {{end}}
	}
}
`
