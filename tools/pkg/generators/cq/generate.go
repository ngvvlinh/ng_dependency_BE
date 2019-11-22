package cq

import (
	"fmt"
	"go/types"
	"strings"

	"etop.vn/backend/tools/pkg/generator"
	"etop.vn/backend/tools/pkg/genutil"
)

func processService(w *MultiWriter, ng generator.Engine, def ServiceDef) {
	switch def.Kind {
	case QueryService:
		processQueryService(w, ng, def.PkgPath, def.Name, def.Type)
	case Aggregate:
		processAggregate(w, ng, def.PkgPath, def.Name, def.Type)
	default:
		panic("unexpected")
	}
}

func processQueryService(w *MultiWriter, ng generator.Engine, pkgPath string, name string, typ *types.Interface) {
	defs := ExtractHandlerDefs(ng, pkgPath, name, typ)
	generateQueries(w, name, defs)
	mustNoError("type %v.%v:\n", pkgPath, name)
}

func processAggregate(w *MultiWriter, ng generator.Engine, pkgPath string, name string, typ *types.Interface) {
	defs := ExtractHandlerDefs(ng, pkgPath, name, typ)
	generateCommands(w, name, defs)
	mustNoError("type %v.%v:\n", pkgPath, name)
}

func ExtractHandlerDefs(ng generator.Engine, pkgPath string, name string, typ *types.Interface) (defs []*HandlerDef) {
	n := typ.NumMethods()
	for i := 0; i < n; i++ {
		method := typ.Method(i)
		if !method.Exported() {
			continue
		}
		def, err := ExtractHandlerDef(ng, method)
		if err != nil {
			errorf("%v", err)
			continue
		}
		defs = append(defs, def)
	}
	mustNoError("type %v.%v:\n", pkgPath, name)
	return defs
}

func ExtractHandlerDef(ng generator.Engine, method *types.Func) (*HandlerDef, error) {
	mtyp := method.Type()
	styp := mtyp.(*types.Signature)
	params := styp.Params()
	results := styp.Results()
	requests, responses, err := CheckMethodSignature(method.Name(), params, results)
	if err != nil {
		return nil, fmt.Errorf("%v: %v", method.Name(), err)
	}
	return &HandlerDef{
		Name:     method.Name(),
		Comment:  ng.GetComment(method).Text(),
		Method:   method,
		Request:  requests,
		Response: responses,
	}, nil
}

func CheckMethodSignature(name string, params *types.Tuple, results *types.Tuple) (request, response *Message, err error) {
	if params.Len() == 0 {
		err = generator.Errorf(nil, "expect at least 1 param")
		return
	}
	if results.Len() == 0 {
		err = generator.Errorf(nil, "expect at least 1 param")
		return
	}
	var requestItems, responseItems []*ArgItem
	{
		t := params.At(0)
		if t.Type().String() != "context.Context" {
			err = generator.Errorf(nil, "expect the first param is context.Context")
			return
		}
	}
	{
		t := results.At(results.Len() - 1)
		if t.Type().String() != "error" {
			err = generator.Errorf(nil, "expect the last return value is error")
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
			requestItems = append(requestItems, arg)
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
			responseItems = append(responseItems, arg)
		}
		if len(responseItems) > 1 {
			for _, arg := range responseItems {
				if arg.Name == "" || strings.HasPrefix(arg.Name, "_") {
					errorf("%v: must provide name for result %v", name, arg.Type)
				}
			}
		}
	}
	request = &Message{Items: requestItems}
	response = &Message{Items: responseItems}
	return request, response, nil
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
	w.Import("capi", "etop.vn/capi")

	tmpl := `
type Command interface { command() }
type Query interface { query() }
type CommandBus struct { bus capi.Bus }
type QueryBus struct { bus capi.Bus }

func NewCommandBus(bus capi.Bus) CommandBus                          { return CommandBus{bus} }
func NewQueryBus(bus capi.Bus) QueryBus                              { return QueryBus{bus} }
func (c CommandBus) Dispatch(ctx context.Context, msg Command) error { return c.bus.Dispatch(ctx, msg) }
func (c QueryBus) Dispatch(ctx context.Context, msg Query) error     { return c.bus.Dispatch(ctx, msg) }
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

func generateQueries(w *MultiWriter, serviceName string, defs []*HandlerDef) {
	w.Import("context", "context")

	genHandlerName := serviceName + "Handler"
	{
		tmpl := `
type %v struct {
	inner %v
}

func New%v(service %v) %v { return %v{service} }

func (h %v) RegisterHandlers(b interface{
	capi.Bus
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
			generateStruct(w.Writer, item.Request)
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
			generateGetArgs(w2, genQueryName, item.Request)
			generateSetArgs(w2, genQueryName, item.Request)
		}
		// implement Query
		{
			w2 := &w.WriteIface
			p(w2, "func (q *%v) query() {}\n", genQueryName)
		}

	}
}

