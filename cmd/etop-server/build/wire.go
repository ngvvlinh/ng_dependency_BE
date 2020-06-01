// +build wireinject

package build

import (
	"context"

	"github.com/google/wire"

	"o.o/backend/cmd/etop-server/config"
	"o.o/backend/pkg/common/apifw/captcha"
	"o.o/backend/pkg/common/apifw/health"
	"o.o/backend/pkg/common/authorization/auth"
	"o.o/backend/pkg/common/lifecycle"
	"o.o/backend/pkg/common/redis"
	"o.o/backend/pkg/etop/api"
	"o.o/backend/pkg/etop/api/admin"
	affapi "o.o/backend/pkg/etop/api/affiliate"
	"o.o/backend/pkg/etop/api/export"
	"o.o/backend/pkg/etop/api/integration"
	"o.o/backend/pkg/etop/api/sadmin"
	"o.o/backend/pkg/etop/api/shop"
	"o.o/backend/pkg/etop/apix/partner"
	"o.o/backend/pkg/etop/apix/partnercarrier"
	"o.o/backend/pkg/etop/apix/partnerimport"
	xshipping "o.o/backend/pkg/etop/apix/shipping"
	xshop "o.o/backend/pkg/etop/apix/shop"
	"o.o/backend/pkg/etop/apix/webhook"
	"o.o/backend/pkg/etop/authorize/middleware"
	"o.o/backend/pkg/etop/authorize/tokens"
	"o.o/backend/pkg/etop/eventstream"
	logicsummary "o.o/backend/pkg/etop/logic/summary"
	"o.o/backend/pkg/etop/sqlstore"
	saffapi "o.o/backend/pkg/services/affiliate/api"
	"o.o/capi"
)

func Servers(
	ctx context.Context,
	cfg config.Config,
	eventBus capi.EventBus,
	healthServer *health.Service,
	partnerAuthURL partner.AuthURL,
) ([]lifecycle.HTTPServer, func(), error) {
	panic(wire.Build(
		wire.FieldsOf(&cfg,
			"email",
			"sms",
			"invitation",
			"secret",
			"ghn",
			"ghtk",
			"vtpost",
			"ahamove",
			"vtpay",
			"redis",
			"export",
			"captcha",
		),
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
		export.WireSet,
		webhook.WireSet,
		middleware.WireSet,
		tokens.NewTokenStore,
		redis.Connect,
		auth.NewGenerator,
		logicsummary.WireSet,
		wire.InterfaceValue(new(eventstream.Publisher), new(eventstream.EventStream)),
		sqlstore.WireSet,
		captcha.WireSet,
		WireSAdminToken,
		NewSession,
		BindProducer,
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
