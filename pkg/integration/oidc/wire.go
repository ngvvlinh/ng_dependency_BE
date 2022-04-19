package oidc

import (
	"github.com/google/wire"
	"o.o/backend/pkg/integration/oidc/client"
)

var WireSet = wire.NewSet(
	client.New,
)
