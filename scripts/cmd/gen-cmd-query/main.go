package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/types"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"k8s.io/code-generator/third_party/forked/golang/reflect"

	"golang.org/x/tools/go/packages"
)

var flV = flag.Bool("v", false, "verbose")
var flOut = flag.String("format-out", "types.d.go", "format of generated declaration Go files")

func usage() {
	const text = `
gen-cmd-query finds service definitions and generate code to dispatch queries to
corresponding method.

Usage: gen-cmd-query [OPTION] PACKAGE ...

Options:
`

	fmt.Print(text[1:])
	flag.PrintDefaults()
}

func init() {
	flag.Usage = usage
}

func debugf(format string, args ...interface{}) {
	if *flV {
		_, _ = fmt.Fprintf(os.Stderr, format, args...)
	}
}

func fatalf(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}

var storedErrors []string

func errorf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	storedErrors = append(storedErrors, msg)
}

func mustNoError(format string, args ...interface{}) {
	count := len(storedErrors)
	if count == 0 {
		return
	}

	_, _ = fmt.Fprintf(os.Stderr, format, args...)
	for _, msg := range storedErrors {
		msg = strings.TrimRight(msg, "\n")
		_, _ = fmt.Fprintf(os.Stderr, "    %v\n", msg)
	}
	switch {
	case count == 1:
		fatalf("stopped due to %v error\n", count)
	case count > 1:
		fatalf("stopped due to %v errors\n", count)
	}
}

func must(err error) {
	if err != nil {
		fatalf("%v\n", err)
	}
}

func mustWrite(w io.Writer, p []byte) {
	_, err := w.Write(p)
	must(err)
}

func formatFileName(format string, fileName string) string {
	if strings.Contains(format, "{}") {
		return strings.Replace(format, "{}", fileName, 1)
	}
	if strings.HasPrefix(format, ".") {
		return fileName + format
	}
	return format
}

func p(w io.Writer, format string, args ...interface{}) {
	_, err := fmt.Fprintf(w, format, args...)
	must(err)
}

const QueryService = "QueryService"
const Aggregate = "Aggregate"

type ServiceDef struct {
	Kind    string
	Package *packages.Package
	Ident   *ast.Ident
	Type    *types.Interface
}

type HandlerDef struct {
	Method  *types.Func
	Request *types.Struct

	RequestNamed  *types.Named
	ResponseNamed *types.Named
}

type MultiWriter struct {
	*Writer
	WriteArgs     bytes.Buffer
	WriteIface    bytes.Buffer
	WriteDispatch bytes.Buffer
}

func main() {
	flag.Parse()
	formatOut := *flOut
	if formatOut == "" {
		fatalf("invalid output format")
	}

	pkgPaths := flag.Args()
	if len(pkgPaths) == 0 {
		flag.Usage()
		os.Exit(2)
	}

	cfg := &packages.Config{Mode: packages.LoadAllSyntax}
	pkgs, err := packages.Load(cfg, pkgPaths...)
	must(err)

	if packages.PrintErrors(pkgs) > 0 {
		os.Exit(1)
	}

	var services []ServiceDef
	var generatedFiles []string
	kinds := []string{QueryService, Aggregate}
	for _, pkg := range pkgs {
		debugf("processing package %v\n", pkg.PkgPath)
		services = nil

		defs := pkg.TypesInfo.Defs
		for ident, obj := range defs {
			if obj != nil {
				debugf("    %v : %v", ident.Name, obj.String())
			} else {
				debugf("    %v : nil", ident.Name)
			}
			kind, svc, err := checkService(kinds, ident, obj)
			if err != nil {
				errorf("%v\n", err)
				continue
			}
			if svc != nil {
				services = append(services, ServiceDef{
					Kind:    kind,
					Package: pkg,
					Ident:   ident,
					Type:    svc,
				})
			}
		}
		mustNoError("package %v:\n", pkg.PkgPath)
		if len(services) == 0 {
			fmt.Printf("  skipped %v\n", pkg.PkgPath)
			continue
		}

		sort.Slice(services, func(i, j int) bool {
			return services[i].Ident.Name < services[j].Ident.Name
		})

		w := NewWriter(pkg.Name, pkg.PkgPath)
		ws := &MultiWriter{Writer: w}
		writeCommonDeclaration(ws)
		for _, item := range services {
			debugf("processing service %v\n", item.Ident.Name)
			processService(ws, item)
		}

		p(w, "\n// implement interfaces\n\n")
		mustWrite(w, ws.WriteIface.Bytes())
		p(w, "\n// implement conversion\n\n")
		mustWrite(w, ws.WriteArgs.Bytes())
		p(w, "\n// implement dispatching\n\n")
		mustWrite(w, ws.WriteDispatch.Bytes())

		fileName := formatFileName(formatOut, pkg.Name)
		dirPath := filepath.Dir(pkg.GoFiles[0])
		absFileName := filepath.Join(dirPath, fileName)
		w.WriteFile(absFileName, 0666)
		generatedFiles = append(generatedFiles, absFileName)
		fmt.Printf("generated %v\n", absFileName)
	}

	execGoimport(generatedFiles)
}

