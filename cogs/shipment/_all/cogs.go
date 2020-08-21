package shipment_all

import (
	"context"
	"strconv"

	"o.o/api/main/connectioning"
	"o.o/api/main/location"
	"o.o/api/top/types/etc/connection_type"
	sptypes "o.o/api/top/types/etc/shipping_provider"
	_ghtk "o.o/backend/cogs/shipment/ghtk"
	_vtpost "o.o/backend/cogs/shipment/vtpost"
	carriertypes "o.o/backend/com/main/shipping/carrier/types"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/cmenv"
	"o.o/backend/pkg/etop/logic/shipping_provider"
	"o.o/backend/pkg/etop/sqlstore"
	directclient "o.o/backend/pkg/integration/shipping/direct/client"
	directdriver "o.o/backend/pkg/integration/shipping/direct/driver"
	"o.o/backend/pkg/integration/shipping/ghn"
	ghnclient "o.o/backend/pkg/integration/shipping/ghn/client"
	ghnclientv2 "o.o/backend/pkg/integration/shipping/ghn/clientv2"
	ghndriver "o.o/backend/pkg/integration/shipping/ghn/driver"
	ghndriverv2 "o.o/backend/pkg/integration/shipping/ghn/driverv2"
	"o.o/backend/pkg/integration/shipping/ghtk"
	ghtkclient "o.o/backend/pkg/integration/shipping/ghtk/client"
	ghtkdriver "o.o/backend/pkg/integration/shipping/ghtk/driver"
	"o.o/backend/pkg/integration/shipping/services"
	"o.o/backend/pkg/integration/shipping/vtpost"
	vtpostclient "o.o/backend/pkg/integration/shipping/vtpost/client"
	vtpostdriver "o.o/backend/pkg/integration/shipping/vtpost/driver"
	"o.o/common/l"
)

var ll = l.New()

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

func SupportedShippingCarrierConfig(cfg Config) carriertypes.Config {
	return carriertypes.Config{
		Endpoints: []carriertypes.ConfigEndpoint{
			{
				connection_type.ConnectionProviderGHN, cfg.GHNWebhook.Endpoint,
			},
		},
	}
}

type CarrierDriver struct {
}

func SupportedCarrierDriver() carriertypes.Driver {
	return CarrierDriver{}
}

func (d CarrierDriver) GetShipmentDriver(
	env string, locationQS location.QueryBus,
	connection *connectioning.Connection,
	shopConnection *connectioning.ShopConnection,
	endpoints carriertypes.ConfigEndpoints,
) (carriertypes.ShipmentCarrier, error) {
	etopAffiliateAccount := connection.EtopAffiliateAccount

	switch connection.ConnectionProvider {
	case connection_type.ConnectionProviderGHN:
		version := connection.Version
		if shopConnection.ExternalData == nil {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "Connection ExternalData is missing (connection_id = %v, shop_id = %v)", connection.ID, shopConnection.ShopID)
		}
		userID := shopConnection.ExternalData.UserID
		clientID, err := strconv.Atoi(userID)
		if err != nil {
			return nil, cm.Errorf(cm.InvalidArgument, err, "Connection ExternalData: UserID is invalid")
		}

		switch version {
		case "v2":
			shopIDConnectionStr := shopConnection.ExternalData.ShopID
			shopIDConnection, err := strconv.Atoi(shopIDConnectionStr)
			if err != nil {
				return nil, cm.Errorf(cm.InvalidArgument, err, "Connection ExternalData: ShopID is invalid")
			}
			cfg := ghnclientv2.GHNAccountCfg{
				ClientID: clientID,
				ShopID:   shopIDConnection,
				Token:    shopConnection.Token,
			}
			if affiliateID, err := strconv.Atoi(etopAffiliateAccount.UserID); err == nil {
				cfg.AffiliateID = affiliateID
			}
			driver := ghndriverv2.New(env, cfg, locationQS)
			return driver, nil
		default:
			cfg := ghnclient.GHNAccountCfg{
				ClientID: clientID,
				Token:    shopConnection.Token,
			}
			if etopAffiliateAccount != nil {
				if affiliateID, err := strconv.Atoi(etopAffiliateAccount.UserID); err == nil {
					cfg.AffiliateID = affiliateID
				}
			}
			webhookEndpoint, ok := endpoints.GetByCarrier(connection_type.ConnectionProviderGHN)
			if !ok {
				return nil, cm.Errorf(cm.Internal, nil, "no webhook endpoint")
			}
			driver := ghndriver.New(env, cfg, locationQS, webhookEndpoint)
			return driver, nil
		}

	case connection_type.ConnectionProviderGHTK:
		cfg := ghtkclient.GhtkAccount{
			Token: shopConnection.Token,
		}
		if etopAffiliateAccount != nil {
			cfg.AffiliateID = etopAffiliateAccount.UserID
			cfg.B2CToken = etopAffiliateAccount.Token
		}
		driver := ghtkdriver.New(env, cfg, locationQS)
		return driver, nil

	case connection_type.ConnectionProviderVTP:
		driver := vtpostdriver.New(env, shopConnection.Token, locationQS)
		return driver, nil

	case connection_type.ConnectionProviderPartner:
		cfg := directclient.PartnerAccountCfg{
			Token:      shopConnection.Token,
			Connection: connection,
		}
		if etopAffiliateAccount != nil {
			cfg.AffiliateID = etopAffiliateAccount.UserID
		}
		return directdriver.New(locationQS, cfg)
	default:
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Connection không hợp lệ")
	}
}

