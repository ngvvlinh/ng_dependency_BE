package convert

import (
	etoptypes "etop.vn/api/main/etop"
	"etop.vn/api/main/identity"
	"etop.vn/backend/pkg/etop/model"
)

func ShopToModel(in *identity.Shop) (out *model.Shop) {
	if in == nil {
		return nil
	}
	out = &model.Shop{
		ID:                            in.ID,
		Name:                          in.Name,
		OwnerID:                       in.OwnerID,
		IsTest:                        in.IsTest,
		AddressID:                     in.AddressID,
		ShipToAddressID:               in.ShipToAddressID,
		ShipFromAddressID:             in.ShipFromAddressID,
		Phone:                         in.Phone,
		BankAccount:                   nil,
		WebsiteURL:                    in.WebsiteURL,
		ImageURL:                      in.ImageURL,
		Email:                         in.Email,
		Code:                          in.Code,
		AutoCreateFFM:                 in.AutoCreateFFM,
		Status:                        model.Status3(in.Status),
		CreatedAt:                     in.CreatedAt,
		UpdatedAt:                     in.UpdatedAt,
		DeletedAt:                     in.DeletedAt,
		Address:                       nil,
		RecognizedHosts:               nil,
		GhnNoteCode:                   "",
		TryOn:                         "",
		CompanyInfo:                   nil,
		MoneyTransactionRRule:         "",
		SurveyInfo:                    nil,
		ShippingServiceSelectStrategy: nil,
	}
	return
}

func Shop(in *model.Shop) (out *identity.Shop) {
	if in == nil {
		return nil
	}
	out = &identity.Shop{
		ID:                in.ID,
		Name:              in.Name,
		OwnerID:           in.OwnerID,
		IsTest:            in.IsTest,
		AddressID:         in.AddressID,
		ShipToAddressID:   in.ShipToAddressID,
		ShipFromAddressID: in.ShipFromAddressID,
		Phone:             in.Phone,
		WebsiteURL:        in.WebsiteURL,
		ImageURL:          in.ImageURL,
		Email:             in.Email,
		Code:              in.Code,
		AutoCreateFFM:     in.AutoCreateFFM,
		Status:            etoptypes.Status3(in.Status),
		CreatedAt:         in.CreatedAt,
		UpdatedAt:         in.UpdatedAt,
		DeletedAt:         in.DeletedAt,
	}
	return
}
