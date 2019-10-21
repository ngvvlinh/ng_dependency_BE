package apix

import (
	"fmt"
	"go/types"
	"strings"

	"golang.org/x/tools/go/packages"

	"etop.vn/backend/tools/pkg/generator"
	"etop.vn/common/l"
)

var ll = l.New()
var _ generator.Plugin = &plugin{}

type plugin struct {
	generator.Filter
	ng generator.Engine
}

func New() generator.Plugin {
	return &plugin{
		Filter: generator.FilterByCommand("gen:apix"),
	}
}

func (p *plugin) Name() string { return "apix" }

func (p *plugin) Generate(ng generator.Engine) error {
	p.ng = ng
	return ng.GenerateEachPackage(p.generatePackage)
}

func (p *plugin) generatePackage(ng generator.Engine, pkg *packages.Package, printer generator.Printer) error {
	ll.V(2).Debugf("apix: generating package %v", pkg.PkgPath)
	objects := ng.GetObjectsByPackage(pkg)
	var services []*Service
	for _, obj := range objects {
		ll.V(2).Debugf("  object %v: %v", obj.Name(), obj.Type())
		directives := ng.GetDirectives(obj)
		switch obj := obj.(type) {
		case *types.TypeName:
			ll.V(2).Debugf("  type %v", obj.Name())
			switch typ := obj.Type().(type) {

			case *types.Named:
				switch typ := typ.Underlying().(type) {
				case *types.Interface:
					if !strings.HasSuffix(obj.Name(), "API") {
						ll.V(1).Debugf("ignore unrecognized interface %v", obj.Name())
						continue
					}
					service, err := p.parseService(typ)
					if err != nil {
						return generator.Errorf(err, "service %v: %v", obj.Name(), err)
					}
					service.Name = strings.TrimSuffix(obj.Name(), "API")
					service.APIPath = parseAPIPath(directives)
					services = append(services, service)
				}
			}
		}
	}
	for _, s := range services {
		fmt.Println("service", s.Name, s.APIPath)
		for _, m := range s.Methods {
			fmt.Println("  method  ", m.Name, strings.ReplaceAll(m.Comment, "\n", `\n`))
			fmt.Println("  request ", m.Request.Type.Underlying())
			fmt.Println("  response", m.Response.Name)
		}
	}
	return p.generateServices(printer, services)
}

func parseAPIPath(directives []generator.Directive) string {
	for _, d := range directives {
		if d.Cmd == "apix:path" {
			return d.Arg
		}
	}
	return ""
}

func (p *plugin) parseService(iface *types.Interface) (*Service, error) {
	ng := p.ng
	methods := make([]*Method, 0, iface.NumMethods())
	for i, n := 0, iface.NumMethods(); i < n; i++ {
		method := iface.Method(i)
		req, resp, err := p.parseMethod(method)
		if err != nil {
			return nil, generator.Errorf(err, "method %v: %v", method.Name(), err)
		}
		m := &Method{
			Name:     method.Name(),
			Comment:  ng.GetComment(method).Text(),
			Request:  req,
			Response: resp,
		}
		methods = append(methods, m)
	}
	return &Service{Methods: methods}, nil
}

func (p *plugin) parseMethod(method *types.Func) (req, resp Message, err error) {
	sign := method.Type().(*types.Signature)
	{
		params := sign.Params()
		if params.Len() != 2 {
			err = generator.Errorf(nil, "must have 2 params")
			return
		}
		if err = validateType(params.At(0).Type(), "context", "Context"); err != nil {
			return
		}
		req, err = parseMessage(params.At(1))
		if err != nil {
			err = generator.Errorf(err, "param %v", err)
			return
		}
	}
	{
		results := sign.Results()
		if results.Len() != 2 {
			err = generator.Errorf(nil, "must have 2 results")
			return
		}
		if err = validateType(results.At(1).Type(), "", "error"); err != nil {
			err = generator.Errorf(err, "result %v", err)
			return
		}
		resp, err = parseMessage(results.At(0))
		if err != nil {
			return
		}
	}
	return
}

func parseMessage(m *types.Var) (_ Message, err error) {
	typ := m.Type()
	ptr, ok := typ.(*types.Pointer)
	if !ok {
		err = generator.Errorf(nil, "must be pointer")
		return
	}
	typ = ptr.Elem()
	named, ok := typ.(*types.Named)
	if !ok {
		err = generator.Errorf(nil, "must be pointer to named type")
		return
	}
	return Message{
		Type:    named,
		PkgPath: m.Pkg().Path(),
		Name:    named.Obj().Name(),
	}, nil

	// TODO: parse struct
}

func validateType(typ types.Type, pkgPath, name string) error {
	named, ok := typ.(*types.Named)
	if !ok {
		return generator.Errorf(nil, "must be a named type")
	}
	var typPkgPath string
	if named.Obj().Pkg() != nil {
		typPkgPath = named.Obj().Pkg().Path()
	}
	if typPkgPath == pkgPath && named.Obj().Name() == name {
		return nil
	}
	return generator.Errorf(nil, "must be %v.%v (got %v)", pkgPath, name, typ)
}
