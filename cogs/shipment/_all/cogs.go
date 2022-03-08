package shipment_all

import (
	_ntx "o.o/backend/cogs/shipment/ntx"
	ntxclient "o.o/backend/pkg/integration/shipping/ntx/client"
	"strconv"

	"o.o/api/main/connectioning"
	"o.o/api/main/identity"
	"o.o/api/main/location"
	"o.o/api/main/shippingcode"
	"o.o/api/top/types/etc/connection_type"
	sptypes "o.o/api/top/types/etc/shipping_provider"
	_ghtk "o.o/backend/cogs/shipment/ghtk"
	_vtpost "o.o/backend/cogs/shipment/vtpost"
	carriertypes "o.o/backend/com/main/shipping/carrier/types"
	cm "o.o/backend/pkg/common"
	dhlclient "o.o/backend/pkg/integration/shipping/dhl/client"
	dhldriver "o.o/backend/pkg/integration/shipping/dhl/driver"
	directclient "o.o/backend/pkg/integration/shipping/direct/client"
	directdriver "o.o/backend/pkg/integration/shipping/direct/driver"
	"o.o/backend/pkg/integration/shipping/ghn"
	ghnclient "o.o/backend/pkg/integration/shipping/ghn/client"
	ghnclientv2 "o.o/backend/pkg/integration/shipping/ghn/clientv2"
	ghndriver "o.o/backend/pkg/integration/shipping/ghn/driver"
	ghndriverv2 "o.o/backend/pkg/integration/shipping/ghn/driverv2"
	"o.o/backend/pkg/integration/shipping/ghn/driverv2/etop"
	"o.o/backend/pkg/integration/shipping/ghtk"
	ghtkclient "o.o/backend/pkg/integration/shipping/ghtk/client"
	ghtkdriver "o.o/backend/pkg/integration/shipping/ghtk/driver"
	"o.o/backend/pkg/integration/shipping/ninjavan"
	ninjavanclient "o.o/backend/pkg/integration/shipping/ninjavan/client"
	ninjavandriver "o.o/backend/pkg/integration/shipping/ninjavan/driver"
	//ntxclient "o.o/backend/pkg/integration/shipping/ntx/client"
	ntxdriver "o.o/backend/pkg/integration/shipping/ntx/driver"
	"o.o/backend/pkg/integration/shipping/services"
	"o.o/backend/pkg/integration/shipping/vtpost"
	vtpostclient "o.o/backend/pkg/integration/shipping/vtpost/client"
	vtpostdriver "o.o/backend/pkg/integration/shipping/vtpost/driver"
	"o.o/capi"
	"o.o/common/l"
)

var ll = l.New()

type Config struct {
	GHN             ghn.Config             `yaml:"ghn"`
	GHNWebhook      ghn.WebhookConfig      `yaml:"ghn_webhook"`
	GHTK            ghtk.Config            `yaml:"ghtk"`
	GHTKWebhook     _ghtk.WebhookConfig    `yaml:"ghtk_webhook"`
	VTPost          vtpost.Config          `yaml:"vtpost"`
	VTPostWebhook   _vtpost.WebhookConfig  `yaml:"vtpost_webhook"`
	NTX             ntxclient.Config       `yaml:"ntx"`
	NTXWebhook      _ntx.WebhookConfig     `yaml:"ntx_webhook"`
	NinjaVanWebhook ninjavan.WebhookConfig `yaml:"ninjavan_webhook"`
}

func (cfg *Config) MustLoadEnv() {
	cfg.GHN.MustLoadEnv()
	cfg.GHTK.MustLoadEnv()
	cfg.VTPost.MustLoadEnv()
	cfg.NTX.MustLoadEnv()
}

func DefaultConfig() Config {
	return Config{
		GHN:             ghn.DefaultConfig(),
		GHNWebhook:      ghn.DefaultWebhookConfig(),
		GHTK:            ghtk.DefaultConfig(),
		GHTKWebhook:     _ghtk.WebhookConfig{Port: 9032},
		VTPost:          vtpost.DefaultConfig(),
		VTPostWebhook:   _vtpost.WebhookConfig{Port: 9042},
		NinjaVanWebhook: ninjavan.DefaultWebhookConfig(),
		NTX:             ntxclient.DefaultConfig(),
		NTXWebhook:      _ntx.WebhookConfig{Port: 9062},
	}
}

type CarrierDriver struct {
	eventBus capi.EventBus
}

func SupportedCarrierDriver(eventBus capi.EventBus) carriertypes.Driver {
	return CarrierDriver{
		eventBus: eventBus,
	}
}

