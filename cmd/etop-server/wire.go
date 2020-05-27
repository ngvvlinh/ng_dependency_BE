// +build wireinject

package main

import (
	"github.com/google/wire"

	"o.o/api/external/payment/manager"
	"o.o/api/main/shipnow/carrier"
	"o.o/backend/cmd/etop-server/config"
	"o.o/backend/com/etc/logging/smslog"
	"o.o/backend/com/main/authorization"
	"o.o/backend/com/main/invitation"
	"o.o/backend/com/services/affiliate"
	cmservice "o.o/backend/pkg/common/apifw/service"
	"o.o/backend/pkg/common/authorization/auth"
	"o.o/backend/pkg/common/extservice/telebot"
	"o.o/backend/pkg/common/redis"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/etop/api"
	"o.o/backend/pkg/etop/api/shop"
	"o.o/backend/pkg/etop/apix/partner"
	"o.o/backend/pkg/etop/apix/partnercarrier"
	xshipping "o.o/backend/pkg/etop/apix/shipping"
	xshop "o.o/backend/pkg/etop/apix/shop"
	"o.o/backend/pkg/etop/authorize/session"
	logicorder "o.o/backend/pkg/etop/logic/orders"
	affapi "o.o/backend/pkg/services/affiliate/api"
	"o.o/capi"
	"o.o/capi/httprpc"
)

var WireSet = wire.NewSet(
	smslog.WireSet,
	authorization.WireSet,
	invitation.WireSet,
	affiliate.WireSet,
	logicorder.WireSet,
)

func BuildServers(
	db *cmsql.Database,
	cfg config.Config,
	bot *telebot.Channel,
	sd cmservice.Shutdowner,
	eventBus capi.EventBus,
	rd redis.Store,
	s auth.Generator,
	ss *session.Session,
	shipnowCarrierManager carrier.Manager,
	paymentManager manager.CommandBus,
	authURL partner.AuthURL,
) []httprpc.Server {
	panic(wire.Build(
		wire.FieldsOf(&cfg, "email", "sms", "invitation", "secret"),
		wire.FieldsOf(&cfg, "FlagApplyShipmentPrice"),
		WireSet,
		api.WireSet,
		shop.WireSet,
		shop.WireDepsSet,
		affapi.WireSet,
		partner.WireSet,
		partnercarrier.WireSet,
		xshop.WireSet,
		xshipping.WireSet,
		SupportedCarrierDrivers,
		NewServers,
	))
}
