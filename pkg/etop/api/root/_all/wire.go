package _all

import (
	"github.com/google/wire"
	apiroot "o.o/backend/pkg/etop/api/root"
)

var WireSet = wire.NewSet(
	wire.Struct(new(apiroot.AccountService), "*"),
	wire.Struct(new(apiroot.AccountRelationshipService), "*"),
	wire.Struct(new(apiroot.AddressService), "*"),
	wire.Struct(new(apiroot.BankService), "*"),
	wire.Struct(new(apiroot.EcomService), "*"),
	wire.Struct(new(apiroot.LocationService), "*"),
	wire.Struct(new(apiroot.MiscService), "*"),
	wire.Struct(new(apiroot.UserService), "*"),
	wire.Struct(new(apiroot.UserRelationshipService), "*"),
	wire.Struct(new(TicketService), "*"),
	NewServers,
)
