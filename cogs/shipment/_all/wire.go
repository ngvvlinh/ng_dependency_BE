package shipment_all

import (
	"context"

	"github.com/google/wire"

	"o.o/api/main/location"
	"o.o/api/top/types/etc/connection_type"
	_shipment "o.o/backend/cogs/shipment"
	_ghn "o.o/backend/cogs/shipment/ghn"
	_ghtk "o.o/backend/cogs/shipment/ghtk"
	_vtpost "o.o/backend/cogs/shipment/vtpost"
	shippingcarrier "o.o/backend/com/main/shipping/carrier"
	"o.o/backend/pkg/common/cmenv"
	imcsvghtk "o.o/backend/pkg/etop/logic/money-transaction/ghtk-imcsv"
	imcsvghn "o.o/backend/pkg/etop/logic/money-transaction/imcsv"
	vtpostimxlsx "o.o/backend/pkg/etop/logic/money-transaction/vtpost-imxlsx"
	"o.o/backend/pkg/etop/logic/shipping_provider"
	"o.o/backend/pkg/etop/sqlstore"
	"o.o/backend/pkg/integration/shipping/ghn"
	"o.o/backend/pkg/integration/shipping/ghtk"
	"o.o/backend/pkg/integration/shipping/vtpost"
	"o.o/common/l"
)

var ll = l.New()

var WireSet = wire.NewSet(
	_shipment.WireSet,
	_ghn.WireSet,
	_ghtk.WireSet,
	_vtpost.WireSet,
	wire.FieldsOf(new(Config), "GHN", "GHNWebhook", "GHTK", "GHTKWebhook", "VTPost", "VTPostWebhook"),
	ghn.New,
	ghtk.New,
	vtpost.New,
	imcsvghn.WireSet,
	imcsvghtk.WireSet,
	vtpostimxlsx.WireSet,
	SupportedCarrierDrivers,
	SupportedShippingCarrierConfig,
)

type Config struct {
	GHN           ghn.Config            `yaml:"ghn"`
	GHNWebhook    ghn.WebhookConfig     `yaml:"ghn_webhook"`
	GHTK          ghtk.Config           `yaml:"ghtk"`
	GHTKWebhook   _ghtk.WebhookConfig   `yaml:"ghtk_webhook"`
	VTPost        vtpost.Config         `yaml:"vtpost"`
	VTPostWebhook _vtpost.WebhookConfig `yaml:"vtpost_webhook"`
}

func (cfg *Config) MustLoadEnv() {
	cfg.GHN.MustLoadEnv()
	cfg.GHTK.MustLoadEnv()
	cfg.VTPost.MustLoadEnv()
}

func DefaultConfig() Config {
	return Config{
		GHN:           ghn.DefaultConfig(),
		GHNWebhook:    ghn.DefaultWebhookConfig(),
		GHTK:          ghtk.DefaultConfig(),
		GHTKWebhook:   _ghtk.WebhookConfig{Port: 9032},
		VTPost:        vtpost.DefaultConfig(),
		VTPostWebhook: _vtpost.WebhookConfig{Port: 9042},
	}
}

// TODO(vu): remove dependence on *sqlstore.Store
func SupportedCarrierDrivers(ctx context.Context, _ *sqlstore.Store, cfg Config, locationBus location.QueryBus) []shipping_provider.CarrierDriver {
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

func SupportedShippingCarrierConfig(cfg Config) shippingcarrier.Config {
	return shippingcarrier.Config{
		Endpoints: []shippingcarrier.ConfigEndpoint{
			{
				connection_type.ConnectionProviderGHN, cfg.GHNWebhook.Endpoint,
			},
		},
	}
}
