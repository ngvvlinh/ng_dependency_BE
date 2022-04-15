// +build wireinject

package server_max

import (
	"github.com/google/wire"
	"o.o/backend/pkg/etop/apix/authx"
	"o.o/backend/pkg/etop/apix/portsip_pbx"
)

var WireSet = wire.NewSet(
	BuildIntHandlers,
	BuildExtHandlers,
	authx.WireSet,
	BuildAuthxHandler,
	portsip_pbx.WireSet,
	BuildPortSipPBXHandler,
)
