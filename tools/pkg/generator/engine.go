package generator

import (
	"go/ast"
	"go/token"
	"go/types"
	"path/filepath"
	"sync"

	"golang.org/x/tools/go/packages"
)

type Object struct {
	Ident      *ast.Ident
	Object     types.Object
	Comment    *Comment
	Directives []Directive
}

type PreparsedPackage struct {
	PkgPath    string
	Imports    map[string]*packages.Package
	Directives []Directive
}

type GeneratingPackage struct {
	Package    *packages.Package
	Directives []Directive

	plugin  *pluginStruct
	engine  *engine
	printer *printer
}

func (g *GeneratingPackage) Generate() Printer {
	if g.printer == nil {
		input := GenerateFileNameInput{PluginName: g.plugin.name}
		filename := g.engine.genFilename(input)
		dir := filepath.Dir(g.Package.CompiledGoFiles[0])
		filePath := filepath.Join(dir, filename)
		g.printer = newPrinter(g.engine, g.plugin, g.Package, filePath)
	}
	return g.printer
}

func (g *GeneratingPackage) Objects() []Object {
	return g.engine.ObjectsByPackage(g.Package)
}

type Engine interface {
	GeneratingPackages() []*GeneratingPackage

	CommentByIdent(*ast.Ident) *Comment
	CommentByObject(types.Object) *Comment
	DirectivesByIdent(*ast.Ident) []Directive
	DirectivesByObject(types.Object) []Directive
	IdentByObject(types.Object) *ast.Ident
	IdentByPos(token.Pos) *ast.Ident
	ObjectByIdent(*ast.Ident) types.Object
	ObjectsByPackage(*packages.Package) []Object
	ObjectsByScope(*types.Scope) []Object
	PackageByIdent(*ast.Ident) *packages.Package
	PackageByPath(string) *packages.Package
}

var _ Engine = &wrapEngine{}
var theEngine = newEngine()

type engine struct {
	plugins        []*pluginStruct
	enabledPlugins []*pluginStruct
	pluginsMap     map[string]*pluginStruct

	xcfg    Config
	xinfo   *extendedInfo
	pkgcfg  packages.Config
	pkgMap  map[string]*packages.Package
	srcMap  map[string][]byte
	bufPool sync.Pool

	cleanedFileNames  map[string]bool
	collectedPackages []PreparsedPackage
	includes          []bool
	generatedFile     []string
}

type wrapEngine struct {
	*engine
	plugin *pluginStruct
	pkgs   []*GeneratingPackage
}

func newEngine() *engine {
	return &engine{
		pkgMap:     make(map[string]*packages.Package),
		pluginsMap: make(map[string]*pluginStruct),
	}
}

func (ng *engine) clone() *engine {
	result := newEngine()
	result.plugins = ng.plugins
	result.pluginsMap = ng.pluginsMap
	result.bufPool = ng.bufPool
	return result
}

func (ng *engine) CommentByIdent(ident *ast.Ident) *Comment {
	cmt, _ := ng.xinfo.GetComment(ident)
	return cmt
}

func (ng *engine) CommentByObject(obj types.Object) *Comment {
	ident := ng.IdentByPos(obj.Pos())
	return ng.CommentByIdent(ident)
}

func (ng *engine) DirectivesByIdent(ident *ast.Ident) []Directive {
	_, directives := ng.xinfo.GetComment(ident)
	return directives
}

func (ng *engine) DirectivesByObject(obj types.Object) []Directive {
	ident := ng.IdentByPos(obj.Pos())
	return ng.DirectivesByIdent(ident)
}

func (ng *engine) IdentByObject(obj types.Object) *ast.Ident {
	return ng.IdentByPos(obj.Pos())
}

func (ng *engine) IdentByPos(pos token.Pos) *ast.Ident {
	return ng.xinfo.Positions[pos]
}

func (ng *engine) ObjectByIdent(ident *ast.Ident) types.Object {
	return ng.xinfo.GetDef(ident)
}

func (ng *engine) PackageByIdent(ident *ast.Ident) *packages.Package {
	decl := ng.xinfo.Declarations[ident]
	if decl == nil {
		return nil
	}
	return decl.Pkg
}

func (ng *engine) PackageByPath(pkgPath string) *packages.Package {
	return ng.pkgMap[pkgPath]
}

func (ng *engine) ObjectsByPackage(pkg *packages.Package) []Object {
	return ng.ObjectsByScope(pkg.Types.Scope())
}

func (ng *engine) ObjectsByScope(s *types.Scope) []Object {
	names := s.Names()
	objs := make([]Object, len(names))
	for i, name := range names {
		obj := s.Lookup(name)
		ident := ng.IdentByPos(obj.Pos())
		cmt, directives := ng.xinfo.GetComment(ident)
		objs[i] = Object{
			Ident:      ident,
			Object:     obj,
			Comment:    cmt,
			Directives: directives,
		}
	}
	return objs
}

func (ng *wrapEngine) GeneratingPackages() []*GeneratingPackage {
	if ng.pkgs == nil {
		ng.pkgs = ng.generatingPackages()
	}
	return ng.pkgs
}

func (ng *wrapEngine) generatingPackages() []*GeneratingPackage {
	pkgs := make([]*GeneratingPackage, 0, ng.plugin.includesN)
	includes := ng.plugin.includes
	for i, ppkg := range ng.collectedPackages {
		if includes[i] {
			pkg := ng.pkgMap[ppkg.PkgPath]
			gpkg := &GeneratingPackage{
				Package:    pkg,
				Directives: ppkg.Directives,
				plugin:     ng.plugin,
				engine:     ng.engine,
			}
			pkgs = append(pkgs, gpkg)
		}
	}
	return pkgs
}
