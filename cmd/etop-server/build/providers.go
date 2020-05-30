package build

import (
	"context"

	"github.com/google/wire"

	"o.o/api/main/location"
	"o.o/api/top/types/etc/connection_type"
	"o.o/backend/cmd/etop-server/config"
	loggingpayment "o.o/backend/com/etc/logging/payment"
	"o.o/backend/com/etc/logging/smslog"
	"o.o/backend/com/external/payment/manager"
	"o.o/backend/com/external/payment/payment"
	paymentvtpay "o.o/backend/com/external/payment/vtpay"
	"o.o/backend/com/handler/etop-handler/intctl"
	com "o.o/backend/com/main"
	"o.o/backend/com/main/authorization"
	"o.o/backend/com/main/invitation"
	"o.o/backend/com/main/moneytx"
	shipnowcarrier "o.o/backend/com/main/shipnow-carrier"
	shippingcarrier "o.o/backend/com/main/shipping/carrier"
	"o.o/backend/com/services/affiliate"
	"o.o/backend/com/web/webserver"
	"o.o/backend/pkg/common/apifw/whitelabel/drivers"
	"o.o/backend/pkg/common/cmenv"
	"o.o/backend/pkg/common/mq"
	"o.o/backend/pkg/common/redis"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/etop/api"
	"o.o/backend/pkg/etop/api/admin"
	affapi "o.o/backend/pkg/etop/api/affiliate"
	"o.o/backend/pkg/etop/api/integration"
	"o.o/backend/pkg/etop/api/sadmin"
	"o.o/backend/pkg/etop/api/shop"
	shopimports "o.o/backend/pkg/etop/api/shop/imports"
	"o.o/backend/pkg/etop/apix/partner"
	"o.o/backend/pkg/etop/apix/partnercarrier"
	"o.o/backend/pkg/etop/apix/partnerimport"
	xshop "o.o/backend/pkg/etop/apix/shop"
	"o.o/backend/pkg/etop/apix/shopping"
	"o.o/backend/pkg/etop/apix/webhook"
	"o.o/backend/pkg/etop/authorize/middleware"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/backend/pkg/etop/authorize/tokens"
	"o.o/backend/pkg/etop/eventstream"
	imcsvghtk "o.o/backend/pkg/etop/logic/money-transaction/ghtk-imcsv"
	imcsvghn "o.o/backend/pkg/etop/logic/money-transaction/imcsv"
	vtpostimxlsx "o.o/backend/pkg/etop/logic/money-transaction/vtpost-imxlsx"
	logicorder "o.o/backend/pkg/etop/logic/orders"
	orderimcsv "o.o/backend/pkg/etop/logic/orders/imcsv"
	productimcsv "o.o/backend/pkg/etop/logic/products/imcsv"
	"o.o/backend/pkg/etop/logic/shipping_provider"
	"o.o/backend/pkg/etop/middlewares"
	"o.o/backend/pkg/etop/model"
	"o.o/backend/pkg/etop/upload"
	"o.o/backend/pkg/integration/payment/vtpay"
	"o.o/backend/pkg/integration/shipnow/ahamove"
	ahamovewebhook "o.o/backend/pkg/integration/shipnow/ahamove/webhook"
	"o.o/backend/pkg/integration/shipping/ghn"
	ghnwebhook "o.o/backend/pkg/integration/shipping/ghn/webhook"
	"o.o/backend/pkg/integration/shipping/ghtk"
	ghtkwebhook "o.o/backend/pkg/integration/shipping/ghtk/webhook"
	"o.o/backend/pkg/integration/shipping/vtpost"
	vtpostwebhook "o.o/backend/pkg/integration/shipping/vtpost/webhook"
	"o.o/backend/pkg/integration/sms"
	imgroupsms "o.o/backend/pkg/integration/sms/imgroup"
	"o.o/backend/pkg/integration/sms/mock"
	"o.o/backend/pkg/integration/sms/vietguys"
	serviceaffapi "o.o/backend/pkg/services/affiliate/api"
	"o.o/backend/tools/pkg/acl"
	"o.o/capi/httprpc"
	"o.o/common/l"
)

var ll = l.New()

