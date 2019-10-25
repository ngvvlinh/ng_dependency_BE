package apix

import (
	"go/types"
	"text/template"

	"etop.vn/backend/tools/pkg/generator"
)

var tpl = template.Must(template.New("tpl").Funcs(funcs).Parse(tplText))
var currentPrinter generator.Printer

var funcs = map[string]interface{}{
	"type": renderType,
}

func (p *plugin) generateServices(printer generator.Printer, services []*Service) error {
	currentPrinter = printer
	printer.Import("context", "context")
	printer.Import("fmt", "fmt")
	printer.Import("http", "net/http")
	printer.Import("proto", "github.com/golang/protobuf/proto")
	printer.Import("httprpc", "etop.vn/backend/pkg/common/httprpc")
	vars := map[string]interface{}{
		"Services": services,
	}
	return tpl.Execute(printer, vars)
}

func renderType(typ types.Type) string {
	return currentPrinter.TypeString(typ)
}
