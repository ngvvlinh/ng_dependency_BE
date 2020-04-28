package api

import (
	"golang.org/x/tools/go/packages"

	"o.o/backend/tools/pkg/generator"
	"o.o/backend/tools/pkg/generators/api/defs"
	"o.o/backend/tools/pkg/generators/api/parse"
	"o.o/backend/tools/pkg/genutil"
)

func New() generator.Plugin {
	return &gen{
		Filterer:  generator.FilterByCommand("gen:api"),
		Qualifier: genutil.Qualifier{},
	}
}

type gen struct {
	generator.Filterer
	generator.Qualifier
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

	currentPrinter = printer
	generate(printer, services)
	return nil
}
