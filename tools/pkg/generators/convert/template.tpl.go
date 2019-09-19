package convert

const tplConvertTypeText = `
func Convert_{{.InStr}}_{{.OutStr}}(in *{{.InType}}, out *{{.OutType}}) *{{.OutType}} {
  {{- if .CustomConversionMode|eq 1}}
    return {{.CustomConversionFuncType}}(in)
  {{- else if .CustomConversionMode|eq 2}}
    if in == nil {
        return nil
    }
    if out == nil {
        out = &{{.OutType}}{}
    }
    {{.CustomConversionFuncType}}(in, out)
    return out
  {{- else if .CustomConversionMode|eq 3}}
    return {{.CustomConversionFuncType}}(in, out)
  {{- else}}
    if in == nil {
        return nil
    }
    if out == nil {
      out = &{{.OutType}}{}
    }
    convert_{{.InStr}}_{{.OutStr}}(in, out)
    return out
  {{- end}}
}

func convert_{{.InStr}}_{{.OutStr}}(in *{{.InType}}, out *{{.OutType}}) {
	{{- range .Fields}}
		out.{{.|fieldName}} = {{.|fieldValue "in"}} {{lastComment -}}
  {{end}}
}

func Convert_{{.InStr|plural}}_{{.OutStr|plural}}(ins []*{{.InType}})(outs []*{{.OutType}}) {
  tmps := make([]{{.OutType}}, len(ins))
  outs = make([]*{{.OutType}}, len(ins))
	for i := range tmps {
		outs[i] = Convert_{{.InStr}}_{{.OutStr}}(ins[i], &tmps[i])
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
func apply_{{.ArgStr}}(arg *{{.ArgType}}, out *{{.BaseType}}) {
  {{- range .Fields}}
	out.{{.|fieldName}} = {{.|fieldApply "arg"}} {{lastComment -}}
  {{end}}
}
`
