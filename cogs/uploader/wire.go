// +build wireinject

package _uploader

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewUploader,
)
