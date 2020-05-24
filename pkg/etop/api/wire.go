// +build wireinject

package api

import (
	"github.com/google/wire"
	apisms "o.o/api/etc/logging/smslog"
	"o.o/api/main/authorization"
	"o.o/api/main/identity"
	"o.o/api/main/invitation"
	"o.o/api/main/location"
	cmservice "o.o/backend/pkg/common/apifw/service"
	"o.o/backend/pkg/common/authorization/auth"
	"o.o/backend/pkg/common/extservice/telebot"
	"o.o/backend/pkg/common/redis"
	"o.o/backend/pkg/integration/sms"
	"o.o/capi"
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
	NewServers,
)

func BuildServers(
	bus capi.EventBus,
	locationQuery location.QueryBus,
	smsCommandBus apisms.CommandBus,
	identityCommandBus identity.CommandBus,
	identityQueryBus identity.QueryBus,
	invitationAggr invitation.CommandBus,
	invitationQS invitation.QueryBus,
	authorizationQ authorization.QueryBus,
	authorizationA authorization.CommandBus,
	sd cmservice.Shutdowner,
	rd redis.Store,
	s auth.Generator,
	_cfgEmail EmailConfig,
	_cfgSMS sms.Config,
	bot *telebot.Channel,
) Servers {
	panic(wire.Build(WireSet))
}
