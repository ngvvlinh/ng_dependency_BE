package shipnow_all

import (
	carriertypes "o.o/api/main/shipnow/carrier/types"
	"o.o/backend/cmd/etop-server/config"
	"o.o/backend/com/main/shipnowcarrier"
	"o.o/backend/pkg/integration/shipnow/ahamove"
	"o.o/backend/pkg/integration/shipnow/ahamove/server"
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
