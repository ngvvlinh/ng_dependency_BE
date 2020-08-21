// +build wireinject

package summary

import "github.com/google/wire"

var WireSet = wire.NewSet(
	New,
)
