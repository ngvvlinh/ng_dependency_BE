// +build wireinject

package build

import (
	"github.com/google/wire"

	"o.o/backend/pkg/common/lifecycle"
	"o.o/backend/zexp/sample/calc3/config"
	"o.o/backend/zexp/sample/calc3/service"
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
