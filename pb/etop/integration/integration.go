package integration

import (
	"etop.vn/backend/pb/etop"
	"etop.vn/backend/pb/etop/etc/status3"
	"etop.vn/backend/pkg/etop/model"
)

func PbPartnerUserInfo(m *model.User) *PartnerUserLogin {
	if m == nil {
		return nil
	}
	return &PartnerUserLogin{
		Id:        m.ID,
		FullName:  m.FullName,
		ShortName: m.ShortName,
		Phone:     m.Phone,
		Email:     m.Email,
	}
}

func PbPartnerShopInfo(m *model.Shop) *PartnerShopInfo {
	return &PartnerShopInfo{
		Id:     m.ID,
		Name:   m.Name,
		Status: status3.Pb(m.Status),
		IsTest: m.IsTest != 0,
		// Address:           m.Address,
		// Phone:             m.
		ImageUrl: m.ImageURL,
		// Email:             "",
		ShipFromAddressId: 0,
	}
}

func PbPartnerShopLoginAccount(m *model.Shop, accessToken string, expiresIn int32) *PartnerShopLoginAccount {
	return &PartnerShopLoginAccount{
		Id:          m.ID,
		Name:        m.Name,
		Type:        etop.PbAccountType(model.TypeShop),
		AccessToken: accessToken,
		ExpiresIn:   expiresIn,
		ImageUrl:    m.ImageURL,
	}
}