func generateCommands(w *MultiWriter, serviceName string, defs []*HandlerDef) {
	w.Import("context", "context")

	genHandlerName := serviceName + "Handler"
	{
		tmpl := `
type %v struct {
	inner %v
}

func New%v(service %v) %v { return %v{service} }

func (h %v) RegisterHandlers(b interface{
	capi.Bus
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
		generateStruct(w.Writer, item.Request)
		mustNoError("method %v:\n", item.Method)
		generateResult(w, item)
		p(w, "}\n\n")

		// implement GetArgs()
		{
			w2 := w.GetImportWriter(&w.WriteArgs)
			generateGetArgs(w2, genCommandName, item.Request)
			generateSetArgs(w2, genCommandName, item.Request)
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

func generateGetArgs(w ImportWriter, wrapperName string, requests *Message) {
	p(w, "func (q *%v) GetArgs(ctx context.Context) (_ context.Context, ", wrapperName)
	generateArgList(w, requests.Items)
	p(w, ") {\n")
	p(w, "\treturn ctx,\n")

	comma := false
	inline := false
	err := requests.Items.Walk(
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

func generateSetArgs(w ImportWriter, wrapperName string, requests *Message) {
	for _, req := range requests.Items {
		if !req.Inline {
			continue
		}
		p(w, "func (q *%v) Set%v(args %v) {\n", wrapperName, renderTypeName(req.Type), renderType(w, req.Type, false))
		for i, n := 0, req.Struct.NumFields(); i < n; i++ {
			field := req.Struct.Field(i)
			p(w, "q.%v = args.%v\n", field.Name(), field.Name())
		}
		p(w, "}\n\n")
	}
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

func renderTypeName(typ types.Type) string {
	ptr, ok := typ.(*types.Pointer)
	if ok {
		typ = ptr.Elem()
	}
	return typ.(*types.Named).Obj().Name()
}

func generateStruct(w ImportWriter, args *Message) {
	err := args.Items.Walk(
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

func generateResult(w ImportWriter, item *HandlerDef) {
	items := item.Response.Items
	if len(items) == 1 {
		p(w, "\nResult %v `json:\"-\"`\n", renderType(w, items[0].Type, false))
	} else {
		p(w, "\nResult struct {\n")
		for _, arg := range items {
			p(w, "%v %v\n", arg.Name, renderType(w, arg.Type, false))
		}
		p(w, "} `json:\"-\"`\n")
	}
}

func generateHandle(w ImportWriter, item *HandlerDef, methodName, genHandlerName, genQueryName string) {
	p(w, "\nfunc (h %v) Handle%v(ctx context.Context, msg *%v) (err error) {\n", genHandlerName, methodName, genQueryName)
	switch len(item.Response.Items) {
	case 0:
		p(w, "return h.inner.%v(msg.GetArgs(ctx))\n", methodName)
	case 1:
		p(w, "msg.Result, err = h.inner.%v(msg.GetArgs(ctx))\n", methodName)
		p(w, "return err\n")
	default:
		for _, arg := range item.Response.Items {
			p(w, "msg.Result.%v, ", arg.Var.Name())
		}
		p(w, "err = h.inner.%v(msg.GetArgs(ctx))\n", methodName)
		p(w, "return err")
	}
	p(w, "}\n")
}

func processTag(tag string) (string, error) {
	stag, err := genutil.ParseStructTags(tag)
	if err != nil {
		return "", err
	}
	if strings.Contains(tag, "`") {
		return "", generator.Errorf(nil, "backquote (`) is not supported in tag")
	}

	result := make(genutil.StructTags, 0, len(stag))
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
