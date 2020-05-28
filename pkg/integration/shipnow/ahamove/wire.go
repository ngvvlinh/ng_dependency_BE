package ahamove

import (
	"github.com/google/wire"

	"o.o/backend/pkg/integration/shipnow/ahamove/client"
)

var WireSet = wire.NewSet(
	client.New,
	NewCarrierAccount,
	New,
)
