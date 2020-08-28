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

func UserFtRefSaff(in *identitymodel.UserFtRefSaff) *identity.UserFtRefSaff {
	if in == nil {
		return nil
	}
	out := &identity.UserFtRefSaff{
		User: &identity.User{
			ID:                      in.ID,
			FullName:                in.FullName,
			ShortName:               in.ShortName,
			Email:                   in.Email,
			Phone:                   in.Phone,
			Status:                  in.Status,
			EmailVerifiedAt:         in.EmailVerifiedAt,
			PhoneVerifiedAt:         in.PhoneVerifiedAt,
			EmailVerificationSentAt: in.EmailVerificationSentAt,
			PhoneVerificationSentAt: in.PhoneVerificationSentAt,
			IsTest:                  in.IsTest,
			CreatedAt:               in.CreatedAt,
			UpdatedAt:               in.UpdatedAt,
			WLPartnerID:             in.WLPartnerID,
			Source:                  in.Source,
			BlockedAt:               in.BlockedAt,
			BlockedBy:               in.BlockedBy,
			BlockReason:             in.BlockReason,
		},
		RefSale: in.RefSale,
		RefAff:  in.RefAff,
	}
	return out
}

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

func AccountUser(in *identitymodel.AccountUser) *identity.AccountUser {
	if in == nil {
		return nil
	}
	out := &identity.AccountUser{}
	convert_identitymodel_AccountUser_identity_AccountUser(in, out)
	out.Permission.Roles = in.Roles
	out.Permission.Permissions = in.Permissions
	return out
}
