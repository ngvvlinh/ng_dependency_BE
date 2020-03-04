// +build !generator

// Code generated by generator convert. DO NOT EDIT.

package convert

import (
	time "time"

	identity "etop.vn/api/main/identity"
	identitytypes "etop.vn/api/main/identity/types"
	inttypes "etop.vn/api/top/int/types"
	identitymodel "etop.vn/backend/com/main/identity/model"
	sharemodel "etop.vn/backend/com/main/identity/sharemodel"
	conversion "etop.vn/backend/pkg/common/conversion"
	dot "etop.vn/capi/dot"
)

/*
Custom conversions:
    Affiliate          // in use
    AffiliateDB        // in use
    BankAccount        // in use
    BankAccountDB      // in use
    Shop               // in use
    ShopDB             // in use
    User               // in use
    XAccountAhamove    // in use

Ignored functions:
    Affiliates           // params are not pointer to named types
    Permission           // params are not pointer to named types
    PermissionToModel    // params are not pointer to named types
*/

func RegisterConversions(s *conversion.Scheme) {
	registerConversions(s)
}

func registerConversions(s *conversion.Scheme) {
	s.Register((*identitymodel.Affiliate)(nil), (*identity.Affiliate)(nil), func(arg, out interface{}) error {
		Convert_identitymodel_Affiliate_identity_Affiliate(arg.(*identitymodel.Affiliate), out.(*identity.Affiliate))
		return nil
	})
	s.Register(([]*identitymodel.Affiliate)(nil), (*[]*identity.Affiliate)(nil), func(arg, out interface{}) error {
		out0 := Convert_identitymodel_Affiliates_identity_Affiliates(arg.([]*identitymodel.Affiliate))
		*out.(*[]*identity.Affiliate) = out0
		return nil
	})
	s.Register((*identity.Affiliate)(nil), (*identitymodel.Affiliate)(nil), func(arg, out interface{}) error {
		Convert_identity_Affiliate_identitymodel_Affiliate(arg.(*identity.Affiliate), out.(*identitymodel.Affiliate))
		return nil
	})
	s.Register(([]*identity.Affiliate)(nil), (*[]*identitymodel.Affiliate)(nil), func(arg, out interface{}) error {
		out0 := Convert_identity_Affiliates_identitymodel_Affiliates(arg.([]*identity.Affiliate))
		*out.(*[]*identitymodel.Affiliate) = out0
		return nil
	})
	s.Register((*identitymodel.ExternalAccountAhamove)(nil), (*identity.ExternalAccountAhamove)(nil), func(arg, out interface{}) error {
		Convert_identitymodel_ExternalAccountAhamove_identity_ExternalAccountAhamove(arg.(*identitymodel.ExternalAccountAhamove), out.(*identity.ExternalAccountAhamove))
		return nil
	})
	s.Register(([]*identitymodel.ExternalAccountAhamove)(nil), (*[]*identity.ExternalAccountAhamove)(nil), func(arg, out interface{}) error {
		out0 := Convert_identitymodel_ExternalAccountAhamoves_identity_ExternalAccountAhamoves(arg.([]*identitymodel.ExternalAccountAhamove))
		*out.(*[]*identity.ExternalAccountAhamove) = out0
		return nil
	})
	s.Register((*identity.ExternalAccountAhamove)(nil), (*identitymodel.ExternalAccountAhamove)(nil), func(arg, out interface{}) error {
		Convert_identity_ExternalAccountAhamove_identitymodel_ExternalAccountAhamove(arg.(*identity.ExternalAccountAhamove), out.(*identitymodel.ExternalAccountAhamove))
		return nil
	})
	s.Register(([]*identity.ExternalAccountAhamove)(nil), (*[]*identitymodel.ExternalAccountAhamove)(nil), func(arg, out interface{}) error {
		out0 := Convert_identity_ExternalAccountAhamoves_identitymodel_ExternalAccountAhamoves(arg.([]*identity.ExternalAccountAhamove))
		*out.(*[]*identitymodel.ExternalAccountAhamove) = out0
		return nil
	})
	s.Register((*identitymodel.Partner)(nil), (*identity.Partner)(nil), func(arg, out interface{}) error {
		Convert_identitymodel_Partner_identity_Partner(arg.(*identitymodel.Partner), out.(*identity.Partner))
		return nil
	})
	s.Register(([]*identitymodel.Partner)(nil), (*[]*identity.Partner)(nil), func(arg, out interface{}) error {
		out0 := Convert_identitymodel_Partners_identity_Partners(arg.([]*identitymodel.Partner))
		*out.(*[]*identity.Partner) = out0
		return nil
	})
	s.Register((*identity.Partner)(nil), (*identitymodel.Partner)(nil), func(arg, out interface{}) error {
		Convert_identity_Partner_identitymodel_Partner(arg.(*identity.Partner), out.(*identitymodel.Partner))
		return nil
	})
	s.Register(([]*identity.Partner)(nil), (*[]*identitymodel.Partner)(nil), func(arg, out interface{}) error {
		out0 := Convert_identity_Partners_identitymodel_Partners(arg.([]*identity.Partner))
		*out.(*[]*identitymodel.Partner) = out0
		return nil
	})
	s.Register((*identitymodel.Permission)(nil), (*identity.Permission)(nil), func(arg, out interface{}) error {
		Convert_identitymodel_Permission_identity_Permission(arg.(*identitymodel.Permission), out.(*identity.Permission))
		return nil
	})
	s.Register(([]*identitymodel.Permission)(nil), (*[]*identity.Permission)(nil), func(arg, out interface{}) error {
		out0 := Convert_identitymodel_Permissions_identity_Permissions(arg.([]*identitymodel.Permission))
		*out.(*[]*identity.Permission) = out0
		return nil
	})
	s.Register((*identity.Permission)(nil), (*identitymodel.Permission)(nil), func(arg, out interface{}) error {
		Convert_identity_Permission_identitymodel_Permission(arg.(*identity.Permission), out.(*identitymodel.Permission))
		return nil
	})
	s.Register(([]*identity.Permission)(nil), (*[]*identitymodel.Permission)(nil), func(arg, out interface{}) error {
		out0 := Convert_identity_Permissions_identitymodel_Permissions(arg.([]*identity.Permission))
		*out.(*[]*identitymodel.Permission) = out0
		return nil
	})
	s.Register((*identitymodel.Shop)(nil), (*identity.Shop)(nil), func(arg, out interface{}) error {
		Convert_identitymodel_Shop_identity_Shop(arg.(*identitymodel.Shop), out.(*identity.Shop))
		return nil
	})
	s.Register(([]*identitymodel.Shop)(nil), (*[]*identity.Shop)(nil), func(arg, out interface{}) error {
		out0 := Convert_identitymodel_Shops_identity_Shops(arg.([]*identitymodel.Shop))
		*out.(*[]*identity.Shop) = out0
		return nil
	})
	s.Register((*identity.Shop)(nil), (*identitymodel.Shop)(nil), func(arg, out interface{}) error {
		Convert_identity_Shop_identitymodel_Shop(arg.(*identity.Shop), out.(*identitymodel.Shop))
		return nil
	})
	s.Register(([]*identity.Shop)(nil), (*[]*identitymodel.Shop)(nil), func(arg, out interface{}) error {
		out0 := Convert_identity_Shops_identitymodel_Shops(arg.([]*identity.Shop))
		*out.(*[]*identitymodel.Shop) = out0
		return nil
	})
	s.Register((*identitymodel.User)(nil), (*identity.User)(nil), func(arg, out interface{}) error {
		Convert_identitymodel_User_identity_User(arg.(*identitymodel.User), out.(*identity.User))
		return nil
	})
	s.Register(([]*identitymodel.User)(nil), (*[]*identity.User)(nil), func(arg, out interface{}) error {
		out0 := Convert_identitymodel_Users_identity_Users(arg.([]*identitymodel.User))
		*out.(*[]*identity.User) = out0
		return nil
	})
	s.Register((*identity.User)(nil), (*identitymodel.User)(nil), func(arg, out interface{}) error {
		Convert_identity_User_identitymodel_User(arg.(*identity.User), out.(*identitymodel.User))
		return nil
	})
	s.Register(([]*identity.User)(nil), (*[]*identitymodel.User)(nil), func(arg, out interface{}) error {
		out0 := Convert_identity_Users_identitymodel_Users(arg.([]*identity.User))
		*out.(*[]*identitymodel.User) = out0
		return nil
	})
	s.Register((*sharemodel.BankAccount)(nil), (*identitytypes.BankAccount)(nil), func(arg, out interface{}) error {
		Convert_sharemodel_BankAccount_identitytypes_BankAccount(arg.(*sharemodel.BankAccount), out.(*identitytypes.BankAccount))
		return nil
	})
	s.Register(([]*sharemodel.BankAccount)(nil), (*[]*identitytypes.BankAccount)(nil), func(arg, out interface{}) error {
		out0 := Convert_sharemodel_BankAccounts_identitytypes_BankAccounts(arg.([]*sharemodel.BankAccount))
		*out.(*[]*identitytypes.BankAccount) = out0
		return nil
	})
	s.Register((*identitytypes.BankAccount)(nil), (*sharemodel.BankAccount)(nil), func(arg, out interface{}) error {
		Convert_identitytypes_BankAccount_sharemodel_BankAccount(arg.(*identitytypes.BankAccount), out.(*sharemodel.BankAccount))
		return nil
	})
	s.Register(([]*identitytypes.BankAccount)(nil), (*[]*sharemodel.BankAccount)(nil), func(arg, out interface{}) error {
		out0 := Convert_identitytypes_BankAccounts_sharemodel_BankAccounts(arg.([]*identitytypes.BankAccount))
		*out.(*[]*sharemodel.BankAccount) = out0
		return nil
	})
	s.Register((*sharemodel.AdjustmentLine)(nil), (*inttypes.AdjustmentLine)(nil), func(arg, out interface{}) error {
		Convert_sharemodel_AdjustmentLine_inttypes_AdjustmentLine(arg.(*sharemodel.AdjustmentLine), out.(*inttypes.AdjustmentLine))
		return nil
	})
	s.Register(([]*sharemodel.AdjustmentLine)(nil), (*[]*inttypes.AdjustmentLine)(nil), func(arg, out interface{}) error {
		out0 := Convert_sharemodel_AdjustmentLines_inttypes_AdjustmentLines(arg.([]*sharemodel.AdjustmentLine))
		*out.(*[]*inttypes.AdjustmentLine) = out0
		return nil
	})
	s.Register((*inttypes.AdjustmentLine)(nil), (*sharemodel.AdjustmentLine)(nil), func(arg, out interface{}) error {
		Convert_inttypes_AdjustmentLine_sharemodel_AdjustmentLine(arg.(*inttypes.AdjustmentLine), out.(*sharemodel.AdjustmentLine))
		return nil
	})
	s.Register(([]*inttypes.AdjustmentLine)(nil), (*[]*sharemodel.AdjustmentLine)(nil), func(arg, out interface{}) error {
		out0 := Convert_inttypes_AdjustmentLines_sharemodel_AdjustmentLines(arg.([]*inttypes.AdjustmentLine))
		*out.(*[]*sharemodel.AdjustmentLine) = out0
		return nil
	})
	s.Register((*sharemodel.DiscountLine)(nil), (*inttypes.DiscountLine)(nil), func(arg, out interface{}) error {
		Convert_sharemodel_DiscountLine_inttypes_DiscountLine(arg.(*sharemodel.DiscountLine), out.(*inttypes.DiscountLine))
		return nil
	})
	s.Register(([]*sharemodel.DiscountLine)(nil), (*[]*inttypes.DiscountLine)(nil), func(arg, out interface{}) error {
		out0 := Convert_sharemodel_DiscountLines_inttypes_DiscountLines(arg.([]*sharemodel.DiscountLine))
		*out.(*[]*inttypes.DiscountLine) = out0
		return nil
	})
	s.Register((*inttypes.DiscountLine)(nil), (*sharemodel.DiscountLine)(nil), func(arg, out interface{}) error {
		Convert_inttypes_DiscountLine_sharemodel_DiscountLine(arg.(*inttypes.DiscountLine), out.(*sharemodel.DiscountLine))
		return nil
	})
	s.Register(([]*inttypes.DiscountLine)(nil), (*[]*sharemodel.DiscountLine)(nil), func(arg, out interface{}) error {
		out0 := Convert_inttypes_DiscountLines_sharemodel_DiscountLines(arg.([]*inttypes.DiscountLine))
		*out.(*[]*sharemodel.DiscountLine) = out0
		return nil
	})
	s.Register((*sharemodel.FeeLine)(nil), (*inttypes.FeeLine)(nil), func(arg, out interface{}) error {
		Convert_sharemodel_FeeLine_inttypes_FeeLine(arg.(*sharemodel.FeeLine), out.(*inttypes.FeeLine))
		return nil
	})
	s.Register(([]*sharemodel.FeeLine)(nil), (*[]*inttypes.FeeLine)(nil), func(arg, out interface{}) error {
		out0 := Convert_sharemodel_FeeLines_inttypes_FeeLines(arg.([]*sharemodel.FeeLine))
		*out.(*[]*inttypes.FeeLine) = out0
		return nil
	})
	s.Register((*inttypes.FeeLine)(nil), (*sharemodel.FeeLine)(nil), func(arg, out interface{}) error {
		Convert_inttypes_FeeLine_sharemodel_FeeLine(arg.(*inttypes.FeeLine), out.(*sharemodel.FeeLine))
		return nil
	})
	s.Register(([]*inttypes.FeeLine)(nil), (*[]*sharemodel.FeeLine)(nil), func(arg, out interface{}) error {
		out0 := Convert_inttypes_FeeLines_sharemodel_FeeLines(arg.([]*inttypes.FeeLine))
		*out.(*[]*sharemodel.FeeLine) = out0
		return nil
	})
}

