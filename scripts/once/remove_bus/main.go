package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"io/ioutil"
	"path/filepath"
	"sort"

	"golang.org/x/tools/go/packages"
)

const ppBus = "o.o/backend/pkg/common/bus"

var cfg = &packages.Config{
	Fset: token.NewFileSet(),
	Mode: packages.LoadAllSyntax,
}

var flReplace = flag.Bool("replace", false, "overwrite the file")
var flWrite = flag.Bool("write", false, "write file")
var mapPkg = map[string]*packages.Package{}

func main() {
	flag.Parse()
	patterns := flag.Args()
	if len(patterns) == 0 {
		panic("no dir")
	}

	pkgs, err := packages.Load(cfg, patterns...)
	must(err)
	for _, pkg := range pkgs {
		mapPkg[pkg.PkgPath] = pkg
	}
	processPkgs(pkgs)
}

func processPkgs(pkgs []*packages.Package) {
	astStructs := map[*types.Named]*ast.StructType{}
	var calls []*CallExpr

	for _, pkg := range pkgs {
		busPkg := pkg.Imports[ppBus]
		if busPkg == nil {
			continue
		}

		for ident, item := range pkg.TypesInfo.Defs {
			if item == nil {
				continue
			}
			typ, ok := item.(*types.TypeName)
			if !ok {
				continue
			}
			named := namedStruct(typ.Type())
			if named == nil {
				continue
			}

			astStructs[named] = nil
			fmt.Println("detected struct", ident.Name)
		}

		calls = append(calls, lookupDispatch(pkg)...)
	}
	if !lookupOK {
		panic("*")
	}

	// detect the commands
	mapMsgType := map[*types.Named][]*MessageType{}
	for _, call := range calls {
		pkg := call.Pkg
		arg := call.Call.Args[1]
		switch arg := arg.(type) {
		case *ast.Ident:
			msgType := unwrapIdentType(pkg, arg)
			call.Msg = msgType
			mapMsgType[msgType] = nil

		default:
			fmt.Printf("unexpected call=%v arg=%v (%T)\n", call.Call, arg, arg)
		}
	}

	// detect the struct and method
	var _context, _error types.Type
	for st := range astStructs {
		for i, n := 0, st.NumMethods(); i < n; i++ {
			method := st.Method(i)
			sign := method.Type().(*types.Signature)
			if sign.Params().Len() != 2 {
				continue
			}
			if sign.Results().Len() != 1 {
				continue
			}
			param0 := sign.Params().At(0).Type()
			param1 := sign.Params().At(1).Type()
			result := sign.Results().At(0).Type()

			// cache the types if found
			if _context == nil && param0.String() == "context.Context" {
				_context = param0
			}
			if _error == nil && result.String() == "error" {
				_error = result
			}

			// lookup message
			if _context != param0 || _error != result {
				continue
			}
			msg := ptrNamedStruct(param1)
			if msg == nil {
				continue
			}
			if _, ok := mapMsgType[msg]; !ok {
				continue
			}

			msgType := &MessageType{
				Message: msg,
				Struct:  st,
				Method:  method,
			}
			mapMsgType[msg] = append(mapMsgType[msg], msgType)
		}
	}

	// lookup struct position
	for _, pkg := range pkgs {
		for _, file := range pkg.Syntax {
			ast.Inspect(file, func(node ast.Node) bool {
				switch node := node.(type) {
				case *ast.TypeSpec:
					st, ok := node.Type.(*ast.StructType)
					if !ok {
						return false
					}
					def := pkg.TypesInfo.Defs[node.Name]
					named := def.Type().(*types.Named)
					if _, ok = astStructs[named]; ok {
						astStructs[named] = st
					}
				}
				return true
			})
		}
	}

	// verify that all messages were found
	valid := true
	for msg := range mapMsgType {
		typs := mapMsgType[msg]
		if len(typs) == 0 {
			fmt.Printf("msg %v does not have handler\n", msg)
			valid = false
		}

		// select only exported method
		var choices []*MessageType
		for _, msgType := range typs {
			if msgType.Method.Exported() {
				choices = append(choices, msgType)
			}
		}
		if len(choices) == 1 {
			mapMsgType[msg] = choices
			continue
		}
		valid = false

		if len(choices) == 0 {
			fmt.Printf("found no exported methods for msg %v:\n", msg)
		} else {
			fmt.Printf("found multiple exported methods for msg %v:\n", msg)
		}
		for _, msgType := range typs {
			fmt.Printf("  %v", msgType.Method)
		}
	}
	if !valid {
		return
	}

	usedMessageTypes := map[MessageType]bool{}
	var replaces []Replacement
	for _, call := range calls {
		msgType := mapMsgType[call.Msg][0]
		sign := call.InMethod.Type().(*types.Signature)
		recv := sign.Recv().Name()
		usedMessageTypes[*msgType] = true

		replaces = append(replaces, includeInterface(astStructs, call.InStruct, msgType.Struct)...)
		newCode := fmt.Sprintf("%v.%v.%v", recv, msgType.Struct.Obj().Name(), msgType.Method.Name())
		replaces = append(replaces, replaceNode(call.Sel, newCode))
	}
	replaces = append(replaces, addInterfaces(usedMessageTypes)...)
	replaceFiles(cfg.Fset, replaces)
}

