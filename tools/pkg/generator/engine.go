package generator

import (
	"go/ast"
	"go/token"
	"go/types"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"golang.org/x/tools/go/packages"
)

type Positioner interface {
	Pos() token.Pos
}

type PreparsedPackage struct {
	PkgPath    string
	Imports    map[string]*packages.Package
	Directives []Directive
}

type GeneratingPackage struct {
	*packages.Package

	directives []Directive
	plugin     *pluginStruct
	engine     *engine
	printer    *printer
}

func (g *GeneratingPackage) GetPrinter() Printer {
	if g.printer == nil {
		input := GenerateFileNameInput{PluginName: g.plugin.name}
		filename := g.engine.genFilename(input)
		dir := filepath.Dir(g.Package.CompiledGoFiles[0])
		filePath := filepath.Join(dir, filename)
		g.printer = newPrinter(g.engine, g.plugin, g.Package, filePath)
	}
	return g.printer
}

func (g *GeneratingPackage) GetDirectives() []Directive {
	return cloneDirectives(g.directives)
}

func (g *GeneratingPackage) GetObjects() []types.Object {
	return g.engine.GetObjectsByPackage(g.Package)
}

type Engine interface {
	GenerateEachPackage(func(Engine, *packages.Package, Printer) error) error
	GeneratingPackages() []*GeneratingPackage

	GetComment(Positioner) Comment
	GetDirectives(Positioner) []Directive
	GetDirectivesByPackage(*packages.Package) []Directive
	GetIdent(Positioner) *ast.Ident
	GetObject(Positioner) types.Object
	GetObjectsByPackage(*packages.Package) []types.Object
	GetObjectsByScope(*types.Scope) []types.Object
	GetPackage(Positioner) *packages.Package
	GetPackageByPath(string) *packages.Package
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
	mapPkgDirectives  map[string][]Directive
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

func (ng *engine) GetComment(p Positioner) Comment {
	cmt := ng.xinfo.GetComment(ng.GetIdent(p))
	return cmt
}

func (ng *engine) CommentByIdent(ident *ast.Ident) Comment {
	cmt := ng.xinfo.GetComment(ident)
	return cmt
}

func (ng *engine) CommentByObject(obj types.Object) Comment {
	ident := ng.GetIdentByPos(obj.Pos())
	return ng.CommentByIdent(ident)
}

func (ng *engine) GetDirectives(p Positioner) []Directive {
	return ng.GetComment(p).Directives
}

func (ng *engine) GetIdent(p Positioner) *ast.Ident {
	return ng.GetIdentByPos(p.Pos())
}

func (ng *engine) GetIdentByObject(obj types.Object) *ast.Ident {
	return ng.GetIdentByPos(obj.Pos())
}

func (ng *engine) GetIdentByPos(pos token.Pos) *ast.Ident {
	return ng.xinfo.Positions[pos]
}

func (ng *engine) GetObject(p Positioner) types.Object {
	return ng.GetObjectByIdent(ng.GetIdent(p))
}

func (ng *engine) GetObjectByIdent(ident *ast.Ident) types.Object {
	return ng.xinfo.GetDef(ident)
}

func (ng *engine) GetPackage(p Positioner) *packages.Package {
	return ng.GetPackageByIdent(ng.GetIdent(p))
}

func (ng *engine) GetPackageByIdent(ident *ast.Ident) *packages.Package {
	decl := ng.xinfo.Declarations[ident]
	if decl == nil {
		return nil
	}
	return decl.Pkg
}

func (ng *engine) GetPackageByPath(pkgPath string) *packages.Package {
	return ng.pkgMap[pkgPath]
}

func (ng *engine) GetObjectsByPackage(pkg *packages.Package) []types.Object {
	return ng.GetObjectsByScope(pkg.Types.Scope())
}

func (ng *engine) GetObjectsByScope(s *types.Scope) []types.Object {
	names := s.Names()
	objs := make([]types.Object, len(names))
	for i, name := range names {
		objs[i] = s.Lookup(name)
	}
	return objs
}

func (ng *wrapEngine) GenerateEachPackage(
	fn func(Engine, *packages.Package, Printer) error,
) error {
	for _, pkg := range ng.generatingPackages() {
		printer := pkg.GetPrinter()
		if err := fn(ng, pkg.Package, printer); err != nil {
			return Errorf(err, "generating package %v: %v", pkg.PkgPath, err)
		}
		if err := printer.Close(); err != nil {
			return err
		}
	}
	return nil
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
			gpkg := ng.generatingPackage(pkg)
			pkgs = append(pkgs, gpkg)
		}
	}
	return pkgs
}

func (ng *wrapEngine) generatingPackage(pkg *packages.Package) *GeneratingPackage {
	directives := ng.GetDirectivesByPackage(pkg)
	gpkg := &GeneratingPackage{
		Package:    pkg,
		directives: directives,
		plugin:     ng.plugin,
		engine:     ng.engine,
	}
	return gpkg
}

func (ng *wrapEngine) GetDirectivesByPackage(pkg *packages.Package) []Directive {
	directives, ok := ng.mapPkgDirectives[pkg.PkgPath]
	if !ok {
		var ds []Directive
		for _, file := range pkg.GoFiles {
			body, err := ioutil.ReadFile(file)
			if err != nil {
				if os.IsNotExist(err) {
					ll.V(1).Debugf("ignore not found file: %v", file)
					continue
				}
				panic(err)
			}

			var errs []error
			ds, errs = parseDirectivesFromBody(ds, body)
			for _, err := range errs {
				ll.V(1).Debugf("invalid directive from file %v: %v", file, err)
			}
		}
		directives = ds
		ng.mapPkgDirectives[pkg.PkgPath] = ds
	}
	return cloneDirectives(directives)
}

func cloneDirectives(directives []Directive) []Directive {
	result := make([]Directive, len(directives))
	for i, d := range directives {
		result[i] = d
	}
	return result
}
