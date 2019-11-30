package apix

import (
	"encoding/json"
	"fmt"
	"go/types"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/packages"

	"etop.vn/backend/tools/pkg/gen"
	"etop.vn/backend/tools/pkg/generator"
	"etop.vn/backend/tools/pkg/generators/api"
	"etop.vn/backend/tools/pkg/generators/apix/defs"
	"etop.vn/backend/tools/pkg/generators/apix/swagger"
	"etop.vn/backend/tools/pkg/generators/wrapper"
	"etop.vn/common/l"
)

var ll = l.New()
var _ generator.Plugin = &plugin{}

type plugin struct {
	generator.Filterer
	ng generator.Engine
}

func New() generator.Plugin {
	return &plugin{
		Filterer: generator.FilterByCommand("gen:apix"),
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
	docPath := pkgDirectives.GetArg("gen:apix:doc-path")
	if docPath == "" {
		panic(fmt.Sprintf("no doc-path for pkg %v", pkg.Name))
	}

	objects := ng.GetObjectsByPackage(pkg)
	var services []*defs.Service
	for _, obj := range objects {
		ll.V(2).Debugf("  object %v: %v", obj.Name(), obj.Type())
		directives := ng.GetDirectives(obj)
		switch obj := obj.(type) {
		case *types.TypeName:
			ll.V(2).Debugf("  type %v", obj.Name())
			switch typ := obj.Type().(type) {
			case *types.Named:
				switch underlyingType := typ.Underlying().(type) {
				case *types.Interface:
					if !strings.HasSuffix(obj.Name(), wrapper.SuffixService) {
						ll.V(1).Debugf("ignore unrecognized interface %v", obj.Name())
						continue
					}
					service, err := p.parseService(underlyingType)
					if err != nil {
						return generator.Errorf(err, "service %v: %v", obj.Name(), err)
					}
					service.Name = strings.TrimSuffix(obj.Name(), wrapper.SuffixService)
					service.APIPath = directives.GetArg("apix:path")
					if service.APIPath == "" {
						return generator.Errorf(nil, "no api path for %v", obj.Name())
					}
					service.BasePath = basePath
					services = append(services, service)
				}
			}
		}
	}

	if err := p.generateServices(printer, services); err != nil {
		return err
	}
	swaggerDoc, err := swagger.GenerateSwagger(ng, services)
	if err != nil {
		return generator.Errorf(err, "generate swagger: %v", err)
	}
	{
		dir := filepath.Join(gen.ProjectPath(), "doc", docPath)
		filename := filepath.Join(dir, "swagger.json")
		f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		defer func() {
			err := f.Close()
			if _err == nil {
				_err = err
			}
		}()
		encoder := json.NewEncoder(f)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(swaggerDoc); err != nil {
			return generator.Errorf(nil, "generate swagger: %v", err)
		}
	}
	return nil
}

func (p *plugin) parseService(iface *types.Interface) (*defs.Service, error) {
	methods := make([]*defs.Method, 0, iface.NumMethods())
	for i, n := 0, iface.NumMethods(); i < n; i++ {
		method := iface.Method(i)
		if !method.Exported() {
			continue
		}
		m, err := p.parseMethod(method)
		if err != nil {
			return nil, generator.Errorf(err, "method %v: %v", method.Name(), err)
		}
		methods = append(methods, m)
	}
	return &defs.Service{Methods: methods}, nil
}

func (p *plugin) parseMethod(method *types.Func) (_ *defs.Method, err error) {
	def, err := api.ExtractHandlerDef(p.ng, method)
	if err != nil {
		return nil, generator.Errorf(err, "%v", err)
	}
	return def, nil
}
