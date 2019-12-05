package apix

import (
	"golang.org/x/tools/go/packages"

	"etop.vn/backend/tools/pkg/generator"
	"etop.vn/backend/tools/pkg/generators/api/defs"
	"etop.vn/backend/tools/pkg/generators/api/parse"
	"etop.vn/backend/tools/pkg/genutil"
	"etop.vn/common/l"
)

var ll = l.New()
var _ generator.Plugin = &plugin{}

type plugin struct {
	generator.Filterer
	generator.Qualifier
	ng generator.Engine
}

func New() generator.Plugin {
	return &plugin{
		Filterer:  generator.FilterByCommand("gen:apix"),
		Qualifier: genutil.Qualifier{},
	}
}

func (p *plugin) Name() string { return "apix" }

func (p *plugin) Generate(ng generator.Engine) error {
	p.ng = ng
	return ng.GenerateEachPackage(p.generatePackage)
}

func (p *plugin) generatePackage(ng generator.Engine, pkg *packages.Package, printer generator.Printer) (_err error) {
	ll.V(2).Debugf("apix: generating package %v", pkg.PkgPath)

	pkgDirectives := ng.GetDirectivesByPackage(pkg)
	basePath := pkgDirectives.GetArg("gen:apix:base-path")
	if basePath == "" {
		basePath = "/api"
	}
	opts := Opts{
		BasePath: basePath,
	}

	services, err := parse.Services(ng, pkg, []defs.Kind{defs.KindService})
	if err != nil {
		return err
	}
	for _, service := range services {
		if service.APIPath == "" {
			return generator.Errorf(nil, "no api path for %v", service.Name)
		}
	}
	if err := generateServices(printer, opts, services); err != nil {
		return err
	}
	return nil
}
