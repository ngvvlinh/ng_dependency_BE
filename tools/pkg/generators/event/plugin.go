package event

import (
	"fmt"
	"go/types"
	"io"
	"strings"

	"golang.org/x/tools/go/packages"

	"etop.vn/backend/tools/pkg/generator"
)

const CmdTopic = "gen:event:topic"

func New() generator.Plugin {
	return &gen{
		Filterer: generator.FilterByCommand(CmdTopic),
	}
}

type gen struct {
	generator.Filterer
}

func (g gen) Name() string { return "event" }

func (g gen) Generate(ng generator.Engine) error {
	return ng.GenerateEachPackage(generatePackage)
}

func parseTopic(directives []generator.Directive) (topic string, _ error) {
	for _, d := range directives {
		if d.Cmd == CmdTopic {
			topic = d.Arg
			break
		}
	}
	if topic == "" {
		return "", generator.Errorf(nil, "no topic")
	}
	return
}

func generatePackage(ng generator.Engine, pkg *packages.Package, p generator.Printer) error {
	directives := ng.GetDirectivesByPackage(pkg)
	topic, err := parseTopic(directives)
	if err != nil {
		return err
	}

	for _, object := range ng.GetObjectsByPackage(pkg) {
		switch obj := object.(type) {
		case *types.TypeName:
			if strings.HasSuffix(obj.Name(), "Event") {
				w(p, "func (e *%v) GetTopic() string { return %q }\n",
					obj.Name(), topic)
			}
		}
	}
	return nil
}

func w(w io.Writer, format string, args ...interface{}) {
	_, _ = fmt.Fprintf(w, format, args...)
}
