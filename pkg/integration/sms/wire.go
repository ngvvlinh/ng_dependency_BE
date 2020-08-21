// +build wireinject

package sms

import "github.com/google/wire"

var WireSet = wire.NewSet(
	New,
)
