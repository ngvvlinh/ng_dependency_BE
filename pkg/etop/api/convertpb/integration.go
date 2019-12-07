package convertpb

import (
	integration "etop.vn/api/top/int/integration"
	"etop.vn/api/top/types/etc/account_type"
	"etop.vn/backend/pkg/etop/model"
)

func PbPartnerUserInfo(m *model.User) *integration.PartnerUserLogin {
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

func PbPartnerShopInfo(m *model.Shop) *integration.PartnerShopInfo {
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

func PbPartnerShopLoginAccount(m *model.Shop, accessToken string, expiresIn int) *integration.PartnerShopLoginAccount {
	return &integration.PartnerShopLoginAccount{
		Id:          m.ID,
		Name:        m.Name,
		Type:        PbAccountType(account_type.Shop),
		AccessToken: accessToken,
		ExpiresIn:   expiresIn,
		ImageUrl:    m.ImageURL,
	}
}
