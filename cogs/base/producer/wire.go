// +build wireinject

package _producer

import (
	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	SupportedProducers,
)
