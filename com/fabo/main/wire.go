// +build wireinject

package aggregate

import (
	"github.com/google/wire"

	"o.o/backend/com/fabo/main/fbmessaging"
	"o.o/backend/com/fabo/main/fbpage"
	"o.o/backend/com/fabo/main/fbuser"
)

var WireSet = wire.NewSet(
	fbmessaging.WireSet,
	fbpage.WireSet,
	fbuser.WireSet,
)
