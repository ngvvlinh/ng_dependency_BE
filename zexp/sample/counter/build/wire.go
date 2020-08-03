//+build wireinject

package build

import (
	"github.com/google/wire"
	"o.o/backend/pkg/common/lifecycle"
	"o.o/backend/zexp/sample/counter/config"
	"o.o/backend/zexp/sample/counter/service"
)

func Build(cfg config.Config) (lifecycle.HTTPServer, error) {
	wire.Build(
		wire.FieldsOf(
			&cfg,
			"Postgres"),
		service.WireSet,
		BuildServer,
	)
	return lifecycle.HTTPServer{}, nil
}
