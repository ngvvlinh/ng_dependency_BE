package grpcgen

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"etop.vn/backend/tools/pkg/gen"
)

// Method ...
type Method struct {
	Service    *Service
	Name       string
	InputType  string
	OutputType string
}

// FullPath ...
func (m *Method) FullPath() string {
	s := m.Service.Package + "." + m.Service.Name + "/" + m.Name
	if m.Service.PkgPrefix != "" {
		s = m.Service.PkgPrefix + "/" + s
	}
	return s
}

// Result ...
type Result struct {
	Package     string
	PackagePath string
	Services    []*Service
	Imports     []Import
	Interfaces  []string
}

// ExtraImports returns imports which specific to the project.
func (r Result) ExtraImports() []Import {
	var ims []Import
	for _, im := range r.Imports {
		if !strings.HasPrefix(im.Path, gen.ProjectImport_) {
			continue
		}
		if r.ImportInUse(im) {
			ims = append(ims, im)
		}
	}
	return ims
}

// ImportInUse reports whether given import is used in the file.
func (r Result) ImportInUse(im Import) bool {
	imPrefix := "*" + im.Name + "."
	if im.Name == "" {
		imPrefix = "*" + filepath.Base(im.Path) + "."
	}
	for _, s := range r.Services {
		for _, m := range s.Methods {
			if strings.HasPrefix(m.InputType, imPrefix) || strings.HasPrefix(m.OutputType, imPrefix) {
				return true
			}
		}
	}
	return false
}

// Service ...
type Service struct {
	PkgPrefix  string
	Package    string
	Name       string
	MapMethods map[string]*Method
	Methods    []*Method
}

// Import ...
type Import struct {
	Full string
	Name string
	Path string
}

// Options ...
type Options struct {
	IncludeInterface     func(name string) bool
	ImportCurrentPackage bool
}

// ParseServiceFile ...
func ParseServiceFile(inputPath string, opt Options) Result {
	p := newParser()
	return p.parse(inputPath, opt)
}

type sortMethodsType []*Method

func (a sortMethodsType) Len() int           { return len(a) }
func (a sortMethodsType) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a sortMethodsType) Less(i, j int) bool { return a[i].Name < a[j].Name }

// SortMethods ...
func SortMethods(methods map[string]*Method) []*Method {
	sortMethods := make(sortMethodsType, 0, len(methods))
	for _, m := range methods {
		sortMethods = append(sortMethods, m)
	}
	sort.Sort(sortMethods)
	return sortMethods
}

type parserStruct struct {
	options       Options
	inputFilePath string

	fset *token.FileSet
}

func newParser() *parserStruct {
	return &parserStruct{
		fset: token.NewFileSet(),
	}
}

// inputPath is relative to projectPath (we don't allow inputPath outside the project)
func (p *parserStruct) parse(inputPath string, opt Options) Result {
	absPath := gen.GetAbsPath(inputPath)
	f, err := parser.ParseFile(p.fset, absPath, nil, 0)
	if err != nil {
		p.Fatalf("Unable to parse file `%v`.\n  Error: %v\n", absPath, err)
		return Result{}
	}

	data, _ := ioutil.ReadFile(absPath)
	p.options = opt
	p.inputFilePath = inputPath

	pkg := getPackageName(data)
	pkgPath, _ := filepath.Rel(gen.ProjectPath(), filepath.Dir(absPath))
	pkgPath = gen.ProjectImport + "/" + pkgPath

	services := p.extract(f)
	for _, s := range services {
		s.Package = pkg
		for _, m := range s.Methods {
			if opt.ImportCurrentPackage {
				m.InputType, _ = prependPackage(pkg, m.InputType)
				m.OutputType, _ = prependPackage(pkg, m.OutputType)
			}
		}
	}
	imports := extractImports(f.Imports)
	if opt.ImportCurrentPackage {
		imports = append(imports, Import{
			Full: pkg + ` "` + pkgPath + `"`,
			Name: pkg,
			Path: pkgPath,
		})
	}
	return Result{
		Interfaces:  getInterfaceNames(data),
		Imports:     imports,
		Package:     pkg,
		PackagePath: pkgPath,
		Services:    services,
	}
}