//-- convert etop.vn/api/main/identity.Affiliate --//

func Convert_identitymodel_Affiliate_identity_Affiliate(arg *identitymodel.Affiliate, out *identity.Affiliate) *identity.Affiliate {
	return Affiliate(arg)
}

func convert_identitymodel_Affiliate_identity_Affiliate(arg *identitymodel.Affiliate, out *identity.Affiliate) {
	out.ID = arg.ID               // simple assign
	out.OwnerID = arg.OwnerID     // simple assign
	out.Name = arg.Name           // simple assign
	out.Phone = arg.Phone         // simple assign
	out.Email = arg.Email         // simple assign
	out.IsTest = arg.IsTest       // simple assign
	out.Status = arg.Status       // simple assign
	out.CreatedAt = arg.CreatedAt // simple assign
	out.UpdatedAt = arg.UpdatedAt // simple assign
	out.DeletedAt = arg.DeletedAt // simple assign
	out.BankAccount = Convert_sharemodel_BankAccount_identitytypes_BankAccount(arg.BankAccount, nil)
}

func Convert_identitymodel_Affiliates_identity_Affiliates(args []*identitymodel.Affiliate) (outs []*identity.Affiliate) {
	tmps := make([]identity.Affiliate, len(args))
	outs = make([]*identity.Affiliate, len(args))
	for i := range tmps {
		outs[i] = Convert_identitymodel_Affiliate_identity_Affiliate(args[i], &tmps[i])
	}
	return outs
}

