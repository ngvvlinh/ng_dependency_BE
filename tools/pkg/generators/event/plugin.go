package event

import (
	"fmt"
	"go/types"
	"io"
	"strings"

	"etop.vn/backend/tools/pkg/generator"
)

const CmdTopic = "gen:event:topic"

func New() generator.Plugin {
	return &gen{
		Filter: generator.FilterByCommand(CmdTopic),
	}
}

type gen struct {
	generator.Filter
}

func (g gen) Name() string { return "event" }

func (g gen) Generate(ng generator.Engine) error {
	pkgs := ng.GeneratingPackages()
	for _, pkg := range pkgs {
		err := generatePackage(pkg)
		if err != nil {
			return err
		}
	}
	return nil
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

func generatePackage(gpkg *generator.GeneratingPackage) error {
	topic, err := parseTopic(gpkg.Directives)
	if err != nil {
		return err
	}

	p := gpkg.Generate()
	for _, object := range gpkg.Objects() {
		switch obj := object.Object.(type) {
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
