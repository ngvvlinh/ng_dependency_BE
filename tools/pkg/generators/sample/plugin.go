package sample

import (
	"fmt"

	"o.o/backend/tools/pkg/generator"
)

func New() generator.Plugin {
	return plugin{}
}

var _ generator.Filterer = plugin{}

type plugin struct{}

func (p plugin) Name() string    { return "sample" }
func (p plugin) Command() string { return "sample" }

func (p plugin) Filter(_ generator.FilterEngine) error {
	return nil
}

func (p plugin) Generate(ng generator.Engine) error {
	pkgs := ng.GeneratingPackages()
	for _, gpkg := range pkgs {
		fmt.Printf("package %v:\n", gpkg.Package.PkgPath)
		objects := gpkg.GetObjects()
		for _, obj := range objects {
			fmt.Printf("  %v\t%v\n", obj.Name(), obj.Type())
		}
		fmt.Println()
	}
	return nil
}