func Convert_identity_Affiliate_identitymodel_Affiliate(arg *identity.Affiliate, out *identitymodel.Affiliate) *identitymodel.Affiliate {
	return AffiliateDB(arg)
}

func convert_identity_Affiliate_identitymodel_Affiliate(arg *identity.Affiliate, out *identitymodel.Affiliate) {
	out.ID = arg.ID               // simple assign
	out.OwnerID = arg.OwnerID     // simple assign
	out.Name = arg.Name           // simple assign
	out.Phone = arg.Phone         // simple assign
	out.Email = arg.Email         // simple assign
	out.IsTest = arg.IsTest       // simple assign
	out.Status = arg.Status       // simple assign
	out.CreatedAt = arg.CreatedAt // simple assign
	out.UpdatedAt = arg.UpdatedAt // simple assign
	out.DeletedAt = arg.DeletedAt // simple assign
	out.BankAccount = Convert_identitytypes_BankAccount_sharemodel_BankAccount(arg.BankAccount, nil)
}

func Convert_identity_Affiliates_identitymodel_Affiliates(args []*identity.Affiliate) (outs []*identitymodel.Affiliate) {
	tmps := make([]identitymodel.Affiliate, len(args))
	outs = make([]*identitymodel.Affiliate, len(args))
	for i := range tmps {
		outs[i] = Convert_identity_Affiliate_identitymodel_Affiliate(args[i], &tmps[i])
	}
	return outs
}

