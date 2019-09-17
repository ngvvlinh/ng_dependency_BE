package cq

import (
	"errors"
	"fmt"
	"go/types"
	"strings"

	"etop.vn/backend/tools/pkg/reflect"
)

func processService(w *MultiWriter, def ServiceDef) {
	switch def.Kind {
	case QueryService:
		processQueryService(w, def.PkgPath, def.Name, def.Type)
	case Aggregate:
		processAggregate(w, def.PkgPath, def.Name, def.Type)
	default:
		panic("unexpected")
	}
}

func processQueryService(w *MultiWriter, pkgPath string, name string, typ *types.Interface) {
	defs := extractHandlerDefs(pkgPath, name, typ)
	generateQueries(w, name, defs)
	mustNoError("type %v.%v:\n", pkgPath, name)
}

func processAggregate(w *MultiWriter, pkgPath string, name string, typ *types.Interface) {
	defs := extractHandlerDefs(pkgPath, name, typ)
	generateCommands(w, name, defs)
	mustNoError("type %v.%v:\n", pkgPath, name)
}

func extractHandlerDefs(pkgPath string, name string, typ *types.Interface) (defs []HandlerDef) {
	n := typ.NumMethods()
	for i := 0; i < n; i++ {
		method := typ.Method(i)
		if !method.Exported() {
			continue
		}

		mtyp := method.Type()
		styp := mtyp.(*types.Signature)
		params := styp.Params()
		results := styp.Results()
		requests, responses, err := checkMethodSignature(method.Name(), params, results)
		if err != nil {
			errorf("%v: %v", method.Name(), err)
			continue
		}
		defs = append(defs, HandlerDef{
			Method:    method,
			Requests:  requests,
			Responses: responses,
		})
	}
	mustNoError("type %v.%v:\n", pkgPath, name)
	return defs
}

func checkMethodSignature(name string, params *types.Tuple, results *types.Tuple) (requests []*ArgItem, responses []*ArgItem, err error) {
	if params.Len() == 0 {
		err = errors.New("expect at least 1 param")
		return
	}
	if results.Len() == 0 {
		err = errors.New("expect at least 1 param")
		return
	}
	{
		t := params.At(0)
		if t.Type().String() != "context.Context" {
			err = errors.New("expect the first param is context.Context")
			return
		}
	}
	{
		t := results.At(results.Len() - 1)
		if t.Type().String() != "error" {
			err = errors.New("expect the last return value is error")
			return
		}
	}
	{
		// skip the first param (context.Context)
		for i, n := 1, params.Len(); i < n; i++ {
			arg, err := checkArg(params.At(i), n == 2)
			if err != nil {
				errorf("%v: %v", name, err)
			}
			requests = append(requests, arg)
			if !arg.Inline && arg.Name == "" {
				errorf("%v: must provide name for param %v", name, arg.Type)
			}
		}
	}
	{
		// skip the last result (error)
		for i, n := 0, results.Len()-1; i < n; i++ {
			arg, err := checkArg(results.At(i), n == 2)
			if err != nil {
				errorf("%v: %v", name, err)
			}
			responses = append(responses, arg)
		}
		if len(responses) > 1 {
			for _, arg := range responses {
				if arg.Name == "" || strings.HasPrefix(arg.Name, "_") {
					errorf("%v: must provide name for result %v", name, arg.Type)
				}
			}
		}
	}
	return
}

func checkArg(v *types.Var, autoInline bool) (*ArgItem, error) {
	arg := &ArgItem{
		Inline: v.Name() == "_" || v.Name() == "" && autoInline,
		Name:   toTitle(v.Name()),
		Var:    v,
		Type:   v.Type(),
	}
	// when inline, the param must be struct or pointer to struct
	if arg.Inline {
		var err error
		arg.Struct, arg.Ptr, err = checkStruct(v.Type())
		if err != nil {
			return nil, fmt.Errorf("type must be a struct or a pointer to struct to be inline: %v", err)
		}
	}
	return arg, nil
}

