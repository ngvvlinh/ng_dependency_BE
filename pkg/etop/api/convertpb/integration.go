package convertpb

import (
	integration "o.o/api/top/int/integration"
	"o.o/api/top/types/etc/account_type"
	identitymodel "o.o/backend/com/main/identity/model"
)

func PbPartnerUserInfo(m *identitymodel.User) *integration.PartnerUserLogin {
	if m == nil {
		return nil
	}
	return &integration.PartnerUserLogin{
		Id:        m.ID,
		FullName:  m.FullName,
		ShortName: m.ShortName,
		Phone:     m.Phone,
		Email:     m.Email,
	}
}

func PbPartnerShopInfo(m *identitymodel.Shop) *integration.PartnerShopInfo {
	return &integration.PartnerShopInfo{
		Id:     m.ID,
		Name:   m.Name,
		Status: m.Status,
		IsTest: m.IsTest != 0,
		// Address:           m.Address,
		// Phone:             m.
		ImageUrl: m.ImageURL,
		// Email:             "",
		ShipFromAddressId: 0,
	}
}

func PbPartnerShopLoginAccount(m *identitymodel.Shop, accessToken string, expiresIn int) *integration.PartnerShopLoginAccount {
	return &integration.PartnerShopLoginAccount{
		Id:          m.ID,
		Name:        m.Name,
		Type:        PbAccountType(account_type.Shop),
		AccessToken: accessToken,
		ExpiresIn:   expiresIn,
		ImageUrl:    m.ImageURL,
	}
}
