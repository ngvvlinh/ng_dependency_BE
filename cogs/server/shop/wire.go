// +build wireinject

package server_shop

import (
	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	BuildImportHandler,
	BuildEventStreamHandler,
	BuildDownloadHandler,
)