//-- convert etop.vn/api/main/identity.ExternalAccountAhamove --//

func Convert_identitymodel_ExternalAccountAhamove_identity_ExternalAccountAhamove(arg *identitymodel.ExternalAccountAhamove, out *identity.ExternalAccountAhamove) *identity.ExternalAccountAhamove {
	return XAccountAhamove(arg)
}

func convert_identitymodel_ExternalAccountAhamove_identity_ExternalAccountAhamove(arg *identitymodel.ExternalAccountAhamove, out *identity.ExternalAccountAhamove) {
	out.ID = arg.ID                                   // simple assign
	out.Phone = arg.Phone                             // simple assign
	out.Name = arg.Name                               // simple assign
	out.ExternalID = arg.ExternalID                   // simple assign
	out.ExternalToken = arg.ExternalToken             // simple assign
	out.ExternalVerified = arg.ExternalVerified       // simple assign
	out.CreatedAt = arg.CreatedAt                     // simple assign
	out.UpdatedAt = arg.UpdatedAt                     // simple assign
	out.ExternalCreatedAt = arg.ExternalCreatedAt     // simple assign
	out.LastSendVerifiedAt = arg.LastSendVerifiedAt   // simple assign
	out.ExternalTicketID = arg.ExternalTicketID       // simple assign
	out.IDCardFrontImg = arg.IDCardFrontImg           // simple assign
	out.IDCardBackImg = arg.IDCardBackImg             // simple assign
	out.PortraitImg = arg.PortraitImg                 // simple assign
	out.WebsiteURL = arg.WebsiteURL                   // simple assign
	out.FanpageURL = arg.FanpageURL                   // simple assign
	out.CompanyImgs = arg.CompanyImgs                 // simple assign
	out.BusinessLicenseImgs = arg.BusinessLicenseImgs // simple assign
	out.UploadedAt = arg.UploadedAt                   // simple assign
}

func Convert_identitymodel_ExternalAccountAhamoves_identity_ExternalAccountAhamoves(args []*identitymodel.ExternalAccountAhamove) (outs []*identity.ExternalAccountAhamove) {
	tmps := make([]identity.ExternalAccountAhamove, len(args))
	outs = make([]*identity.ExternalAccountAhamove, len(args))
	for i := range tmps {
		outs[i] = Convert_identitymodel_ExternalAccountAhamove_identity_ExternalAccountAhamove(args[i], &tmps[i])
	}
	return outs
}

func Convert_identity_ExternalAccountAhamove_identitymodel_ExternalAccountAhamove(arg *identity.ExternalAccountAhamove, out *identitymodel.ExternalAccountAhamove) *identitymodel.ExternalAccountAhamove {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &identitymodel.ExternalAccountAhamove{}
	}
	convert_identity_ExternalAccountAhamove_identitymodel_ExternalAccountAhamove(arg, out)
	return out
}

func convert_identity_ExternalAccountAhamove_identitymodel_ExternalAccountAhamove(arg *identity.ExternalAccountAhamove, out *identitymodel.ExternalAccountAhamove) {
	out.ID = arg.ID                                   // simple assign
	out.OwnerID = 0                                   // zero value
	out.Phone = arg.Phone                             // simple assign
	out.Name = arg.Name                               // simple assign
	out.ExternalID = arg.ExternalID                   // simple assign
	out.ExternalVerified = arg.ExternalVerified       // simple assign
	out.ExternalCreatedAt = arg.ExternalCreatedAt     // simple assign
	out.ExternalToken = arg.ExternalToken             // simple assign
	out.CreatedAt = arg.CreatedAt                     // simple assign
	out.UpdatedAt = arg.UpdatedAt                     // simple assign
	out.LastSendVerifiedAt = arg.LastSendVerifiedAt   // simple assign
	out.ExternalTicketID = arg.ExternalTicketID       // simple assign
	out.IDCardFrontImg = arg.IDCardFrontImg           // simple assign
	out.IDCardBackImg = arg.IDCardBackImg             // simple assign
	out.PortraitImg = arg.PortraitImg                 // simple assign
	out.WebsiteURL = arg.WebsiteURL                   // simple assign
	out.FanpageURL = arg.FanpageURL                   // simple assign
	out.CompanyImgs = arg.CompanyImgs                 // simple assign
	out.BusinessLicenseImgs = arg.BusinessLicenseImgs // simple assign
	out.ExternalDataVerified = nil                    // zero value
	out.UploadedAt = arg.UploadedAt                   // simple assign
}