func checkStruct(t types.Type) (_ *types.Struct, ptr bool, _ error) {
	p, ptr := t.(*types.Pointer)
	if ptr {
		t = p.Elem()
	}

underlying:
	switch typ := t.(type) {
	case *types.Pointer:
		return nil, false, fmt.Errorf("got double pointer (%v)", t)

	case *types.Named:
		t = typ.Underlying()
		goto underlying

	case *types.Struct:
		return typ, ptr, nil

	default:
		return nil, false, fmt.Errorf("got %v", typ)
	}
}

func writeCommonDeclaration(w *MultiWriter) {
	w.Import("meta", "etop.vn/api/meta")

	tmpl := `
type Command interface { command() }
type Query interface { query() }
type CommandBus struct { bus meta.Bus }
type QueryBus struct { bus meta.Bus }

func (c CommandBus) Dispatch(ctx context.Context, msg Command) error {
	return c.bus.Dispatch(ctx, msg)
}
func (c QueryBus) Dispatch(ctx context.Context, msg Query) error {
	return c.bus.Dispatch(ctx, msg)
}
func (c CommandBus) DispatchAll(ctx context.Context, msgs ...Command) error {
	for _, msg := range msgs {
		if err := c.bus.Dispatch(ctx, msg); err != nil {
			return err
		}
	}
	return nil
}
func (c QueryBus) DispatchAll(ctx context.Context, msgs ...Query) error {
	for _, msg := range msgs {
		if err := c.bus.Dispatch(ctx, msg); err != nil {
			return err
		}
	}
	return nil
}
`
	mustWrite(w, []byte(tmpl))
}

func generateQueries(w *MultiWriter, serviceName string, defs []HandlerDef) {
	w.Import("context", "context")

	genHandlerName := serviceName + "Handler"
	{
		tmpl := `
type %v struct {
	inner %v
}

func New%v(service %v) %v { return %v{service} }

func (h %v) RegisterHandlers(b interface{
	meta.Bus
	AddHandler(handler interface{})
}) QueryBus {
`
		w2 := &w.WriteDispatch
		p(w2, tmpl, genHandlerName, serviceName, genHandlerName, serviceName, genHandlerName, genHandlerName, genHandlerName)
		for _, item := range defs {
			p(w2, "\tb.AddHandler(h.Handle%v)\n", item.Method.Name())
		}
		p(w2, "\treturn QueryBus{b}\n")
		p(w2, "}\n")
	}

	for _, item := range defs {
		methodName := item.Method.Name()
		genQueryName := item.Method.Name() + "Query"

		// generate declaration
		{
			p(w, "type %v struct {\n", genQueryName)
			generateStruct(w.Writer, item.Requests)
			mustNoError("method %v:\n", item.Method)
			generateResult(w, item)
			p(w, "}\n\n")
		}
		// implement Handle()
		{
			generateHandle(w, item, methodName, genHandlerName, genQueryName)
		}
		// implement GetArgs()
		{
			w2 := w.GetImportWriter(&w.WriteArgs)
			generateGetArgs(w2, genQueryName, item.Requests)
		}
		// implement Query
		{
			w2 := &w.WriteIface
			p(w2, "func (q *%v) query() {}\n", genQueryName)
		}

	}
}

func generateCommands(w *MultiWriter, serviceName string, defs []HandlerDef) {
	w.Import("context", "context")

	genHandlerName := serviceName + "Handler"
	{
		tmpl := `
type %v struct {
	inner %v
}

func New%v(service %v) %v { return %v{service} }

func (h %v) RegisterHandlers(b interface{
	meta.Bus
	AddHandler(handler interface{})
}) CommandBus {
`
		w2 := &w.WriteDispatch
		p(w2, tmpl, genHandlerName, serviceName, genHandlerName, serviceName, genHandlerName, genHandlerName, genHandlerName)
		for _, item := range defs {
			p(w2, "\tb.AddHandler(h.Handle%v)\n", item.Method.Name())
		}
		p(w2, "\treturn CommandBus{b}")
		p(w2, "}\n")
	}

	for _, item := range defs {
		methodName := item.Method.Name()
		genCommandName := item.Method.Name() + "Command"

		p(w, "type %v struct {\n", genCommandName)
		generateStruct(w.Writer, item.Requests)
		mustNoError("method %v:\n", item.Method)
		generateResult(w, item)
		p(w, "}\n\n")

		// implement GetArgs()
		{
			w2 := w.GetImportWriter(&w.WriteArgs)
			generateGetArgs(w2, genCommandName, item.Requests)
		}
		// implement Handle()
		{
			generateHandle(w, item, methodName, genHandlerName, genCommandName)
		}
		// implement Command
		{
			w2 := &w.WriteIface
			p(w2, "func (q *%v) command() {}\n", genCommandName)
		}
	}
}

