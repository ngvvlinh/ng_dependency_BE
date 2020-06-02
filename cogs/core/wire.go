package _core

import (
	"github.com/google/wire"

	"o.o/backend/com/main/authorization"
	"o.o/backend/com/main/invitation"
)

var WireSet = wire.NewSet(
	authorization.WireSet,
	invitation.WireSet,
)