func Convert_identity_ExternalAccountAhamoves_identitymodel_ExternalAccountAhamoves(args []*identity.ExternalAccountAhamove) (outs []*identitymodel.ExternalAccountAhamove) {
	tmps := make([]identitymodel.ExternalAccountAhamove, len(args))
	outs = make([]*identitymodel.ExternalAccountAhamove, len(args))
	for i := range tmps {
		outs[i] = Convert_identity_ExternalAccountAhamove_identitymodel_ExternalAccountAhamove(args[i], &tmps[i])
	}
	return outs
}

//-- convert etop.vn/api/main/identity.Partner --//

func Convert_identitymodel_Partner_identity_Partner(arg *identitymodel.Partner, out *identity.Partner) *identity.Partner {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &identity.Partner{}
	}
	convert_identitymodel_Partner_identity_Partner(arg, out)
	return out
}

func convert_identitymodel_Partner_identity_Partner(arg *identitymodel.Partner, out *identity.Partner) {
	out.ID = arg.ID                       // simple assign
	out.OwnerID = arg.OwnerID             // simple assign
	out.Name = arg.Name                   // simple assign
	out.PublicName = arg.PublicName       // simple assign
	out.ImageURL = arg.ImageURL           // simple assign
	out.WebsiteURL = arg.WebsiteURL       // simple assign
	out.WhiteLabelKey = arg.WhiteLabelKey // simple assign
}

func Convert_identitymodel_Partners_identity_Partners(args []*identitymodel.Partner) (outs []*identity.Partner) {
	tmps := make([]identity.Partner, len(args))
	outs = make([]*identity.Partner, len(args))
	for i := range tmps {
		outs[i] = Convert_identitymodel_Partner_identity_Partner(args[i], &tmps[i])
	}
	return outs
}

func Convert_identity_Partner_identitymodel_Partner(arg *identity.Partner, out *identitymodel.Partner) *identitymodel.Partner {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &identitymodel.Partner{}
	}
	convert_identity_Partner_identitymodel_Partner(arg, out)
	return out
}

func convert_identity_Partner_identitymodel_Partner(arg *identity.Partner, out *identitymodel.Partner) {
	out.ID = arg.ID                       // simple assign
	out.OwnerID = arg.OwnerID             // simple assign
	out.Status = 0                        // zero value
	out.IsTest = 0                        // zero value
	out.Name = arg.Name                   // simple assign
	out.PublicName = arg.PublicName       // simple assign
	out.Phone = ""                        // zero value
	out.Email = ""                        // zero value
	out.ImageURL = arg.ImageURL           // simple assign
	out.WebsiteURL = arg.WebsiteURL       // simple assign
	out.ContactPersons = nil              // zero value
	out.RecognizedHosts = nil             // zero value
	out.RedirectURLs = nil                // zero value
	out.AvailableFromEtop = false         // zero value
	out.AvailableFromEtopConfig = nil     // zero value
	out.WhiteLabelKey = arg.WhiteLabelKey // simple assign
	out.CreatedAt = time.Time{}           // zero value
	out.UpdatedAt = time.Time{}           // zero value
	out.DeletedAt = time.Time{}           // zero value
}

func Convert_identity_Partners_identitymodel_Partners(args []*identity.Partner) (outs []*identitymodel.Partner) {
	tmps := make([]identitymodel.Partner, len(args))
	outs = make([]*identitymodel.Partner, len(args))
	for i := range tmps {
		outs[i] = Convert_identity_Partner_identitymodel_Partner(args[i], &tmps[i])
	}
	return outs
}

//-- convert etop.vn/api/main/identity.Permission --//

func Convert_identitymodel_Permission_identity_Permission(arg *identitymodel.Permission, out *identity.Permission) *identity.Permission {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &identity.Permission{}
	}
	convert_identitymodel_Permission_identity_Permission(arg, out)
	return out
}

func convert_identitymodel_Permission_identity_Permission(arg *identitymodel.Permission, out *identity.Permission) {
	out.Roles = arg.Roles             // simple assign
	out.Permissions = arg.Permissions // simple assign
}

func Convert_identitymodel_Permissions_identity_Permissions(args []*identitymodel.Permission) (outs []*identity.Permission) {
	tmps := make([]identity.Permission, len(args))
	outs = make([]*identity.Permission, len(args))
	for i := range tmps {
		outs[i] = Convert_identitymodel_Permission_identity_Permission(args[i], &tmps[i])
	}
	return outs
}

func Convert_identity_Permission_identitymodel_Permission(arg *identity.Permission, out *identitymodel.Permission) *identitymodel.Permission {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &identitymodel.Permission{}
	}
	convert_identity_Permission_identitymodel_Permission(arg, out)
	return out
}

func convert_identity_Permission_identitymodel_Permission(arg *identity.Permission, out *identitymodel.Permission) {
	out.Roles = arg.Roles             // simple assign
	out.Permissions = arg.Permissions // simple assign
}

func Convert_identity_Permissions_identitymodel_Permissions(args []*identity.Permission) (outs []*identitymodel.Permission) {
	tmps := make([]identitymodel.Permission, len(args))
	outs = make([]*identitymodel.Permission, len(args))
	for i := range tmps {
		outs[i] = Convert_identity_Permission_identitymodel_Permission(args[i], &tmps[i])
	}
	return outs
}

//-- convert etop.vn/api/main/identity.Shop --//