func checkService(
	suffixes []string, ident *ast.Ident, obj types.Object,
) (kind string, _ *types.Interface, _ error) {
	for _, suffix := range suffixes {
		if strings.HasSuffix(ident.Name, suffix) {
			kind = suffix
			break
		}
	}
	if kind == "" {
		return
	}

	if obj == nil {
		return "", nil, fmt.Errorf("%v: can not load definition", ident.Name)
	}
	typ := obj.Type()
	if typ == nil {
		return "", nil, fmt.Errorf("%v: can not load type information", ident.Name)
	}
	if typ, ok := typ.(*types.Named); ok {
		if typ, ok := typ.Underlying().(*types.Interface); ok {
			return kind, typ, nil
		}
	}
	return "", nil, fmt.Errorf("%v: must be an interface", ident.Name)
}

func processService(w *MultiWriter, def ServiceDef) {
	switch def.Kind {
	case QueryService:
		processQueryService(w, def.Package, def.Ident, def.Type)
	case Aggregate:
		processAggregate(w, def.Package, def.Ident, def.Type)
	default:
		panic("unexpected")
	}
}

func processQueryService(w *MultiWriter, pkg *packages.Package, ident *ast.Ident, typ *types.Interface) {
	defs := extractHandlerDefs(pkg, ident, typ)
	generateQueries(w, ident.Name, defs)
	mustNoError("type %v.%v:\n", pkg.PkgPath, ident.Name)
}

func processAggregate(w *MultiWriter, pkg *packages.Package, ident *ast.Ident, typ *types.Interface) {
	defs := extractHandlerDefs(pkg, ident, typ)
	generateCommands(w, ident.Name, defs)
	mustNoError("type %v.%v:\n", pkg.PkgPath, ident.Name)
}

func extractHandlerDefs(pkg *packages.Package, ident *ast.Ident, typ *types.Interface) (defs []HandlerDef) {
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
		req, _, err := checkMethodSignature(params, results)
		if err != nil {
			errorf("%v: %v\n", method.Name(), err)
			continue
		}
		defs = append(defs, HandlerDef{
			Method:  method,
			Request: req,

			RequestNamed:  params.At(1).Type().(*types.Pointer).Elem().(*types.Named),
			ResponseNamed: results.At(0).Type().(*types.Pointer).Elem().(*types.Named),
		})
	}
	mustNoError("type %v.%v:\n", pkg.PkgPath, ident.Name)
	return defs
}

func checkMethodSignature(params *types.Tuple, results *types.Tuple) (request *types.Struct, response *types.Struct, err error) {
	if params.Len() != 2 {
		err = errors.New("expect 2 params")
		return
	}
	if results.Len() != 2 {
		err = errors.New("expect 2 return values")
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
		t := results.At(1)
		if t.Type().String() != "error" {
			err = errors.New("expect the second return value is error")
			return
		}
	}
	{
		t := params.At(1)
		typ, _err := checkPtrStruct(t.Type())
		if _err != nil {
			err = fmt.Errorf("expect the second param is a pointer to struct (%v)", _err)
			return
		}
		request = typ
	}
	{
		t := results.At(0)
		typ, _err := checkPtrStruct(t.Type())
		if _err != nil {
			err = fmt.Errorf("expect the first return value is a pointer to struct (%v)", _err)
			return
		}
		response = typ
	}
	return
}

func checkPtrStruct(t types.Type) (*types.Struct, error) {
	ptr, ok := t.(*types.Pointer)
	if !ok {
		return nil, errors.New("must be explicit pointer (i.e. *Type)")
	}
	t = ptr.Elem()

	for {
		switch typ := t.(type) {
		case *types.Pointer:
			return nil, errors.New("got double pointer")

		case *types.Named:
			t = typ.Underlying()
			continue

		case *types.Struct:
			return typ, nil

		default:
			return nil, fmt.Errorf("got %v", typ)
		}
	}
}

