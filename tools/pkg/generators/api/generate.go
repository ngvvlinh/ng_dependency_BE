package api

import (
	"go/types"
	"strings"
	"text/template"

	"o.o/backend/tools/pkg/generator"
	"o.o/backend/tools/pkg/generators/api/defs"
	"o.o/backend/tools/pkg/generators/api/parse"
)

var meta = make(parse.Meta)
var currentPrinter generator.Printer
var tpl = template.Must(template.New("template").Funcs(funcs).Parse(tplText))

var funcs = map[string]interface{}{
	"busName":         renderBusName,
	"generateGetArgs": generateGetArgs,
	"generateHandle":  generateHandle,
	"generateResult":  generateResult,
	"generateSetArgs": generateSetArgs,
	"generateStruct":  generateStruct,
	"interfaceMethod": renderInterfaceMethod,
	"messageName":     renderMessageName,
}

func generate(w generator.Printer, services []*defs.Service) {
	w.Import("context", "context")
	w.Import("capi", "o.o/capi")
	vars := map[string]interface{}{
		"Services": services,
	}
	err := tpl.Execute(w, vars)
	must(err)
}

func generateGetArgs(m *defs.Method) string {
	var b strings.Builder
	p(&b, "func (q *%v) GetArgs(ctx context.Context) (_ context.Context, ", renderMessageName(m))
	p(&b, generateArgList(m.Request.Items))
	p(&b, ") {\n")
	p(&b, "\treturn ctx,\n")

	comma := false
	inline := false
	err := m.Request.Items.Walk(
		func(node defs.NodeType, name string, field *types.Var, tag string) error {
			if comma {
				p(&b, ",\n")
				comma = false
			}

			switch node {
			case defs.NodeStartInline:
				inline = true
				p(&b, "%v{\n", renderNew(field.Type()))

			case defs.NodeEndInline:
				inline = false
				p(&b, "}\n")
				comma = true

			case defs.NodeField:
				if inline {
					p(&b, "\t%v: q.%v", name, name)
				} else {
					p(&b, "q.%v", name)
				}
				comma = true

			default:
				panic("unexpected")
			}
			return nil
		})
	must(err)
	p(&b, "}\n\n")
	return b.String()
}

func generateSetArgs(m *defs.Method) string {
	var b strings.Builder
	for _, req := range m.Request.Items {
		if !req.Inline {
			continue
		}
		p(&b, "func (q *%v) Set%v(args %v) {\n", renderMessageName(m), renderTypeName(req.Type), renderType(req.Type))
		for i, n := 0, req.Struct.NumFields(); i < n; i++ {
			field := req.Struct.Field(i)
			p(&b, "q.%v = args.%v\n", field.Name(), field.Name())
		}
		p(&b, "}\n\n")
	}
	return b.String()
}

func generateArgList(args []*defs.ArgItem) string {
	var b strings.Builder
	for i, arg := range args {
		if i > 0 {
			p(&b, ", ")
		}
		name := arg.Var.Name()
		if name == "" {
			name = "_"
		}
		p(&b, "%v %v", name, renderType(arg.Type))
	}
	return b.String()
}

type keyBusName struct{}

func renderBusName(s *defs.Service) string {
	return cache(s, keyBusName{}, func() interface{} {
		switch s.Kind {
		case defs.KindQuery:
			return s.Name + "QueryBus"
		case defs.KindAggregate:
			return s.Name + "CommandBus"
		case defs.KindService:
			return s.Name + "Bus"
		default:
			return "<invalid>"
		}
	}).(string)
}

type keyInterfaceMethod struct{}

func renderInterfaceMethod(s *defs.Service) string {
	return cache(s, keyInterfaceMethod{}, func() interface{} {
		switch s.Kind {
		case defs.KindQuery:
			return "query" + s.Name
		case defs.KindAggregate:
			return "command" + s.Name
		case defs.KindService:
			return "request" + s.Name
		default:
			return "<invalid>"
		}
	}).(string)
}

type keyMessageName struct{}

func renderMessageName(m *defs.Method) string {
	return cache(m, keyMessageName{}, func() interface{} {
		switch m.Service.Kind {
		case defs.KindQuery:
			return m.Name + "Query"
		case defs.KindAggregate:
			return m.Name + "Command"
		case defs.KindService:
			return m.Name + "Request"
		default:
			return "<invalid>"
		}
	}).(string)
}

func renderType(typ types.Type) string {
	return currentPrinter.TypeString(typ)
}

func renderNew(typ types.Type) string {
	if ptr, ok := typ.(*types.Pointer); ok {
		return "&" + currentPrinter.TypeString(ptr.Elem())
	}
	return currentPrinter.TypeString(typ)
}

func renderTypeName(typ types.Type) string {
	ptr, ok := typ.(*types.Pointer)
	if ok {
		typ = ptr.Elem()
	}
	return typ.(*types.Named).Obj().Name()
}

func generateStruct(m *defs.Method) string {
	var b strings.Builder
	err := m.Request.Items.Walk(
		func(node defs.NodeType, name string, field *types.Var, tag string) error {
			switch node {
			case defs.NodeField:
				processedTag, err := processTag(tag)
				if err != nil {
					errorf("field %v: incorrect tag format (%v)\n", field.Name(), err)
					return nil
				}
				p(&b, "%v %v %v\n", name, renderType(field.Type()), processedTag)
			}
			return nil
		})
	must(err)
	return b.String()
}

func generateResult(m *defs.Method) string {
	var b strings.Builder
	items := m.Response.Items
	if len(items) == 1 {
		p(&b, "\nResult %v `json:\"-\"`\n", renderType(items[0].Type))
	} else {
		p(&b, "\nResult struct {\n")
		for _, arg := range items {
			p(&b, "%v %v\n", arg.Name, renderType(arg.Type))
		}
		p(&b, "} `json:\"-\"`\n")
	}
	return b.String()
}

func generateHandle(m *defs.Method) string {
	var b strings.Builder
	p(&b, "\nfunc (h %vHandler) Handle%v(ctx context.Context, msg *%v) (err error) {\n", m.Service.FullName, m.Name, renderMessageName(m))
	switch len(m.Response.Items) {
	case 0:
		p(&b, "return h.inner.%v(msg.GetArgs(ctx))\n", m.Name)
	case 1:
		p(&b, "msg.Result, err = h.inner.%v(msg.GetArgs(ctx))\n", m.Name)
		p(&b, "return err\n")
	default:
		for _, arg := range m.Response.Items {
			p(&b, "msg.Result.%v, ", arg.Var.Name())
		}
		p(&b, "err = h.inner.%v(msg.GetArgs(ctx))\n", m.Name)
		p(&b, "return err")
	}
	p(&b, "}\n")
	return b.String()
}

func processTag(tag string) (string, error) {
	if tag == "" {
		return "", nil
	}
	return "`" + tag + "`", nil
}

func cache(item, key interface{}, fn func() interface{}) interface{} {
	return meta.Cache(item, key, fn)
}