type keyStructStruct struct {
	a, b *types.Named
}

var mapStructStruct = map[keyStructStruct]bool{} // add receiver struct
var mapIncludeStruct = map[*types.Named]bool{}

func includeInterface(astStructs map[*types.Named]*ast.StructType, inStruct, msgReceiver *types.Named) (result []Replacement) {
	k := keyStructStruct{inStruct, msgReceiver}
	if mapStructStruct[k] {
		return nil
	}
	mapStructStruct[k] = true

	stPos := astStructs[inStruct].Fields.Closing

	// include a todo message
	if !mapIncludeStruct[inStruct] {
		mapIncludeStruct[inStruct] = true
		replace := Replacement{
			Start:      stPos,
			End:        stPos,
			ReplacedBy: "\n\t// TODO(vu): fix wire\n",
			Priority:   -10,
		}
		result = append(result, replace)
	}

	name := msgReceiver.Obj().Name()
	pkgName := msgReceiver.Obj().Pkg().Name()
	newCode := fmt.Sprintf("\t%v %v.%vInterface\n", name, pkgName, name)
	replace := Replacement{
		Start:      stPos,
		End:        stPos,
		ReplacedBy: newCode,
	}
	result = append(result, replace)
	return
}

type CallExpr struct {
	Pkg  *packages.Package
	Call *ast.CallExpr
	Sel  *ast.SelectorExpr
	Msg  *types.Named

	InMethod *types.Func
	InStruct *types.Named
}

func ptrNamedStruct(typ types.Type) *types.Named {
	ptr, ok := typ.(*types.Pointer)
	if !ok {
		return namedStruct(ptr)
	}
	return namedStruct(ptr.Elem())
}

func namedStruct(typ types.Type) *types.Named {
	named, ok := typ.(*types.Named)
	if !ok {
		return nil
	}
	_, ok = named.Underlying().(*types.Struct)
	if !ok {
		return nil
	}
	return named
}

func lookupDispatch(pkg *packages.Package) (calls []*CallExpr) {
	for _, file := range pkg.Syntax {

		var decl *ast.FuncDecl // for retriving scope
		ast.Inspect(file, func(node ast.Node) bool {
			switch node := node.(type) {
			case *ast.FuncDecl:
				decl = node

			case *ast.CallExpr:
				switch sel := node.Fun.(type) {
				case *ast.SelectorExpr:
					if sel.Sel.Name == "Dispatch" {
						switch x := sel.X.(type) {
						case (*ast.Ident):
							if x.Name == "bus" && x.Obj == nil {
								callExpr := &CallExpr{
									Pkg:  pkg,
									Call: node,
									Sel:  sel,
								}
								callExpr.InMethod, callExpr.InStruct = lookupMethodStruct(pkg, decl, node)
								calls = append(calls, callExpr)
								return false
							}
						}
					}
				}
			}
			return true
		})
	}
	return
}

var lookupOK = true

func lookupMethodStruct(pkg *packages.Package, decl *ast.FuncDecl, callNode *ast.CallExpr) (*types.Func, *types.Named) {
	method := pkg.TypesInfo.Defs[decl.Name]
	fn := method.(*types.Func)
	sign := fn.Type().(*types.Signature)
	if sign.Recv() == nil {
		fmt.Println("bus.Dispatch in", method)
		must(ast.Print(cfg.Fset, callNode))
		lookupOK = false
		return nil, nil
	}
	st := ptrNamedStruct(sign.Recv().Type())
	return fn, st
}

