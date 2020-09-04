package shipment_all

import (
	"strconv"

	"o.o/api/main/connectioning"
	"o.o/api/main/location"
	"o.o/api/top/types/etc/connection_type"
	sptypes "o.o/api/top/types/etc/shipping_provider"
	carriertypes "o.o/backend/com/main/shipping/carrier/types"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/integration/shipping/ghn"
	ghnclient "o.o/backend/pkg/integration/shipping/ghn/client"
	ghnclientv2 "o.o/backend/pkg/integration/shipping/ghn/clientv2"
	ghndriver "o.o/backend/pkg/integration/shipping/ghn/driver"
	ghndriverv2 "o.o/backend/pkg/integration/shipping/ghn/driverv2"
	"o.o/backend/pkg/integration/shipping/ghn/driverv2/fabo"
	"o.o/backend/pkg/integration/shipping/services"
	"o.o/common/l"
)

var ll = l.New()

type Config struct {
	GHN        ghn.Config        `yaml:"ghn"`
	GHNWebhook ghn.WebhookConfig `yaml:"ghn_webhook"`
}

func (cfg *Config) MustLoadEnv() {
	cfg.GHN.MustLoadEnv()
}

func DefaultConfig() Config {
	return Config{
		GHN:        ghn.DefaultConfig(),
		GHNWebhook: ghn.DefaultWebhookConfig(),
	}
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
			supportedFaboDriverV2 := fabo.NewFaboSupportedGHNDriver(env, cfg)
			driver := ghndriverv2.New(env, cfg, locationQS, supportedFaboDriverV2)
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

	default:
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Connection không hợp lệ")
	}
}

func (d CarrierDriver) GetAffiliateShipmentDriver(env string, locationQS location.QueryBus,
	conn *connectioning.Connection,
	endpoints carriertypes.ConfigEndpoints,
) (carriertypes.ShipmentCarrier, error) {
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
			supportedFaboGHNDriverV2 := fabo.NewFaboSupportedGHNDriver(env, cfg)
			driver := ghndriverv2.New(env, cfg, locationQS, supportedFaboGHNDriverV2)
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
}
