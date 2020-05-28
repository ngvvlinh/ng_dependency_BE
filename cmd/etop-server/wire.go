// +build wireinject

package main

import (
	"github.com/google/wire"

	"o.o/backend/cmd/etop-server/config"
	cmservice "o.o/backend/pkg/common/apifw/service"
	"o.o/backend/pkg/common/authorization/auth"
	"o.o/backend/pkg/common/extservice/telebot"
	"o.o/backend/pkg/common/redis"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/etop/api"
	"o.o/backend/pkg/etop/api/admin"
	affapi "o.o/backend/pkg/etop/api/affiliate"
	"o.o/backend/pkg/etop/api/integration"
	"o.o/backend/pkg/etop/api/sadmin"
	"o.o/backend/pkg/etop/api/shop"
	"o.o/backend/pkg/etop/apix/partner"
	"o.o/backend/pkg/etop/apix/partnercarrier"
	"o.o/backend/pkg/etop/apix/partnerimport"
	xshipping "o.o/backend/pkg/etop/apix/shipping"
	xshop "o.o/backend/pkg/etop/apix/shop"
	"o.o/backend/pkg/etop/authorize/session"
	saffapi "o.o/backend/pkg/services/affiliate/api"
	"o.o/capi"
)

func BuildServers(
	db *cmsql.Database,
	cfg config.Config,
	bot *telebot.Channel,
	sd cmservice.Shutdowner,
	eventBus capi.EventBus,
	rd redis.Store,
	s auth.Generator,
	ss session.Session,
	authURL partner.AuthURL,
) ([]Server, error) {
	panic(wire.Build(
		wire.FieldsOf(&cfg, "email", "sms", "invitation", "secret", "ghn", "ghtk", "vtpost", "ahamove", "vtpay"),
		wire.FieldsOf(&cfg, "FlagApplyShipmentPrice"),
		WireSet,
		api.WireSet,
		shop.WireSet,
		shop.WireDepsSet,
		saffapi.WireSet,
		affapi.WireSet,
		partner.WireSet,
		partnercarrier.WireSet,
		xshop.WireSet,
		xshipping.WireSet,
		partnerimport.WireSet,
		admin.WireSet,
		sadmin.WireSet,
		integration.WireSet,
		AhamoveConfig,
		NewAhamoveVerificationFileServer,
		NewUploader,
		SupportedSMSDrivers,
		SupportedCarriers,
		SupportedCarrierDrivers,
		SupportedShipnowCarriers,
		SupportedShippingCarrierConfig,
		SupportedPaymentProvider,
		NewEtopHandlers,
		NewServers,
		NewEtopServer,
		NewWebServer,
		NewGHNWebhookServer,
		NewGHTKWebhookServer,
		NewVTPostWebhookServer,
		NewAhamoveWebhookServer,
	))
}
