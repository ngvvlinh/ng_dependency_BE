package sample

import (
	"fmt"

	"etop.vn/backend/tools/pkg/generator"
)

func New() generator.Plugin {
	return plugin{}
}

var _ generator.Filter = plugin{}

type plugin struct{}

func (p plugin) Name() string    { return "sample" }
func (p plugin) Command() string { return "sample" }

func (p plugin) FilterPackage(*generator.PreparsedPackage) (bool, error) {
	return true, nil
}

func (p plugin) Generate(ng generator.Engine) error {
	pkgs := ng.GeneratingPackages()
	for _, gpkg := range pkgs {
		fmt.Printf("package %v:\n", gpkg.Package.PkgPath)
		objects := gpkg.Objects()
		for _, object := range objects {
			obj := object.Object
			fmt.Printf("  %v\t%v\n", obj.Name(), obj.Type())
		}
		fmt.Println()
	}
	return nil
}
