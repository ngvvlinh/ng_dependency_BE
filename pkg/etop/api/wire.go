// +build wireinject

package api

import (
	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	wire.Struct(new(AccountService), "*"),
	wire.Struct(new(AccountRelationshipService), "*"),
	wire.Struct(new(AddressService), "*"),
	wire.Struct(new(BankService), "*"),
	wire.Struct(new(EcomService), "*"),
	wire.Struct(new(LocationService), "*"),
	wire.Struct(new(MiscService), "*"),
	wire.Struct(new(UserService), "*"),
	wire.Struct(new(UserRelationshipService), "*"),
	wire.Struct(new(TicketService), "*"),
	NewServers,
)
