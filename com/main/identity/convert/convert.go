package convert

import (
	"o.o/api/main/identity"
	identitytypes "o.o/api/main/identity/types"
	"o.o/api/top/types/etc/status3"
	identitymodel "o.o/backend/com/main/identity/model"
	identitysharemodel "o.o/backend/com/main/identity/sharemodel"
)

// +gen:convert: o.o/backend/com/main/identity/sharemodel -> o.o/api/main/identity/types
// +gen:convert: o.o/backend/com/main/identity/model -> o.o/api/main/identity
// +gen:convert: o.o/api/main/identity
// +gen:convert: o.o/backend/com/main/identity/sharemodel -> o.o/api/top/int/types

func ShopDB(in *identity.Shop) *identitymodel.Shop {
	if in == nil {
		return nil
	}
	out := &identitymodel.Shop{}
	convert_identity_Shop_identitymodel_Shop(in, out)
	return out
}

func Shop(in *identitymodel.Shop) *identity.Shop {
	if in == nil {
		return nil
	}
	out := &identity.Shop{}
	convert_identitymodel_Shop_identity_Shop(in, out)
	return out
}

func User(in *identitymodel.User) (out *identity.User) {
	if in == nil {
		return nil
	}
	out = &identity.User{
		ID:              in.ID,
		FullName:        in.FullName,
		ShortName:       in.ShortName,
		Email:           in.Email,
		Phone:           in.Phone,
		Status:          in.Status,
		PhoneVerifiedAt: in.PhoneVerifiedAt,
		EmailVerifiedAt: in.EmailVerifiedAt,
		CreatedAt:       in.CreatedAt,
		UpdatedAt:       in.UpdatedAt,
		Source:          in.Source,
		BlockReason:     in.BlockReason,
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

func Permission(in identitymodel.Permission) identity.Permission {
	return identity.Permission{
		Roles:       in.Roles,
		Permissions: in.Permissions,
	}
}

func PermissionToModel(in identity.Permission) identitymodel.Permission {
	return identitymodel.Permission{
		Roles:       in.Roles,
		Permissions: in.Permissions,
	}
}

func BankAccount(in *identitysharemodel.BankAccount) *identitytypes.BankAccount {
	if in == nil {
		return nil
	}
	out := &identitytypes.BankAccount{}
	convert_sharemodel_BankAccount_identitytypes_BankAccount(in, out)
	return out
}

func BankAccountDB(in *identitytypes.BankAccount) *identitysharemodel.BankAccount {
	if in == nil {
		return nil
	}
	out := &identitysharemodel.BankAccount{}
	convert_identitytypes_BankAccount_sharemodel_BankAccount(in, out)
	return out
}
