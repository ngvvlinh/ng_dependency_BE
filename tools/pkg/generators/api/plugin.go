package api

import (
	"etop.vn/backend/tools/pkg/generators/api/defs"
	"etop.vn/backend/tools/pkg/generators/api/parse"

	"golang.org/x/tools/go/packages"

	"etop.vn/backend/tools/pkg/generator"
)

func New() generator.Plugin {
	return &gen{
		Filterer: generator.FilterByCommand("gen:api"),
	}
}

type gen struct {
	generator.Filterer
}

func (g *gen) Name() string { return "api" }

func (g *gen) Generate(ng generator.Engine) error {
	return ng.GenerateEachPackage(generatePackage)
}

func generatePackage(ng generator.Engine, pkg *packages.Package, printer generator.Printer) error {
	kinds := []defs.Kind{defs.KindQuery, defs.KindAggregate}
	services, err := parse.Services(ng, pkg, kinds)
	if err != nil {
		return err
	}

	w := NewWriter(pkg.Name, pkg.PkgPath, printer, printer)
	ws := &MultiWriter{Writer: w}
	writeCommonDeclaration(ws)
	for _, item := range services {
		debugf("processing service %v", item.Name)
		switch item.Kind {
		case defs.KindQuery:
			generateQueries(ws, item.Name, item.Methods)
		case defs.KindAggregate:
			generateCommands(ws, item.Name, item.Methods)
		}
	}

	p(w, "\n// implement interfaces\n\n")
	mustWrite(w, ws.WriteIface.Bytes())
	p(w, "\n// implement conversion\n\n")
	mustWrite(w, ws.WriteArgs.Bytes())
	p(w, "\n// implement dispatching\n\n")
	mustWrite(w, ws.WriteDispatch.Bytes())
	return nil
}