func Convert_identitymodel_Shop_identity_Shop(arg *identitymodel.Shop, out *identity.Shop) *identity.Shop {
	return Shop(arg)
}

func convert_identitymodel_Shop_identity_Shop(arg *identitymodel.Shop, out *identity.Shop) {
	out.ID = arg.ID                               // simple assign
	out.Name = arg.Name                           // simple assign
	out.OwnerID = arg.OwnerID                     // simple assign
	out.IsTest = arg.IsTest                       // simple assign
	out.AddressID = arg.AddressID                 // simple assign
	out.ShipToAddressID = arg.ShipToAddressID     // simple assign
	out.ShipFromAddressID = arg.ShipFromAddressID // simple assign
	out.Phone = arg.Phone                         // simple assign
	out.WebsiteURL = arg.WebsiteURL               // simple assign
	out.ImageURL = arg.ImageURL                   // simple assign
	out.Email = arg.Email                         // simple assign
	out.Code = arg.Code                           // simple assign
	out.AutoCreateFFM = arg.AutoCreateFFM         // simple assign
	out.Status = arg.Status                       // simple assign
	out.CreatedAt = arg.CreatedAt                 // simple assign
	out.UpdatedAt = arg.UpdatedAt                 // simple assign
	out.DeletedAt = arg.DeletedAt                 // simple assign
	out.BankAccount = Convert_sharemodel_BankAccount_identitytypes_BankAccount(arg.BankAccount, nil)
	out.TryOn = arg.TryOn // simple assign
}

func Convert_identitymodel_Shops_identity_Shops(args []*identitymodel.Shop) (outs []*identity.Shop) {
	tmps := make([]identity.Shop, len(args))
	outs = make([]*identity.Shop, len(args))
	for i := range tmps {
		outs[i] = Convert_identitymodel_Shop_identity_Shop(args[i], &tmps[i])
	}
	return outs
}

func Convert_identity_Shop_identitymodel_Shop(arg *identity.Shop, out *identitymodel.Shop) *identitymodel.Shop {
	return ShopDB(arg)
}

func convert_identity_Shop_identitymodel_Shop(arg *identity.Shop, out *identitymodel.Shop) {
	out.ID = arg.ID                               // simple assign
	out.Name = arg.Name                           // simple assign
	out.OwnerID = arg.OwnerID                     // simple assign
	out.IsTest = arg.IsTest                       // simple assign
	out.AddressID = arg.AddressID                 // simple assign
	out.ShipToAddressID = arg.ShipToAddressID     // simple assign
	out.ShipFromAddressID = arg.ShipFromAddressID // simple assign
	out.Phone = arg.Phone                         // simple assign
	out.BankAccount = Convert_identitytypes_BankAccount_sharemodel_BankAccount(arg.BankAccount, nil)
	out.WebsiteURL = arg.WebsiteURL         // simple assign
	out.ImageURL = arg.ImageURL             // simple assign
	out.Email = arg.Email                   // simple assign
	out.Code = arg.Code                     // simple assign
	out.AutoCreateFFM = arg.AutoCreateFFM   // simple assign
	out.OrderSourceID = 0                   // zero value
	out.Status = arg.Status                 // simple assign
	out.CreatedAt = arg.CreatedAt           // simple assign
	out.UpdatedAt = arg.UpdatedAt           // simple assign
	out.DeletedAt = arg.DeletedAt           // simple assign
	out.Address = nil                       // zero value
	out.RecognizedHosts = nil               // zero value
	out.GhnNoteCode = 0                     // zero value
	out.TryOn = arg.TryOn                   // simple assign
	out.CompanyInfo = nil                   // zero value
	out.MoneyTransactionRRule = ""          // zero value
	out.SurveyInfo = nil                    // zero value
	out.ShippingServiceSelectStrategy = nil // zero value
	out.InventoryOverstock = dot.NullBool{} // zero value
	out.Rid = 0                             // zero value
}

func Convert_identity_Shops_identitymodel_Shops(args []*identity.Shop) (outs []*identitymodel.Shop) {
	tmps := make([]identitymodel.Shop, len(args))
	outs = make([]*identitymodel.Shop, len(args))
	for i := range tmps {
		outs[i] = Convert_identity_Shop_identitymodel_Shop(args[i], &tmps[i])
	}
	return outs
}

//-- convert etop.vn/api/main/identity.User --//

func Convert_identitymodel_User_identity_User(arg *identitymodel.User, out *identity.User) *identity.User {
	return User(arg)
}

func convert_identitymodel_User_identity_User(arg *identitymodel.User, out *identity.User) {
	out.ID = arg.ID                           // simple assign
	out.FullName = ""                         // zero value
	out.ShortName = ""                        // zero value
	out.Email = ""                            // zero value
	out.Phone = ""                            // zero value
	out.Status = arg.Status                   // simple assign
	out.EmailVerifiedAt = arg.EmailVerifiedAt // simple assign
	out.PhoneVerifiedAt = arg.PhoneVerifiedAt // simple assign
	out.CreatedAt = arg.CreatedAt             // simple assign
	out.UpdatedAt = arg.UpdatedAt             // simple assign
	out.RefUserID = arg.RefUserID             // simple assign
	out.RefSaleID = arg.RefSaleID             // simple assign
	out.WLPartnerID = arg.WLPartnerID         // simple assign
}

