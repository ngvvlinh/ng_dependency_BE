package shipping_provider

import (
	"time"

	"etop.vn/backend/pkg/etop/model"
	vtpostclient "etop.vn/backend/pkg/integration/vtpost/client"
)

func Convert_vtpost_ClientStates_To_model_ShippingSourceInternal(in *vtpostclient.ClientStates) *model.ShippingSourceInternal {
	return &model.ShippingSourceInternal{
		ID:        0,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},

		LastSyncAt:  in.AccessTokenCreatedAt,
		AccessToken: in.AccessToken,
		ExpiresAt:   in.ExpiresAt,
		Secret: &model.ShippingSourceSecret{
			CustomerID:     in.CustomerID,
			GroupAddressID: in.GroupAddressID,
		},
	}
}

func Convert_model_ShippingSourceInternal_To_vtpost_ClientStates(in *model.ShippingSourceInternal) *vtpostclient.ClientStates {
	return &vtpostclient.ClientStates{
		AccessToken:          in.AccessToken,
		ExpiresAt:            in.ExpiresAt,
		AccessTokenCreatedAt: in.LastSyncAt,
		CustomerID:           in.Secret.CustomerID,
		GroupAddressID:       in.Secret.GroupAddressID,
	}
}
