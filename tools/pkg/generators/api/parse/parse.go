package parse

import (
	"fmt"
	"go/types"
	"strings"

	"golang.org/x/tools/go/packages"

	"etop.vn/backend/tools/pkg/generator"
	"etop.vn/backend/tools/pkg/generators/api/defs"
	"etop.vn/common/l"
)

var ll = l.New()

func Services(ng generator.Engine, pkg *packages.Package, kinds []defs.Kind) (services []*defs.Service, _ error) {
	objects := ng.GetObjectsByPackage(pkg)
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
					kind := parseKind(kinds, obj.Name())
					if kind == "" {
						ll.V(1).Debugf("ignore unrecognized interface %v", obj.Name())
						continue
					}
					methods, err := parseService(ng, underlyingType)
					if err != nil {
						return nil, generator.Errorf(err, "service %v: %v", obj.Name(), err)
					}

					service := &defs.Service{
						Kind:     kind,
						Name:     strings.TrimSuffix(obj.Name(), string(kind)),
						FullName: obj.Name(),
						APIPath:  directives.GetArg("apix:path"),
						Methods:  methods,
					}
					services = append(services, service)
					for _, m := range methods {
						m.Service = service
					}
				}
			}
		}
	}
	return services, nil
}

func parseKind(kinds []defs.Kind, name string) defs.Kind {
	for _, kind := range kinds {
		suffix := string(kind)
		if strings.HasSuffix(name, suffix) {
			return kind
		}
	}
	return ""
}

func parseService(ng generator.Engine, iface *types.Interface) ([]*defs.Method, error) {
	methods := make([]*defs.Method, 0, iface.NumMethods())
	for i, n := 0, iface.NumMethods(); i < n; i++ {
		method := iface.Method(i)
		if !method.Exported() {
			continue
		}
		m, err := parseMethod(ng, method)
		if err != nil {
			return nil, generator.Errorf(err, "method %v: %v", method.Name(), err)
		}
		methods = append(methods, m)
	}
	return methods, nil
}

func parseMethod(ng generator.Engine, method *types.Func) (_ *defs.Method, err error) {
	mtyp := method.Type()
	styp := mtyp.(*types.Signature)
	params := styp.Params()
	results := styp.Results()
	requests, responses, err := checkMethodSignature(method.Name(), params, results)
	if err != nil {
		return nil, fmt.Errorf("%v: %v", method.Name(), err)
	}
	return &defs.Method{
		Name:     method.Name(),
		Comment:  ng.GetComment(method).Text(),
		Method:   method,
		Request:  requests,
		Response: responses,
	}, nil
}

func checkMethodSignature(name string, params *types.Tuple, results *types.Tuple) (request, response *defs.Message, err error) {
	if params.Len() == 0 {
		err = generator.Errorf(nil, "expect at least 1 param")
		return
	}
	if results.Len() == 0 {
		err = generator.Errorf(nil, "expect at least 1 param")
		return
	}
	var requestItems, responseItems []*defs.ArgItem
	{
		t := params.At(0)
		if t.Type().String() != "context.Context" {
			err = generator.Errorf(nil, "expect the first param is context.Context")
			return
		}
	}
	{
		t := results.At(results.Len() - 1)
		if t.Type().String() != "error" {
			err = generator.Errorf(nil, "expect the last return value is error")
			return
		}
	}
	{
		// skip the first param (context.Context)
		for i, n := 1, params.Len(); i < n; i++ {
			arg, err := checkArg(params.At(i), n == 2)
			if err != nil {
				return nil, nil, generator.Errorf(err, "%v: %v", name, err)
			}
			requestItems = append(requestItems, arg)
			if !arg.Inline && arg.Name == "" {
				return nil, nil, generator.Errorf(err, "%v: must provide name for param %v", name, arg.Type)
			}
		}
	}
	{
		// skip the last result (error)
		for i, n := 0, results.Len()-1; i < n; i++ {
			arg, err := checkArg(results.At(i), n == 2)
			if err != nil {
				return nil, nil, generator.Errorf(err, "%v: %v", name, err)
			}
			responseItems = append(responseItems, arg)
		}
		if len(responseItems) > 1 {
			for _, arg := range responseItems {
				if arg.Name == "" || strings.HasPrefix(arg.Name, "_") {
					return nil, nil, generator.Errorf(err, "%v: must provide name for result %v", name, arg.Type)
				}
			}
		}
	}
	request = &defs.Message{Items: requestItems}
	response = &defs.Message{Items: responseItems}
	return request, response, nil
}

func checkArg(v *types.Var, autoInline bool) (*defs.ArgItem, error) {
	arg := &defs.ArgItem{
		Inline: v.Name() == "_" || v.Name() == "" && autoInline,
		Name:   toTitle(v.Name()),
		Var:    v,
		Type:   v.Type(),
	}
	// when inline, the param must be struct or pointer to struct
	if arg.Inline {
		var err error
		arg.Struct, arg.Ptr, err = checkStruct(v.Type())
		if err != nil {
			return nil, fmt.Errorf("type must be a struct or a pointer to struct to be inline: %v", err)
		}
	}
	return arg, nil
}

func checkStruct(t types.Type) (_ *types.Struct, ptr bool, _ error) {
	p, ptr := t.(*types.Pointer)
	if ptr {
		t = p.Elem()
	}

underlying:
	switch typ := t.(type) {
	case *types.Pointer:
		return nil, false, fmt.Errorf("got double pointer (%v)", t)

	case *types.Named:
		t = typ.Underlying()
		goto underlying

	case *types.Struct:
		return typ, ptr, nil

	default:
		return nil, false, fmt.Errorf("got %v", typ)
	}
}

func toTitle(s string) string {
	s = strings.TrimPrefix(s, "_")
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[0:1]) + s[1:]
}