func Convert_identitymodel_Users_identity_Users(args []*identitymodel.User) (outs []*identity.User) {
	tmps := make([]identity.User, len(args))
	outs = make([]*identity.User, len(args))
	for i := range tmps {
		outs[i] = Convert_identitymodel_User_identity_User(args[i], &tmps[i])
	}
	return outs
}

func Convert_identity_User_identitymodel_User(arg *identity.User, out *identitymodel.User) *identitymodel.User {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &identitymodel.User{}
	}
	convert_identity_User_identitymodel_User(arg, out)
	return out
}

func convert_identity_User_identitymodel_User(arg *identity.User, out *identitymodel.User) {
	out.ID = arg.ID                           // simple assign
	out.UserInner = identitymodel.UserInner{} // zero value
	out.Status = arg.Status                   // simple assign
	out.CreatedAt = arg.CreatedAt             // simple assign
	out.UpdatedAt = arg.UpdatedAt             // simple assign
	out.AgreedTOSAt = time.Time{}             // zero value
	out.AgreedEmailInfoAt = time.Time{}       // zero value
	out.EmailVerifiedAt = arg.EmailVerifiedAt // simple assign
	out.PhoneVerifiedAt = arg.PhoneVerifiedAt // simple assign
	out.EmailVerificationSentAt = time.Time{} // zero value
	out.PhoneVerificationSentAt = time.Time{} // zero value
	out.IsTest = 0                            // zero value
	out.Source = 0                            // zero value
	out.RefUserID = arg.RefUserID             // simple assign
	out.RefSaleID = arg.RefSaleID             // simple assign
	out.WLPartnerID = arg.WLPartnerID         // simple assign
	out.Rid = 0                               // zero value
}

func Convert_identity_Users_identitymodel_Users(args []*identity.User) (outs []*identitymodel.User) {
	tmps := make([]identitymodel.User, len(args))
	outs = make([]*identitymodel.User, len(args))
	for i := range tmps {
		outs[i] = Convert_identity_User_identitymodel_User(args[i], &tmps[i])
	}
	return outs
}

//-- convert etop.vn/api/main/identity/types.BankAccount --//

func Convert_sharemodel_BankAccount_identitytypes_BankAccount(arg *sharemodel.BankAccount, out *identitytypes.BankAccount) *identitytypes.BankAccount {
	return BankAccount(arg)
}

func convert_sharemodel_BankAccount_identitytypes_BankAccount(arg *sharemodel.BankAccount, out *identitytypes.BankAccount) {
	out.Name = arg.Name                   // simple assign
	out.Province = arg.Province           // simple assign
	out.Branch = arg.Branch               // simple assign
	out.AccountNumber = arg.AccountNumber // simple assign
	out.AccountName = arg.AccountName     // simple assign
}

func Convert_sharemodel_BankAccounts_identitytypes_BankAccounts(args []*sharemodel.BankAccount) (outs []*identitytypes.BankAccount) {
	tmps := make([]identitytypes.BankAccount, len(args))
	outs = make([]*identitytypes.BankAccount, len(args))
	for i := range tmps {
		outs[i] = Convert_sharemodel_BankAccount_identitytypes_BankAccount(args[i], &tmps[i])
	}
	return outs
}

func Convert_identitytypes_BankAccount_sharemodel_BankAccount(arg *identitytypes.BankAccount, out *sharemodel.BankAccount) *sharemodel.BankAccount {
	return BankAccountDB(arg)
}

func convert_identitytypes_BankAccount_sharemodel_BankAccount(arg *identitytypes.BankAccount, out *sharemodel.BankAccount) {
	out.Name = arg.Name                   // simple assign
	out.Province = arg.Province           // simple assign
	out.Branch = arg.Branch               // simple assign
	out.AccountNumber = arg.AccountNumber // simple assign
	out.AccountName = arg.AccountName     // simple assign
}

func Convert_identitytypes_BankAccounts_sharemodel_BankAccounts(args []*identitytypes.BankAccount) (outs []*sharemodel.BankAccount) {
	tmps := make([]sharemodel.BankAccount, len(args))
	outs = make([]*sharemodel.BankAccount, len(args))
	for i := range tmps {
		outs[i] = Convert_identitytypes_BankAccount_sharemodel_BankAccount(args[i], &tmps[i])
	}
	return outs
}

//-- convert etop.vn/api/top/int/types.AdjustmentLine --//

func Convert_sharemodel_AdjustmentLine_inttypes_AdjustmentLine(arg *sharemodel.AdjustmentLine, out *inttypes.AdjustmentLine) *inttypes.AdjustmentLine {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &inttypes.AdjustmentLine{}
	}
	convert_sharemodel_AdjustmentLine_inttypes_AdjustmentLine(arg, out)
	return out
}

func convert_sharemodel_AdjustmentLine_inttypes_AdjustmentLine(arg *sharemodel.AdjustmentLine, out *inttypes.AdjustmentLine) {
	out.Note = arg.Note     // simple assign
	out.Amount = arg.Amount // simple assign
}

func Convert_sharemodel_AdjustmentLines_inttypes_AdjustmentLines(args []*sharemodel.AdjustmentLine) (outs []*inttypes.AdjustmentLine) {
	tmps := make([]inttypes.AdjustmentLine, len(args))
	outs = make([]*inttypes.AdjustmentLine, len(args))
	for i := range tmps {
		outs[i] = Convert_sharemodel_AdjustmentLine_inttypes_AdjustmentLine(args[i], &tmps[i])
	}
	return outs
}

