package njvimport

import "github.com/google/wire"

var WireSet = wire.NewSet(
	wire.Struct(new(NJVImporter), "*"),
)