var WireSet = wire.NewSet(
	smslog.WireSet,
	sms.WireSet,
	authorization.WireSet,
	invitation.WireSet,
	affiliate.WireSet,
	logicorder.WireSet,
	moneytx.WireSet,
	payment.WireSet,
	orderimcsv.WireSet,
	productimcsv.WireSet,
	vtpay.WireSet,
	paymentvtpay.WireSet,
	eventstream.WireSet,
	shopimports.WireSet,
	shipnowcarrier.WireSet,
	loggingpayment.WireSet,
	shopping.WireSet,
	ghnwebhook.WireSet,
	ghtkwebhook.WireSet,
	vtpostwebhook.WireSet,
	ahamovewebhook.WireSet,
	wire.FieldsOf(new(Databases), "main", "log", "notifier", "affiliate"),
	BindDatabases,
)

type Databases struct {
	Main      com.MainDB
	Log       com.LogDB
	Notifier  com.NotifierDB
	Affiliate com.AffiliateDB
	WebServer webserver.WebServerDB
}

func BindDatabases(cfg config.Config) (dbs Databases, err error) {
	dbs.Main, err = cmsql.Connect(cfg.Postgres)
	if err != nil {
		return dbs, err
	}
	dbs.Log, err = cmsql.Connect(cfg.PostgresLogs)
	if err != nil {
		return dbs, err
	}
	dbs.Notifier, err = cmsql.Connect(cfg.PostgresNotifier)
	if err != nil {
		return dbs, err
	}
	dbs.Affiliate, err = cmsql.Connect(cfg.PostgresAffiliate)
	if err != nil {
		return dbs, err
	}
	dbs.WebServer, err = cmsql.Connect(cfg.PostgresWebServer)
	if err != nil {
		return dbs, err
	}
	return dbs, nil
}

var SupportedCarriers = wire.NewSet(
	ghn.New,
	ghtk.New,
	vtpost.New,
	imcsvghn.WireSet,
	imcsvghtk.WireSet,
	vtpostimxlsx.WireSet,
	ahamove.WireSet,
)

type EtopHandlers struct {
	IntHandlers []httprpc.Server
	ExtHandlers []httprpc.Server
}

func BindProducer(ctx context.Context, cfg config.Config) (_ webhook.Producer, _ error) {
	var producer *mq.KafkaProducer
	if !cfg.Kafka.Enabled {
		return nil, nil
	}
	producer, err := mq.NewKafkaProducer(ctx, cfg.Kafka.Brokers)
	if err != nil {
		return nil, err
	}
	ctlProducer := producer.WithTopic(intctl.Topic(cfg.Kafka.TopicPrefix))
	return ctlProducer, err
}

func NewSession(cfg config.Config, redisStore redis.Store) session.Session {
	return session.New(
		session.OptValidator(tokens.NewTokenStore(redisStore)),
		session.OptSuperAdmin(cfg.SAdminToken),
	)
}

func WireSAdminToken(cfg config.Config) middleware.SAdminToken {
	return middleware.SAdminToken(cfg.SAdminToken)
}

func NewEtopHandlers(
	rootServers api.Servers,
	shopServers shop.Servers,
	adminServers admin.Servers,
	sadminServers sadmin.Servers,
	integrationServers integration.Servers,

	affServer affapi.Servers,
	saffServer serviceaffapi.Servers,
	partnerServers partner.Servers,
	xshopServers xshop.Servers,
	carrierServers partnercarrier.Servers,
	partnerImportServers partnerimport.Servers,
) (h EtopHandlers) {
	logging := middlewares.NewLogging()
	ssHooks := session.NewHook(acl.GetACL())
	ssExtHooks := session.NewHook(acl.GetExtACL())

	h.IntHandlers = append(h.IntHandlers, rootServers...)
	h.IntHandlers = append(h.IntHandlers, shopServers...)
	h.IntHandlers = append(h.IntHandlers, adminServers...)
	h.IntHandlers = append(h.IntHandlers, httprpc.WithHooks(sadminServers, ssHooks, logging)...)
	h.IntHandlers = append(h.IntHandlers, affServer...)
	h.IntHandlers = append(h.IntHandlers, saffServer...)
	h.IntHandlers = httprpc.WithHooks(h.IntHandlers)

	h.ExtHandlers = append(h.ExtHandlers, partnerServers...)
	h.ExtHandlers = append(h.ExtHandlers, xshopServers...)
	h.ExtHandlers = append(h.ExtHandlers, httprpc.WithHooks(carrierServers, ssExtHooks, logging)...)
	h.ExtHandlers = append(h.ExtHandlers, partnerImportServers...)
	h.ExtHandlers = httprpc.WithHooks(h.ExtHandlers)
	return
}