func Convert_inttypes_AdjustmentLine_sharemodel_AdjustmentLine(arg *inttypes.AdjustmentLine, out *sharemodel.AdjustmentLine) *sharemodel.AdjustmentLine {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &sharemodel.AdjustmentLine{}
	}
	convert_inttypes_AdjustmentLine_sharemodel_AdjustmentLine(arg, out)
	return out
}

func convert_inttypes_AdjustmentLine_sharemodel_AdjustmentLine(arg *inttypes.AdjustmentLine, out *sharemodel.AdjustmentLine) {
	out.Note = arg.Note     // simple assign
	out.Amount = arg.Amount // simple assign
}

func Convert_inttypes_AdjustmentLines_sharemodel_AdjustmentLines(args []*inttypes.AdjustmentLine) (outs []*sharemodel.AdjustmentLine) {
	tmps := make([]sharemodel.AdjustmentLine, len(args))
	outs = make([]*sharemodel.AdjustmentLine, len(args))
	for i := range tmps {
		outs[i] = Convert_inttypes_AdjustmentLine_sharemodel_AdjustmentLine(args[i], &tmps[i])
	}
	return outs
}

//-- convert etop.vn/api/top/int/types.DiscountLine --//

func Convert_sharemodel_DiscountLine_inttypes_DiscountLine(arg *sharemodel.DiscountLine, out *inttypes.DiscountLine) *inttypes.DiscountLine {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &inttypes.DiscountLine{}
	}
	convert_sharemodel_DiscountLine_inttypes_DiscountLine(arg, out)
	return out
}

func convert_sharemodel_DiscountLine_inttypes_DiscountLine(arg *sharemodel.DiscountLine, out *inttypes.DiscountLine) {
	out.Note = arg.Note     // simple assign
	out.Amount = arg.Amount // simple assign
}

func Convert_sharemodel_DiscountLines_inttypes_DiscountLines(args []*sharemodel.DiscountLine) (outs []*inttypes.DiscountLine) {
	tmps := make([]inttypes.DiscountLine, len(args))
	outs = make([]*inttypes.DiscountLine, len(args))
	for i := range tmps {
		outs[i] = Convert_sharemodel_DiscountLine_inttypes_DiscountLine(args[i], &tmps[i])
	}
	return outs
}

func Convert_inttypes_DiscountLine_sharemodel_DiscountLine(arg *inttypes.DiscountLine, out *sharemodel.DiscountLine) *sharemodel.DiscountLine {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &sharemodel.DiscountLine{}
	}
	convert_inttypes_DiscountLine_sharemodel_DiscountLine(arg, out)
	return out
}

func convert_inttypes_DiscountLine_sharemodel_DiscountLine(arg *inttypes.DiscountLine, out *sharemodel.DiscountLine) {
	out.Note = arg.Note     // simple assign
	out.Amount = arg.Amount // simple assign
}

func Convert_inttypes_DiscountLines_sharemodel_DiscountLines(args []*inttypes.DiscountLine) (outs []*sharemodel.DiscountLine) {
	tmps := make([]sharemodel.DiscountLine, len(args))
	outs = make([]*sharemodel.DiscountLine, len(args))
	for i := range tmps {
		outs[i] = Convert_inttypes_DiscountLine_sharemodel_DiscountLine(args[i], &tmps[i])
	}
	return outs
}

//-- convert etop.vn/api/top/int/types.FeeLine --//

func Convert_sharemodel_FeeLine_inttypes_FeeLine(arg *sharemodel.FeeLine, out *inttypes.FeeLine) *inttypes.FeeLine {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &inttypes.FeeLine{}
	}
	convert_sharemodel_FeeLine_inttypes_FeeLine(arg, out)
	return out
}

func convert_sharemodel_FeeLine_inttypes_FeeLine(arg *sharemodel.FeeLine, out *inttypes.FeeLine) {
	out.Note = arg.Note     // simple assign
	out.Amount = arg.Amount // simple assign
}

func Convert_sharemodel_FeeLines_inttypes_FeeLines(args []*sharemodel.FeeLine) (outs []*inttypes.FeeLine) {
	tmps := make([]inttypes.FeeLine, len(args))
	outs = make([]*inttypes.FeeLine, len(args))
	for i := range tmps {
		outs[i] = Convert_sharemodel_FeeLine_inttypes_FeeLine(args[i], &tmps[i])
	}
	return outs
}

func Convert_inttypes_FeeLine_sharemodel_FeeLine(arg *inttypes.FeeLine, out *sharemodel.FeeLine) *sharemodel.FeeLine {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &sharemodel.FeeLine{}
	}
	convert_inttypes_FeeLine_sharemodel_FeeLine(arg, out)
	return out
}

func convert_inttypes_FeeLine_sharemodel_FeeLine(arg *inttypes.FeeLine, out *sharemodel.FeeLine) {
	out.Note = arg.Note     // simple assign
	out.Amount = arg.Amount // simple assign
}

func Convert_inttypes_FeeLines_sharemodel_FeeLines(args []*inttypes.FeeLine) (outs []*sharemodel.FeeLine) {
	tmps := make([]sharemodel.FeeLine, len(args))
	outs = make([]*sharemodel.FeeLine, len(args))
	for i := range tmps {
		outs[i] = Convert_inttypes_FeeLine_sharemodel_FeeLine(args[i], &tmps[i])
	}
	return outs
}
