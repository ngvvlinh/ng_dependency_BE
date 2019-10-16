package cq

import (
	"fmt"
	"go/types"
	"sort"
	"strings"

	"golang.org/x/tools/go/packages"

	"etop.vn/backend/tools/pkg/generator"
)

func New() generator.Plugin {
	return &gen{
		Filter: generator.FilterByCommand("gen:api"),
	}
}

type gen struct {
	generator.Filter
}

func (g *gen) Name() string { return "cq" }

func (g *gen) Generate(ng generator.Engine) error {
	return ng.GenerateEachPackage(generatePackage)
}

func generatePackage(ng generator.Engine, pkg *packages.Package, printer generator.Printer) error {
	var services []ServiceDef
	kinds := []string{QueryService, Aggregate}
	for _, obj := range ng.GetObjectsByPackage(pkg) {
		switch obj := obj.(type) {
		case *types.TypeName:
			kind, iface, err := checkService(kinds, obj, ng.GetComment(obj))
			if err != nil {
				errorf("%v\n", err)
				continue
			}
			if iface != nil {
				services = append(services, ServiceDef{
					Kind: kind,
					Name: obj.Name(),
					Type: iface,
				})
			}
		}
	}
	mustNoError("package %v:\n", pkg.PkgPath)
	if len(services) == 0 {
		fmt.Printf("  skipped %v\n", pkg.PkgPath)
		return nil
	}
	sort.Slice(services, func(i, j int) bool {
		return services[i].Name < services[j].Name
	})

	w := NewWriter(pkg.Name, pkg.PkgPath, printer, printer)
	ws := &MultiWriter{Writer: w}
	writeCommonDeclaration(ws)
	for _, item := range services {
		debugf("processing service %v", item.Name)
		processService(ws, item)
	}

	p(w, "\n// implement interfaces\n\n")
	mustWrite(w, ws.WriteIface.Bytes())
	p(w, "\n// implement conversion\n\n")
	mustWrite(w, ws.WriteArgs.Bytes())
	p(w, "\n// implement dispatching\n\n")
	mustWrite(w, ws.WriteDispatch.Bytes())
	return nil
}

func checkService(kinds []string, obj *types.TypeName, cmt generator.Comment) (kind string, _ *types.Interface, err error) {
	name := obj.Name()
	for _, suffix := range kinds {
		if strings.HasSuffix(name, suffix) {
			kind = suffix
			break
		}
	}
	if kind == "" {
		return
	}

	if obj == nil {
		return "", nil, generator.Errorf(nil, "%v: can not load definition", name)
	}
	typ := obj.Type()
	if typ == nil {
		return "", nil, generator.Errorf(nil, "%v: can not load type information", name)
	}
	if typ, ok := typ.(*types.Named); ok {
		if typ, ok := typ.Underlying().(*types.Interface); ok {
			return kind, typ, nil
		}
	}
	return "", nil, generator.Errorf(nil, "%v: must be an interface", name)
}
