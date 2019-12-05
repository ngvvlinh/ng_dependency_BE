package wrapper

import (
	"go/types"
	"reflect"
	"sort"
	"strings"

	"golang.org/x/tools/go/packages"

	"etop.vn/backend/tools/pkg/generator"
	"etop.vn/backend/tools/pkg/genutil"
	"etop.vn/common/l"
)

var ll = l.New()
var _ generator.Plugin = &plugin{}

const SuffixService = "Service"

type plugin struct {
	generator.Qualifier
}

func New() generator.Plugin {
	return &plugin{Qualifier: genutil.Qualifier{}}
}

func (p *plugin) Name() string { return "wrapper" }

func (p *plugin) Filter(ng generator.FilterEngine) error {
	for _, pkg := range ng.ParsingPackages() {
		ds := pkg.Directives
		d, ok := ds.Get("gen:wrapper")
		if ok {
			pkg.Include()
			ng.ParsePackage(d.Arg)
		}
	}
	return nil
}

func (p *plugin) Generate(ng generator.Engine) error {
	for _, pkg := range ng.GeneratingPackages() {
		printer := pkg.GetPrinter()
		if err := p.generatePackage(ng, pkg.Package, printer); err != nil {
			return err
		}
		if err := printer.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (p *plugin) generatePackage(ng generator.Engine, pkg *packages.Package, printer generator.Printer) error {
	directives := ng.GetDirectivesByPackage(pkg)
	directive, _ := directives.Get("gen:wrapper")
	if directive.Arg == "" {
		return generator.Errorf(nil, "no arg")
	}
	pkgPrefix, _ := directives.Get("gen:wrapper:prefix")
	pkgName, _ := directives.Get("gen:wrapper:package")

	pbPkg := ng.GetPackageByPath(directive.Arg)
	if pbPkg == nil {
		return generator.Errorf(nil, "package %v not found", directive.Arg)
	}
	ifaces := filterServiceInterfaces(ng, ng.GetObjectsByPackage(pbPkg))
	structs := filterServiceStructs(ng.GetObjectsByPackage(pkg))
	structNames := getKeys(structs)
	ifaceNames := getKeys(ifaces)

	{ // compare service names
		ifaceStr := strings.Join(ifaceNames, ",")
		structStr := strings.Join(structNames, ",")
		if structStr != ifaceStr {
			for i, j := 0, 0; i < len(ifaceNames) || j < len(structNames); {
				var si, sj string
				if i < len(ifaceNames) {
					si = ifaceNames[i]
				}
				if j < len(structNames) {
					sj = structNames[j]
				}
				switch {
				case si == sj:
					i++
					j++
				case sj == "":
					i++
					ll.S.Errorf("%v found in %v but not found in %v", si, pbPkg.PkgPath, pkg.PkgPath)
				case si == "":
					j++
					ll.S.Errorf("%v not found in %v but found in %v", sj, pbPkg.PkgPath, pkg.PkgPath)
				case sj > si:
					i++
					ll.S.Errorf("%v found in %v but not found in %v", si, pbPkg.PkgPath, pkg.PkgPath)
				case si > sj:
					j++
					ll.S.Errorf("%v not found in %v but found in %v", sj, pbPkg.PkgPath, pkg.PkgPath)
				}
			}
			return generator.Errorf(nil, "mismatch services between %v and %v", pbPkg.PkgPath, pkg.PkgPath)
		}
	}

	ok := true
	var services []*Service
	for _, name := range ifaceNames {
		iface := ifaces[name]
		named := structs[name]
		ifaceMethods := getMethods(iface)
		structMethods := getMethods(named)

		service := &Service{
			PkgPb:          pbPkg.PkgPath,
			PkgPrefix:      pkgPrefix.Arg,
			PkgName:        pkgName.Arg,
			PkgPath:        pkg.PkgPath,
			EndpointPrefix: iface.EndpointPrefix.Arg,
			Name:           strings.TrimSuffix(name, SuffixService),
		}
		var methods []*Method
		for methodName := range ifaceMethods {
			if structMethods[methodName] == nil {
				ll.S.Errorf("method %v in interface %v not found in struct %v", methodName, name, name)
				ok = false
				continue
			}
			method, err := parseMethod(methodName, ifaceMethods[methodName], structMethods[methodName])
			if err != nil {
				ll.S.Errorf("method %v in struct %v: %v", methodName, name, err)
				ok = false
				continue
			}
			method.Service = service
			methods = append(methods, method)
		}
		sort.Slice(methods, func(i, j int) bool {
			return methods[i].Name < methods[j].Name
		})
		service.Methods = methods
		services = append(services, service)
	}
	sort.Slice(services, func(i, j int) bool {
		return services[i].Name < services[j].Name
	})
	if !ok {
		return generator.Errorf(nil, "validation failed!")
	}

	return generate(printer, pkg, services)
}

type Service struct {
	PkgPb          string
	PkgPrefix      string // hack for etop/apix (external API)
	PkgName        string
	PkgPath        string
	EndpointPrefix string
	Name           string
	Methods        []*Method
}

type Method struct {
	Service *Service
	Name    string
	Kind    int // 1: old, 2: new
	Req     types.Type
	Resp    types.Type
}

func (m *Method) FullPath() string {
	s := m.Service.PkgName + "." + m.Service.Name + "/" + m.Name
	if m.Service.PkgPrefix != "" {
		s = m.Service.PkgPrefix + "/" + s
	}
	return s
}

func parseMethod(name string, ifaceMethod, method *types.Func) (*Method, error) {
	isign := ifaceMethod.Type().(*types.Signature)
	m := &Method{
		Name: name,
		Req:  isign.Params().At(1).Type(),
		Resp: isign.Results().At(0).Type(),
	}
	sign := method.Type().(*types.Signature)
	params := sign.Params()
	results := sign.Results()
	switch {
	case params.Len() == 2 && results.Len() == 1:
		m.Kind = 1

	case params.Len() == 3 && results.Len() == 2:
		m.Kind = 2

	default:
		return nil, generator.Errorf(nil, "invalid signature (%v)", sign)
	}
	return m, nil
}

type Interface struct {
	*types.Interface
	EndpointPrefix generator.Directive
}

func filterServiceInterfaces(ng generator.Engine, objs []types.Object) map[string]*Interface {
	result := map[string]*Interface{}
	for _, obj := range objs {
		endpointPrefix, _ := ng.GetDirectives(obj).Get("wrapper:endpoint-prefix")
		if named, ok := obj.Type().(*types.Named); ok {
			if typ, ok := named.Underlying().(*types.Interface); ok {
				name := obj.Name()
				if strings.HasSuffix(name, SuffixService) && name != "QueryService" {
					result[name] = &Interface{
						Interface:      typ,
						EndpointPrefix: endpointPrefix,
					}
				}
			}
		}
	}
	return result
}

func filterServiceStructs(objs []types.Object) map[string]*types.Named {
	result := map[string]*types.Named{}
	for _, obj := range objs {
		if named, ok := obj.Type().(*types.Named); ok {
			if _, ok := named.Underlying().(*types.Struct); ok {
				name := obj.Name()
				if strings.HasSuffix(name, SuffixService) {
					result[name] = named
				}
			}
		}
	}
	return result
}

func getKeys(m interface{}) []string {
	vkeys := reflect.ValueOf(m).MapKeys()
	keys := make([]string, len(vkeys))
	for i := range vkeys {
		keys[i] = vkeys[i].String()
	}
	sort.Strings(keys)
	return keys
}

type methoder interface {
	NumMethods() int
	Method(int) *types.Func
}

func getMethods(typ methoder) map[string]*types.Func {
	ll.V(3).Debugf("getMethods of type %v (n=%v)", typ, typ.NumMethods())
	result := map[string]*types.Func{}
	for i, n := 0, typ.NumMethods(); i < n; i++ {
		method := typ.Method(i)
		ll.V(4).Debugf("method %v", method)
		if method.Exported() {
			result[method.Name()] = method
		}
	}
	return result
}
