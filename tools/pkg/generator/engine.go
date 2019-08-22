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
	Ident   *ast.Ident
	Object  types.Object
	Comment *Comment
}

type PreparsedPackage struct {
	PkgPath    string
	Imports    map[string]*packages.Package
	Directives []Directive
}

type GeneratingPackage struct {
	Package *packages.Package

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
	return g.engine.ObjectsByScope(g.Package.Types.Scope())
}

type Engine interface {
	GeneratingPackages() []*GeneratingPackage

	CommentByIdent(*ast.Ident) *Comment
	IdentByPos(token.Pos) *ast.Ident
	ObjectByIdent(*ast.Ident) types.Object
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
	return ng.xinfo.GetComment(ident)
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

func (ng *engine) ObjectsByScope(s *types.Scope) []Object {
	names := s.Names()
	objs := make([]Object, len(names))
	for i, name := range names {
		obj := s.Lookup(name)
		ident := ng.IdentByPos(obj.Pos())
		cmt := ng.CommentByIdent(ident)
		objs[i] = Object{
			Ident:   ident,
			Object:  obj,
			Comment: cmt,
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
				Package: pkg,
				plugin:  ng.plugin,
				engine:  ng.engine,
			}
			pkgs = append(pkgs, gpkg)
		}
	}
	return pkgs
}