func generateGetArgs(w ImportWriter, wrapperName string, requests ArgItems) {
	p(w, "func (q *%v) GetArgs(ctx context.Context) (_ context.Context, ", wrapperName)
	generateArgList(w, requests)
	p(w, ") {\n")
	p(w, "\treturn ctx,\n")

	comma := false
	inline := false
	err := requests.Walk(
		func(node NodeType, name string, field *types.Var, tag string) error {
			if comma {
				p(w, ",\n")
				comma = false
			}

			switch node {
			case NodeStartInline:
				inline = true
				p(w, "%v{\n", renderType(w, field.Type(), true))

			case NodeEndInline:
				inline = false
				p(w, "}\n")
				comma = true

			case NodeField:
				if inline {
					p(w, "\t%v: q.%v", name, name)
				} else {
					p(w, "q.%v", name)
				}
				comma = true

			default:
				panic("unexpected")
			}
			return nil
		})
	must(err)
	p(w, "}\n\n")
}

func generateArgList(w ImportWriter, args []*ArgItem) {
	for i, arg := range args {
		if i > 0 {
			p(w, ", ")
		}
		name := arg.Var.Name()
		if name == "" {
			name = "_"
		}
		p(w, "%v %v", name, renderType(w, arg.Type, false))
	}
}

func renderType(w Importer, typ types.Type, literal bool) string {
	result := w.TypeString(typ)
	if literal && result[0] == '*' {
		result = "&" + result[1:]
	}
	return result
}

func generateStruct(w ImportWriter, args ArgItems) {
	err := args.Walk(
		func(node NodeType, name string, field *types.Var, tag string) error {
			switch node {
			case NodeField:
				processedTag, err := processTag(tag)
				if err != nil {
					errorf("field %v: incorrect tag format (%v)\n", field.Name(), err)
					return nil
				}
				p(w, "%v %v %v\n", name, renderType(w, field.Type(), false), processedTag)
			}
			return nil
		})
	must(err)
}

func generateResult(w ImportWriter, item HandlerDef) {
	if len(item.Responses) == 1 {
		p(w, "\nResult %v `json:\"-\"`\n", renderType(w, item.Responses[0].Type, false))
	} else {
		p(w, "\nResult struct {\n")
		for _, arg := range item.Responses {
			p(w, "%v %v\n", arg.Name, renderType(w, arg.Type, false))
		}
		p(w, "} `json:\"-\"`\n")
	}
}

func generateHandle(w ImportWriter, item HandlerDef, methodName, genHandlerName, genQueryName string) {
	p(w, "\nfunc (h %v) Handle%v(ctx context.Context, msg *%v) (err error) {\n", genHandlerName, methodName, genQueryName)
	switch len(item.Responses) {
	case 0:
		p(w, "return h.inner.%v(msg.GetArgs(ctx))\n", methodName)
	case 1:
		p(w, "msg.Result, err = h.inner.%v(msg.GetArgs(ctx))\n", methodName)
		p(w, "return err\n")
	default:
		for _, arg := range item.Responses {
			p(w, "msg.Result.%v, ", arg.Var.Name())
		}
		p(w, "err = h.inner.%v(msg.GetArgs(ctx))\n", methodName)
		p(w, "return err")
	}
	p(w, "}\n")
}

func processTag(tag string) (string, error) {
	stag, err := reflect.ParseStructTags(tag)
	if err != nil {
		return "", err
	}
	if strings.Contains(tag, "`") {
		return "", errors.New("backquote (`) is not supported in tag")
	}

	result := make(reflect.StructTags, 0, len(stag))
	for _, t := range stag {
		if t.Name != "protobuf" {
			result = append(result, t)
		}
	}
	if len(result) == 0 {
		return "", nil
	}
	return result.String(), nil
}
