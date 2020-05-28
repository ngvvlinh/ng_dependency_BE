package main

import (
	"github.com/google/wire"

	"o.o/api/main/location"
	"o.o/api/top/types/etc/connection_type"
	"o.o/backend/cmd/etop-server/config"
	loggingpayment "o.o/backend/com/etc/logging/payment"
	"o.o/backend/com/etc/logging/smslog"
	"o.o/backend/com/external/payment/manager"
	"o.o/backend/com/external/payment/payment"
	paymentvtpay "o.o/backend/com/external/payment/vtpay"
	"o.o/backend/com/main/authorization"
	"o.o/backend/com/main/invitation"
	"o.o/backend/com/main/moneytx"
	shipnowcarrier "o.o/backend/com/main/shipnow-carrier"
	shippingcarrier "o.o/backend/com/main/shipping/carrier"
	"o.o/backend/com/services/affiliate"
	"o.o/backend/pkg/common/apifw/whitelabel/drivers"
	"o.o/backend/pkg/common/cmenv"
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
	"o.o/backend/pkg/etop/authorize/session"
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
	"o.o/backend/pkg/integration/shipping/ghn"
	"o.o/backend/pkg/integration/shipping/ghtk"
	"o.o/backend/pkg/integration/shipping/vtpost"
	"o.o/backend/pkg/integration/sms"
	imgroupsms "o.o/backend/pkg/integration/sms/imgroup"
	"o.o/backend/pkg/integration/sms/mock"
	"o.o/backend/pkg/integration/sms/vietguys"
	serviceaffapi "o.o/backend/pkg/services/affiliate/api"
	"o.o/backend/tools/pkg/acl"
	"o.o/capi/httprpc"
	"o.o/common/l"
)

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
)

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

func SupportedCarrierDrivers(cfg config.Config, locationBus location.QueryBus) []shipping_provider.CarrierDriver {
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
