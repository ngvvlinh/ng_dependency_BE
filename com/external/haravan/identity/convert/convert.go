package convert

import (
	"etop.vn/api/external/haravan/identity"
	identitymodel "etop.vn/backend/com/external/haravan/identity/model"
)

func XAccountHaravan(in *identitymodel.ExternalAccountHaravan) (out *identity.ExternalAccountHaravan) {
	if in == nil {
		return nil
	}
	return &identity.ExternalAccountHaravan{
		ID:                                in.ID,
		ShopID:                            in.ShopID,
		Subdomain:                         in.Subdomain,
		ExternalShopID:                    in.ExternalShopID,
		AccessToken:                       in.AccessToken,
		ExternalCarrierServiceID:          in.ExternalCarrierServiceID,
		ExternalConnectedCarrierServiceAt: in.ExternalConnectedCarrierServiceAt,
		ExpiresAt:                         in.ExpiresAt,
		CreatedAt:                         in.CreatedAt,
		UpdatedAt:                         in.UpdatedAt,
	}
}
