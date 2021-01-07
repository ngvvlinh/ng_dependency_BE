// +build wireinject

package server_max

import (
	"github.com/google/wire"
	"o.o/backend/pkg/etop/apix/authx"
)

var WireSet = wire.NewSet(
	BuildIntHandlers,
	BuildExtHandlers,
	authx.WireSet,
	BuildAuthxHandler,
)
