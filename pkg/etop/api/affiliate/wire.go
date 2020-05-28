package affiliate

import "github.com/google/wire"

var WireSet = wire.NewSet(
	wire.Struct(new(AccountService), "*"),
	wire.Struct(new(MiscService), "*"),
	NewServers,
)
