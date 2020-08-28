package shipnow_all

import (
	"o.o/api/main/accountshipnow"
	"o.o/api/main/connectioning"
	"o.o/api/main/identity"
	"o.o/api/main/location"
	shipnowcarriertypes "o.o/api/main/shipnow/carrier/types"
	"o.o/api/top/types/etc/connection_type"
	"o.o/backend/cmd/etop-server/config"
	comshipnowcarriertypes "o.o/backend/com/main/shipnow/carrier/types"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/integration/shipnow/ahamove"
	ahamoveclient "o.o/backend/pkg/integration/shipnow/ahamove/client"
	"o.o/backend/pkg/integration/shipnow/ahamove/server"
)

func SupportedShipnowCarrierConfig(cfg config.Config) comshipnowcarriertypes.Config {
	return comshipnowcarriertypes.Config{
		ThirdPartyHost: cfg.ThirdPartyHost,
		Paths: []comshipnowcarriertypes.ConfigPath{{
			Carrier:              shipnowcarriertypes.Ahamove,
			PathUserVerification: server.PathAhamoveUserVerification,
		}},
	}
}

type ShipnowCarrierDriver struct{}

func SupportedShipnowCarrierDriver() comshipnowcarriertypes.Driver {
	return ShipnowCarrierDriver{}
}

func (d ShipnowCarrierDriver) GetShipnowDriver(
	env string, locationQS location.QueryBus,
	connection *connectioning.Connection,
	shopConnection *connectioning.ShopConnection,
	identityQS identity.QueryBus,
	accountshipnowQS accountshipnow.QueryBus,
	pathCfgs comshipnowcarriertypes.Config,
) (comshipnowcarriertypes.ShipnowCarrier, error) {
	switch connection.ConnectionProvider {
	case connection_type.ConnectionProviderAhamove:
		cfg := ahamoveclient.Config{
			Env:   env,
			Name:  "",
			Token: shopConnection.Token,
		}
		client := ahamoveclient.New(cfg)
		urlCfg := ahamove.URLConfig{
			ThirdPartyHost:       pathCfgs.ThirdPartyHost,
			PathUserVerification: "",
		}
		pathUserVerification, ok := comshipnowcarriertypes.ConfigPaths(pathCfgs.Paths).GetByCarrier(shipnowcarriertypes.Ahamove)
		if !ok {
			return nil, cm.Errorf(cm.Internal, nil, "ahamove: no path user verification")
		}
		urlCfg.PathUserVerification = pathUserVerification

		driver := ahamove.New(client, urlCfg, locationQS, identityQS, accountshipnowQS)
		return driver, nil
	default:
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Connection không hợp lệ")
	}
}

func (d ShipnowCarrierDriver) GetAffiliateShipnowDriver(
	env string, locationQS location.QueryBus,
	connection *connectioning.Connection,
	identityQS identity.QueryBus,
	accountshipnowQS accountshipnow.QueryBus,
) (comshipnowcarriertypes.ShipnowCarrier, error) {
	if connection.EtopAffiliateAccount == nil || connection.EtopAffiliateAccount.Token == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Etop affiliate account không hợp lệ")
	}
	token := connection.EtopAffiliateAccount.Token

	switch connection.ConnectionProvider {
	case connection_type.ConnectionProviderAhamove:
		cfg := ahamoveclient.Config{
			Env:    env,
			Name:   "",
			ApiKey: token,
		}
		client := ahamoveclient.New(cfg)
		driver := ahamove.New(client, ahamove.URLConfig{}, locationQS, identityQS, accountshipnowQS)
		return driver, nil
	default:
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Connection không hợp lệ")
	}
}
