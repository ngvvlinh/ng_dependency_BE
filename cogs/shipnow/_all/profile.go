package shipnow_all

import (
	"github.com/google/wire"

	carriertypes "o.o/api/main/shipnow/carrier/types"
	"o.o/backend/cmd/etop-server/config"
	_shipnow "o.o/backend/cogs/shipnow"
	_ahamove "o.o/backend/cogs/shipnow/ahamove"
	"o.o/backend/com/main/shipnowcarrier"
	"o.o/backend/pkg/integration/shipnow/ahamove"
	"o.o/backend/pkg/integration/shipnow/ahamove/server"
)

var WireSet = wire.NewSet(
	_shipnow.WireSet,
	_ahamove.WireSet,
	AhamoveConfig,
	AllSupportedShipnowCarriers,
)

func AhamoveConfig(cfg config.Config) ahamove.URLConfig {
	return ahamove.URLConfig{
		ThirdPartyHost:       cfg.ThirdPartyHost,
		PathUserVerification: server.PathAhamoveUserVerification,
	}
}

func AllSupportedShipnowCarriers(
	ahamoveCarrier *ahamove.Carrier,
	ahamoveCarrierAccount *ahamove.CarrierAccount,
) []*shipnowcarrier.Carrier {
	return []*shipnowcarrier.Carrier{
		{
			Code:                  carriertypes.Ahamove,
			ShipnowCarrier:        ahamoveCarrier,
			ShipnowCarrierAccount: ahamoveCarrierAccount,
		},
	}
}