func (p *parserStruct) extract(f *ast.File) (services []*Service) {
	var s *Service
	inspectFunc := func(node ast.Node) bool {
		switch node := node.(type) {
		case *ast.TypeSpec:
			switch typ := node.Type.(type) {
			case *ast.InterfaceType:
				name := node.Name.Name
				if !p.options.IncludeInterface(name) {
					return false
				}

				methods := p.extractMethods(typ)
				s = &Service{
					Name:       name,
					MapMethods: methods,
					Methods:    SortMethods(methods),
				}
				for _, m := range methods {
					m.Service = s
				}
				services = append(services, s)
				return false
			}
		}
		return true
	}

	ast.Inspect(f, inspectFunc)
	return services
}

func (p *parserStruct) extractMethods(typ *ast.InterfaceType) map[string]*Method {
	methods := make(map[string]*Method)
	for _, m := range typ.Methods.List {
		name := m.Names[0].Name

		fnTyp := m.Type.(*ast.FuncType)
		params := fnTyp.Params.List
		if len(params) != 2 {
			p.Fatalf("Error: Method %v must have exactly 2 arguments\n", name)
		}
		inputName := p.extractTypeName(name, params[1].Type)

		results := fnTyp.Results.List
		if len(results) != 2 {
			p.Fatalf("Error: Method %v must have exactly 2 results\n", name)
		}
		outputName := p.extractTypeName(name, results[0].Type)

		methods[name] = &Method{
			Name:       name,
			InputType:  inputName,
			OutputType: outputName,
		}
	}
	return methods
}

func (p *parserStruct) extractTypeName(method string, typ ast.Expr) string {
	s := ""
	if t, ok := typ.(*ast.StarExpr); ok {
		typ = t.X
		s = "*"
	}

	switch typ := typ.(type) {
	case *ast.Ident:
		return s + typ.Name

	case *ast.SelectorExpr:
		x := (typ.X).(*ast.Ident)
		return s + x.Name + "." + typ.Sel.Name

	default:
		err := ast.Print(p.fset, typ)
		if err != nil {
			panic(err)
		}
		p.Fatalf("Unable to parse type")
	}
	return "[ERROR]"
}

func extractImports(ims []*ast.ImportSpec) []Import {
	r := make([]Import, len(ims))
	for i, im := range ims {

		// strip surrounding quotation marks (")
		path := im.Path.Value[1 : len(im.Path.Value)-1]

		name := ""
		full := im.Path.Value
		if im.Name != nil {
			name = im.Name.Name
			full = name + " " + im.Path.Value
		}
		r[i] = Import{
			Full: full,
			Name: name,
			Path: path,
		}
	}
	return r
}

func getPackageName(data []byte) string {
	m := regexp.MustCompile(`package ([A-z0-9_]+)`).FindSubmatch(data)
	if m == nil {
		fmt.Println("Invalid package name")
		os.Exit(1)
	}
	return string(m[1])
}

func getInterfaceNames(data []byte) []string {
	ms := regexp.MustCompile(`type ([A-z0-9_]+) interface`).FindAllSubmatch(data, -1)
	if ms == nil {
		return nil
	}

	interfaceNames := make([]string, 0, len(ms))
	for _, m := range ms {
		interfaceNames = append(interfaceNames, string(m[1]))
	}
	return interfaceNames
}

func (p *parserStruct) Fatalf(format string, args ...interface{}) {
	fmt.Printf("Error parsing file: %v\n", p.inputFilePath)
	if format[len(format)-1] != '\n' {
		format += "\n"
	}
	fmt.Printf(format, args)
	os.Exit(1)
}

func prependPackage(pkg, name string) (string, bool) {
	if strings.Contains(name, ".") {
		return name, false
	}
	return "*" + pkg + "." + name[1:], true
}
