package convert

import (
	"etop.vn/api/main/identity"
	"etop.vn/api/top/types/etc/status3"
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
		Status:                        in.Status,
		CreatedAt:                     in.CreatedAt,
		UpdatedAt:                     in.UpdatedAt,
		DeletedAt:                     in.DeletedAt,
		Address:                       nil,
		RecognizedHosts:               nil,
		GhnNoteCode:                   0,
		TryOn:                         0,
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
		Status:            status3.Status(in.Status),
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
		ID:              in.ID,
		FullName:        in.FullName,
		ShortName:       in.ShortName,
		Email:           in.Email,
		Phone:           in.Phone,
		Status:          status3.Status(in.Status),
		EmailVerifiedAt: in.EmailVerifiedAt,
		CreatedAt:       in.CreatedAt,
		UpdatedAt:       in.UpdatedAt,
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

func Affiliate(in *identitymodel.Affiliate) *identity.Affiliate {
	if in == nil {
		return nil
	}
	return &identity.Affiliate{
		ID:          in.ID,
		OwnerID:     in.OwnerID,
		Name:        in.Name,
		Phone:       in.Phone,
		Email:       in.Email,
		Status:      in.Status,
		IsTest:      in.IsTest,
		CreatedAt:   in.CreatedAt,
		UpdatedAt:   in.UpdatedAt,
		DeletedAt:   in.DeletedAt,
		BankAccount: BankAccount(in.BankAccount),
	}
}

func Affiliates(ins []*identitymodel.Affiliate) []*identity.Affiliate {
	var results []*identity.Affiliate
	if len(ins) == 0 {
		return results
	}
	for _, mAff := range ins {
		results = append(results, Affiliate(mAff))
	}
	return results
}

func AffiliateDB(in *identity.Affiliate) *identitymodel.Affiliate {
	if in == nil {
		return nil
	}
	return &identitymodel.Affiliate{
		ID:          in.ID,
		OwnerID:     in.OwnerID,
		Name:        in.Name,
		Phone:       in.Phone,
		Email:       in.Email,
		Status:      status3.Status(in.Status),
		IsTest:      in.IsTest,
		CreatedAt:   in.CreatedAt,
		UpdatedAt:   in.UpdatedAt,
		DeletedAt:   in.DeletedAt,
		BankAccount: BankAccountDB(in.BankAccount),
	}
}

func Permission(in model.Permission) identity.Permission {
	return identity.Permission{
		Roles:       in.Roles,
		Permissions: in.Permissions,
	}
}

func PermissionToModel(in identity.Permission) model.Permission {
	return model.Permission{
		Roles:       in.Roles,
		Permissions: in.Permissions,
	}
}

func BankAccount(in *model.BankAccount) *identity.BankAccount {
	if in == nil {
		return nil
	}
	return &identity.BankAccount{
		Name:          in.Name,
		Province:      in.Province,
		Branch:        in.Branch,
		AccountNumber: in.AccountNumber,
		AccountName:   in.AccountName,
	}
}

func BankAccountDB(in *identity.BankAccount) *model.BankAccount {
	if in == nil {
		return nil
	}
	return &model.BankAccount{
		Name:          in.Name,
		Province:      in.Province,
		Branch:        in.Branch,
		AccountNumber: in.AccountNumber,
		AccountName:   in.AccountName,
	}
}