func writeCommonDeclaration(w *MultiWriter) {
	w.Import("etop.vn/api/meta")

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
	w.Import("context")

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
		genRequestName := renderNamed(w.Writer, item.RequestNamed)
		genResponseName := renderNamed(w.Writer, item.ResponseNamed)

		p(w, "type %v struct {\n", genQueryName)
		generateStruct(w.Writer, item.Request)
		mustNoError("method %v:\n", item.Method)
		p(w, "\nResult *%v `json:\"-\"`\n", genResponseName)
		p(w, "}\n\n")

		// implement GetArgs()
		{
			w2 := &w.WriteArgs
			generateGetArgs(w2, genQueryName, genRequestName, item.Request)
		}
		// implement Query
		{
			w2 := &w.WriteIface
			p(w2, "func (q *%v) query() {}\n", genQueryName)
		}
		// implement Handle()
		{
			const tmpl = `
func (h %v) Handle%v(ctx context.Context, query *%v) error {
	result, err := h.inner.%v(ctx, query.GetArgs())
	query.Result = result
	return err
}
`
			w2 := &w.WriteDispatch
			p(w2, tmpl, genHandlerName, methodName, genQueryName, methodName)
		}
	}
}

func generateCommands(w *MultiWriter, serviceName string, defs []HandlerDef) {
	w.Import("context")

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
		genRequestName := renderNamed(w.Writer, item.RequestNamed)
		genResponseName := renderNamed(w.Writer, item.ResponseNamed)

		p(w, "type %v struct {\n", genCommandName)
		generateStruct(w.Writer, item.Request)
		mustNoError("method %v:\n", item.Method)
		p(w, "\nResult *%v `json:\"-\"`\n", genResponseName)
		p(w, "}\n\n")

		// implement GetArgs()
		{
			w2 := &w.WriteArgs
			generateGetArgs(w2, genCommandName, genRequestName, item.Request)
		}
		// implement Command
		{
			w2 := &w.WriteIface
			p(w2, "func (q *%v) command() {}\n", genCommandName)
		}
		// implement Handle()
		{
			const tmpl = `
func (h %v) Handle%v(ctx context.Context, cmd *%v) error {
	result, err := h.inner.%v(ctx, cmd.GetArgs())
	cmd.Result = result
	return err
}
`
			w2 := &w.WriteDispatch
			p(w2, tmpl, genHandlerName, methodName, genCommandName, methodName)
		}
	}
}

func generateGetArgs(w io.Writer, wrapperName, requestName string, request *types.Struct) {
	p(w, "func (q *%v) GetArgs() *%v {\n", wrapperName, requestName)
	p(w, "\treturn &%v{\n", requestName)
	for i, n := 0, request.NumFields(); i < n; i++ {
		field := request.Field(i)
		p(w, "\t\t%v: q.%v,\n", field.Name(), field.Name())
	}
	p(w, "\t}\n")
	p(w, "}\n")
}

func renderNamed(w Importer, named *types.Named) string {
	obj := named.Obj()
	pkgAlias := w.Import(obj.Pkg().Path())
	if pkgAlias == "" {
		return obj.Name()
	}
	return pkgAlias + "." + obj.Name()
}

var reTypeImport = regexp.MustCompile(`([0-9A-z/._-]+)\.([0-9A-z]+)`)

func renderType(w Importer, typ types.Type) (string, error) {
	s := typ.String()

	result := reTypeImport.ReplaceAllStringFunc(s,
		func(match string) string {
			parts := reTypeImport.FindStringSubmatch(match)
			if parts == nil {
				panic("unexpected")
			}

			pkgPath := parts[1]
			typeName := parts[2]
			pkgAlias := w.Import(pkgPath)
			if pkgAlias == "" {
				return typeName
			}
			return pkgAlias + "." + typeName
		})
	return result, nil
}

func generateStruct(w *Writer, s *types.Struct) {
	n := s.NumFields()
	for i := 0; i < n; i++ {
		field := s.Field(i)
		if !field.Exported() {
			continue
		}

		if !field.Embedded() {
			p(w, "%v ", field.Name())
		}

		processedTag, err := processTag(s.Tag(i))
		if err != nil {
			errorf("field %v: incorrect tag format (%v)\n", field.Name(), err)
			continue
		}

		ftyp := field.Type()
		renderedType, err := renderType(w, ftyp)
		if err != nil {
			errorf("field %v: %v\n", field.Name(), err)
			continue
		}
		p(w, "%v %v\n", renderedType, processedTag)
	}
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

func execGoimport(files []string) {
	args := []string{"-w"}
	args = append(args, files...)
	cmd := exec.Command("goimports", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fatalf("%s\n\n%s\n", err, out)
	}
	debugf("%s", out)
}
