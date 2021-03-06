package shipment_all

import (
	"strconv"

	"o.o/api/main/connectioning"
	"o.o/api/main/identity"
	"o.o/api/main/location"
	"o.o/api/main/shippingcode"
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

type CarrierDriver struct {
}

func SupportedCarrierDriver() carriertypes.Driver {
	return CarrierDriver{}
}

func (d CarrierDriver) GetShipmentDriver(
	env string, locationQS location.QueryBus,
	identityQS identity.QueryBus,
	connection *connectioning.Connection,
	shopConnection *connectioning.ShopConnection,
	shippingcodeQS shippingcode.QueryBus,
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
			driver := ghndriver.New(env, cfg, locationQS, "")
			return driver, nil
		}

	default:
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Connection kh??ng h???p l???")
	}
}

func (d CarrierDriver) GetAffiliateShipmentDriver(env string, locationQS location.QueryBus,
	identityQS identity.QueryBus,
	conn *connectioning.Connection,
	shippingcodeQS shippingcode.QueryBus,
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
			driver := ghndriver.New(env, cfg, locationQS, "")
			return driver, nil
		}
	default:
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Connection kh??ng h??? tr??? affiliate account (connType = %v, connName = %v)", conn.ConnectionProvider, conn.Name)
	}
}

func SupportedShipmentServices() services.MapShipmentServices {
	return shipmentServicesByCarrier
}

var shipmentServicesByCarrier = services.MapShipmentServices{
	sptypes.GHN: {
		{
			ServiceID: string(ghnclient.ServiceFee6Hours),
			Name:      "G??i 6 gi???",
		}, {
			ServiceID: string(ghnclient.ServiceFee1Day),
			Name:      "G??i 1 ng??y",
		}, {
			ServiceID: string(ghnclient.ServiceFee2Days),
			Name:      "G??i 2 ng??y",
		}, {
			ServiceID: string(ghnclient.ServiceFee3Days),
			Name:      "G??i 3 ng??y",
		}, {
			ServiceID: string(ghnclient.ServiceFee4Days),
			Name:      "G??i 4 ng??y",
		}, {
			ServiceID: string(ghnclient.ServiceFee5Days),
			Name:      "G??i 5 ng??y",
		}, {
			ServiceID: string(ghnclient.ServiceFee6Days),
			Name:      "G??i 6 ng??y",
		},
	},
}