type MessageType struct {
	Message *types.Named
	Struct  *types.Named
	Method  *types.Func
}

func unwrapIdentType(pkg *packages.Package, id *ast.Ident) *types.Named {
	typ := pkg.TypesInfo.TypeOf(id)
	if typ == nil {
		panicf("%v: id %v unknown", pkg.PkgPath, id)
	}
	msg := ptrNamedStruct(typ)
	if msg == nil {
		panicf("%v: id %v must be a pointer to struct (%v)", pkg.PkgPath, id, typ)
	}
	return msg
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func panicf(msg string, args ...interface{}) {
	txt := fmt.Sprintf(msg, args...)
	panic(txt)
}

var loadFiles = map[string][]byte{}

func loadFile(file string) []byte {
	data := loadFiles[file]
	if data != nil {
		return data
	}
	data, err := ioutil.ReadFile(file)
	must(err)
	loadFiles[file] = data
	return data
}

type Replacement struct {
	Start, End token.Pos
	ReplacedBy string
	Priority   int // smaller is more priority
}

func replaceNode(node ast.Node, by string) Replacement {
	return Replacement{
		Start:      node.Pos(),
		End:        node.End(),
		ReplacedBy: by,
	}
}

func replaceFiles(fset *token.FileSet, rs []Replacement) {
	sort.Slice(rs, func(i, j int) bool {
		if rs[i].Start == rs[j].Start { // when equal, order by priority
			return rs[i].Priority < rs[j].Priority
		}
		return rs[i].Start < rs[j].Start
	})

	var rData, wData []byte
	var currFilename string
	var lastOffset int
	for _, r := range rs {
		start := fset.PositionFor(r.Start, false)
		end := fset.PositionFor(r.End, false)
		if start.Filename != end.Filename {
			panic("inconsistent")
		}
		filename := start.Filename
		if filename != currFilename {
			// append the end block then write
			if currFilename != "" && lastOffset != 0 {
				wData = append(wData, rData[lastOffset:]...)
				writeNewFile(currFilename, wData)
			}

			// reset
			lastOffset = 0
			wData = wData[:0]

			// load the next file
			rData = loadFile(filename)
			currFilename = filename
		}

		// do the replacement
		wData = append(wData, rData[lastOffset:start.Offset]...)
		wData = append(wData, r.ReplacedBy...)
		lastOffset = end.Offset
	}

	// append the end block then write
	if currFilename != "" && lastOffset != 0 {
		wData = append(wData, rData[lastOffset:]...)
		writeNewFile(currFilename, wData)
	}
}

func writeNewFile(filename string, data []byte) {
	if !*flReplace {
		dir, base := filepath.Dir(filename), filepath.Base(filename)
		filename = filepath.Join(dir, "__"+base)
	}

	if *flWrite {
		fmt.Println("write file", filename)
		must(ioutil.WriteFile(filename, data, 0644))
	} else {
		fmt.Println("will write file", filename)
	}
}

func addInterfaces(msgs map[MessageType]bool) (result []Replacement) {
	m := map[*types.Named][]*types.Func{}
	for msg := range msgs {
		found := false
		for _, method := range m[msg.Struct] {
			if method == msg.Method {
				found = true
				break
			}
		}
		if !found {
			m[msg.Struct] = append(m[msg.Struct], msg.Method)
		}
	}

	for st, methods := range m {
		pos := st.Obj().Pos() - 5
		result = append(result, Replacement{
			Start:      pos,
			End:        pos,
			ReplacedBy: fmt.Sprintf("type %vInterface interface {\n", st.Obj().Name()),
			Priority:   10,
		})
		result = append(result, Replacement{
			Start:      pos,
			End:        pos,
			ReplacedBy: "}\n\n",
			Priority:   10 + len(methods) + 1,
		})

		sort.Slice(methods, func(i, j int) bool {
			return methods[i].Name() < methods[j].Name()
		})
		for i, method := range methods {
			methodPos := cfg.Fset.PositionFor(method.Pos(), false)
			srcData := loadFile(methodPos.Filename)[methodPos.Offset:]
			idx := bytes.IndexByte(srcData, ')')
			methodTxt := srcData[:idx+1]

			result = append(result, Replacement{
				Start:      pos,
				End:        pos,
				ReplacedBy: fmt.Sprintf("\n\t%s error\n", methodTxt),
				Priority:   10 + i + 1,
			})
		}
	}
	return
}