func NewUploader(cfg config.Config) (*upload.Uploader, error) {
	return upload.NewUploader(map[string]string{
		model.ImportTypeShopOrder.String():   cfg.Upload.DirImportShopOrder,
		model.ImportTypeShopProduct.String(): cfg.Upload.DirImportShopProduct,
	})
}

func AhamoveConfig(cfg config.Config) ahamove.URLConfig {
	return ahamove.URLConfig{
		ThirdPartyHost:       cfg.ThirdPartyHost,
		PathUserVerification: config.PathAhamoveUserVerification,
	}
}

func SupportedShipnowCarriers(
	ahamoveCarrier *ahamove.Carrier,
	ahamoveCarrierAccount *ahamove.CarrierAccount,
) []*shipnowcarrier.Carrier {
	return []*shipnowcarrier.Carrier{
		{
			ShipnowCarrier:        ahamoveCarrier,
			ShipnowCarrierAccount: ahamoveCarrierAccount,
		},
	}
}

func SupportedShippingCarrierConfig(cfg config.Config) shippingcarrier.Config {
	return shippingcarrier.Config{
		Endpoints: []shippingcarrier.ConfigEndpoint{
			{
				connection_type.ConnectionProviderGHN, cfg.GHNWebhook.Endpoint,
			},
		},
	}
}

func SupportedSMSDrivers(mainCfg config.Config, cfg sms.Config) []sms.DriverConfig {
	var imgroupSMSClient *imgroupsms.Client
	if mainCfg.WhiteLabel.IMGroup.SMS.APIKey != "" {
		imgroupSMSClient = imgroupsms.New(mainCfg.WhiteLabel.IMGroup.SMS)
	} else if !cmenv.IsDev() {
		ll.Panic("no sms config for whitelabel/imgroup")
	}

	var mainDriver sms.Driver
	if cfg.Mock {
		mainDriver = mock.GetMock()
	} else {
		mainDriver = vietguys.New(cfg.Vietguys)
	}

	return []sms.DriverConfig{
		{"", mainDriver},
		{drivers.ITopXKey, imgroupSMSClient},
	}
}

func SupportedCarrierDrivers(ctx context.Context, cfg config.Config, locationBus location.QueryBus) []shipping_provider.CarrierDriver {
	var ghnCarrier *ghn.Carrier
	var ghtkCarrier *ghtk.Carrier
	var vtpostCarrier *vtpost.Carrier

	if cfg.GHN.AccountDefault.Token != "" {
		ghnCarrier = ghn.New(cfg.GHN, locationBus)
		if err := ghnCarrier.InitAllClients(ctx); err != nil {
			ll.Fatal("Unable to connect to GHN", l.Error(err))
		}
	} else {
		if cmenv.IsDev() {
			ll.Warn("DEVELOPMENT. Skip connecting to GHN")
		} else {
			ll.Fatal("GHN: No token")
		}
	}

	if cfg.GHTK.AccountDefault.Token != "" {
		ghtkCarrier = ghtk.New(cfg.GHTK, locationBus)
		if err := ghtkCarrier.InitAllClients(ctx); err != nil {
			ll.Fatal("Unable to connect to GHTK", l.Error(err))
		}
	} else {
		if cmenv.IsDev() {
			ll.Warn("DEVELOPMENT. Skip connecting to GHTK.")
		} else {
			ll.Fatal("GHTK: No token")
		}
	}

	if cfg.VTPost.AccountDefault.Username != "" {
		vtpostCarrier = vtpost.New(cfg.VTPost, locationBus)
		if err := vtpostCarrier.InitAllClients(ctx); err != nil {
			ll.Fatal("Unable to connect to VTPost", l.Error(err))
		}
	} else {
		if cmenv.IsDev() {
			ll.Warn("DEVELOPMENT. Skip connecting to VTPost.")
		} else {
			ll.Fatal("VTPost: No token")
		}
	}
	return []shipping_provider.CarrierDriver{ghnCarrier, ghtkCarrier, vtpostCarrier}
}

func SupportedPaymentProvider(
	vtpayProvider *vtpay.Provider,
) []manager.PaymentProvider {
	return []manager.PaymentProvider{
		vtpayProvider,
	}
}