func (d CarrierDriver) GetShipmentDriver(
	env string, locationQS location.QueryBus,
	identityQS identity.QueryBus, connection *connectioning.Connection,
	shopConnection *connectioning.ShopConnection, shippingcodeQS shippingcode.QueryBus,
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
			supportedEtopGHNDriverV2 := etop.NewEtopSupportedGHNDriver(env, cfg)
			driver := ghndriverv2.New(env, cfg, locationQS, supportedEtopGHNDriverV2)
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
			driver := ghndriver.New(env, cfg, locationQS, "")
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
		driver := vtpostdriver.New(env, shopConnection.Token, locationQS, shippingcodeQS)
		return driver, nil

	case connection_type.ConnectionProviderNinjaVan:
		driver := ninjavandriver.New(env, ninjavanclient.NinjaVanCfg{
			Token: shopConnection.Token,
		}, locationQS, identityQS, shippingcodeQS)
		return driver, nil
	case connection_type.ConnectionProviderNTX:
		driver := ntxdriver.New(ntxclient.Config{
			Username:  etopAffiliateAccount.Username,
			Password:  etopAffiliateAccount.Password,
			PartnerID: etopAffiliateAccount.PartnerID,
		}, locationQS, identityQS, shippingcodeQS)
		return driver, nil

	case connection_type.ConnectionProviderDHL:
		shopIDConnectionStr := shopConnection.ExternalData.ShopID
		if shopIDConnectionStr == "" {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "Connection ExternalData: ShopID is invalid")
		}
		cfg := dhlclient.DHLAccountCfg{
			AccountID: shopIDConnectionStr,
			Token:     shopConnection.Token,
		}
		driver := dhldriver.New(env, cfg, locationQS, d.eventBus)
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

func (d CarrierDriver) GetAffiliateShipmentDriver(
	env string, locationQS location.QueryBus,
	identityQS identity.QueryBus, conn *connectioning.Connection,
	shippingcodeQS shippingcode.QueryBus) (carriertypes.ShipmentCarrier, error) {
	var userID, token, shopIDStr, secretKey string
	version := conn.Version
	if conn.EtopAffiliateAccount != nil {
		userID = conn.EtopAffiliateAccount.UserID
		token = conn.EtopAffiliateAccount.Token
		shopIDStr = conn.EtopAffiliateAccount.ShopID
		secretKey = conn.EtopAffiliateAccount.SecretKey
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
			supportedEtopGHNDriverV2 := etop.NewEtopSupportedGHNDriver(env, cfg)
			driver := ghndriverv2.New(env, cfg, locationQS, supportedEtopGHNDriverV2)
			return driver, nil
		default:
			cfg := ghnclient.GHNAccountCfg{
				ClientID: clientID,
				Token:    token,
			}
			driver := ghndriver.New(env, cfg, locationQS, "")
			return driver, nil
		}
	case connection_type.ConnectionProviderGHTK:
		cfg := ghtkclient.GhtkAccount{
			AccountID: userID,
			Token:     token,
		}
		driver := ghtkdriver.New(env, cfg, locationQS)
		return driver, nil
	case connection_type.ConnectionProviderNinjaVan:
		if userID == "" || secretKey == "" {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "Không thể khởi tạo driver cho NJV")
		}
		cfg := ninjavanclient.NinjaVanCfg{
			ClientID:  userID,
			SecretKey: secretKey,
		}
		driver := ninjavandriver.New(env, cfg, locationQS, identityQS, shippingcodeQS)
		return driver, nil

	case connection_type.ConnectionProviderDHL:
		if userID == "" || secretKey == "" {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "Không thể khởi tạo driver cho DHL: UserID or SecretKey can not be empty")
		}
		cfg := dhlclient.DHLAccountCfg{
			ClientID:     userID,
			ClientSecret: secretKey,
		}
		driver := dhldriver.New(env, cfg, locationQS, d.eventBus)
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
			ServiceID: string(vtpostclient.OrderServiceCodeLCOD),
			Name:      "Chậm - LCOD TMĐT bộ",
		},
		{
			ServiceID: string(vtpostclient.OrderServiceCodeNCOD),
			Name:      "Chậm - NCOD TMĐT bay",
		},
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
	sptypes.NinjaVan: {
		{
			ServiceID: string(ninjavanclient.ServiceLevelStandard),
			Name:      "Chuẩn",
		},
	},
	sptypes.NTX: {
		{
			ServiceID: string(ntxclient.OrderServiceCodeCH),
			Name:      "Chuẩn",
		},
		{
			ServiceID: string(ntxclient.OrderServiceCodeNH),
			Name:      "Nhanh",
		},
	},
	sptypes.DHL: {
		{
			ServiceID: string(dhlclient.OrderServiceCodeSPD),
			Name:      "Giao tối ưu",
		},
		{
			ServiceID: string(dhlclient.OrderServiceCodePDE),
			Name:      "Giao nhanh",
		},
		{
			ServiceID: string(dhlclient.OrderServiceCodePDO),
			Name:      "Giao tiêu chuẩn",
		},
	},
}
