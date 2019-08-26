package convert

import (
	etoptypes "etop.vn/api/main/etop"
	"etop.vn/api/main/identity"
	identitymodel "etop.vn/backend/com/main/identity/model"
	"etop.vn/backend/pkg/etop/model"
)

func ShopDB(in *identity.Shop) (out *model.Shop) {
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

func User(in *model.User) (out *identity.User) {
	if in == nil {
		return nil
	}
	out = &identity.User{
		ID:        in.ID,
		FullName:  in.FullName,
		ShortName: in.ShortName,
		Email:     in.Email,
		Phone:     in.Phone,
		Status:    etoptypes.Status3(in.Status),
		CreatedAt: in.CreatedAt,
		UpdatedAt: in.UpdatedAt,
	}
	return
}

func XAccountAhamove(in *identitymodel.ExternalAccountAhamove) *identity.ExternalAccountAhamove {
	return &identity.ExternalAccountAhamove{
		ID:                  in.ID,
		Phone:               in.Phone,
		Name:                in.Name,
		ExternalID:          in.ExternalID,
		ExternalToken:       in.ExternalToken,
		ExternalVerified:    in.ExternalVerified,
		CreatedAt:           in.CreatedAt,
		UpdatedAt:           in.UpdatedAt,
		ExternalCreatedAt:   in.ExternalCreatedAt,
		LastSendVerifiedAt:  in.LastSendVerifiedAt,
		ExternalTicketID:    in.ExternalTicketID,
		IDCardFrontImg:      in.IDCardFrontImg,
		IDCardBackImg:       in.IDCardBackImg,
		PortraitImg:         in.PortraitImg,
		WebsiteURL:          in.WebsiteURL,
		FanpageURL:          in.FanpageURL,
		CompanyImgs:         in.CompanyImgs,
		BusinessLicenseImgs: in.BusinessLicenseImgs,
		UploadedAt:          in.UploadedAt,
	}
}