func (d CarrierDriver) GetAffiliateShipmentDriver(env string, locationQS location.QueryBus,
	conn *connectioning.Connection,
	endpoints carriertypes.ConfigEndpoints) (carriertypes.ShipmentCarrier, error) {
	var userID, token, shopIDStr string
	version := conn.Version
	if conn.EtopAffiliateAccount != nil {
		userID = conn.EtopAffiliateAccount.UserID
		token = conn.EtopAffiliateAccount.Token
		shopIDStr = conn.EtopAffiliateAccount.ShopID
	}

	switch conn.ConnectionProvider {
	case connection_type.ConnectionProviderGHN:
		clientID, err := strconv.Atoi(userID)
		if err != nil {
			return nil, cm.Errorf(cm.InvalidArgument, err, "AffiliateAccount: UserID is invalid")
		}
		switch version {
		case "v2":
			shopID, err := strconv.Atoi(shopIDStr)
			if err != nil {
				return nil, cm.Errorf(cm.InvalidArgument, err, "AffiliateAcount: ShopID is invalid")
			}
			cfg := ghnclientv2.GHNAccountCfg{
				ClientID:    clientID,
				ShopID:      shopID,
				Token:       token,
				AffiliateID: clientID,
			}
			driver := ghndriverv2.New(env, cfg, locationQS)
			return driver, nil
		default:
			cfg := ghnclient.GHNAccountCfg{
				ClientID: clientID,
				Token:    token,
			}
			webhookEndpoint, ok := endpoints.GetByCarrier(connection_type.ConnectionProviderGHN)
			if !ok {
				return nil, cm.Errorf(cm.Internal, nil, "no webhook endpoint")
			}
			driver := ghndriver.New(env, cfg, locationQS, webhookEndpoint)
			return driver, nil
		}
	case connection_type.ConnectionProviderGHTK:
		cfg := ghtkclient.GhtkAccount{
			AccountID: userID,
			Token:     token,
		}
		driver := ghtkdriver.New(env, cfg, locationQS)
		return driver, nil
	case connection_type.ConnectionProviderPartner:
		cfg := directclient.PartnerAccountCfg{
			Connection: conn,
			Token:      token,
		}
		return directdriver.New(locationQS, cfg)
	default:
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Connection không hỗ trợ affiliate account (connType = %v, connName = %v)", conn.ConnectionProvider, conn.Name)
	}
}

func SupportedShipmentServices() services.MapShipmentServices {
	return shipmentServicesByCarrier
}

var shipmentServicesByCarrier = services.MapShipmentServices{
	sptypes.GHN: {
		{
			ServiceID: string(ghnclient.ServiceFee6Hours),
			Name:      "Gói 6 giờ",
		}, {
			ServiceID: string(ghnclient.ServiceFee1Day),
			Name:      "Gói 1 ngày",
		}, {
			ServiceID: string(ghnclient.ServiceFee2Days),
			Name:      "Gói 2 ngày",
		}, {
			ServiceID: string(ghnclient.ServiceFee3Days),
			Name:      "Gói 3 ngày",
		}, {
			ServiceID: string(ghnclient.ServiceFee4Days),
			Name:      "Gói 4 ngày",
		}, {
			ServiceID: string(ghnclient.ServiceFee5Days),
			Name:      "Gói 5 ngày",
		}, {
			ServiceID: string(ghnclient.ServiceFee6Days),
			Name:      "Gói 6 ngày",
		},
	},
	sptypes.GHTK: {
		{
			ServiceID: string(ghtkclient.TransportRoad),
			Name:      "Đường bộ",
		}, {
			ServiceID: string(ghtkclient.TransportFly),
			Name:      "Đường hàng không",
		},
	},
	sptypes.VTPost: {
		{
			ServiceID: string(vtpostclient.OrderServiceCodeSCOD),
			Name:      "Nhanh - SCOD Giao hàng thu tiền",
		}, {
			ServiceID: string(vtpostclient.OrderServiceCodeVCN),
			Name:      "Nhanh - VCN Chuyển phát nhanh - Express dilivery",
		}, {
			ServiceID: string(vtpostclient.OrderServiceCodeVTK),
			Name:      "Chậm - VTK - VTK Tiết kiệm - Express Saver",
		}, {
			ServiceID: string(vtpostclient.OrderServiceCodePHS),
			Name:      "Chậm - PHS Phát hôm sau nội tỉnh",
		}, {
			ServiceID: string(vtpostclient.OrderServiceCodeVVT),
			Name:      "Chậm - VVT Dịch vụ vận tải",
		}, {
			ServiceID: string(vtpostclient.OrderServiceCodeVHT),
			Name:      "Nhanh - VHT Phát Hỏa tốc",
		}, {
			ServiceID: string(vtpostclient.OrderServiceCodePTN),
			Name:      "Nhanh - PTN Phát trong ngày nội tỉnh",
		}, {
			ServiceID: string(vtpostclient.OrderServiceCodePHT),
			Name:      "Nhanh - PHT Phát hỏa tốc nội tỉnh",
		}, {
			ServiceID: string(vtpostclient.OrderServiceCodeVBS),
			Name:      "Nhanh - VBS Nhanh theo hộp",
		}, {
			ServiceID: string(vtpostclient.OrderServiceCodeVBE),
			Name:      "Chậm - VBE Tiết kiệm theo hộp",
		},
	},
}
