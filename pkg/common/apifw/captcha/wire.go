// +build wireinject

package captcha

import "github.com/google/wire"

var WireSet = wire.NewSet(
	New,
)
