// +build wireinject

package storage_all

import (
	"github.com/google/wire"
)

var WireSet = wire.NewSet(Build)
