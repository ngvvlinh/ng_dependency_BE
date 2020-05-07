// +build !generator

// Code generated by generator convert. DO NOT EDIT.

package convert

import (
	time "time"

	fbusering "o.o/api/fabo/fbusering"
	fbusermodel "o.o/backend/com/fabo/main/fbuser/model"
	conversion "o.o/backend/pkg/common/conversion"
)

/*
Custom conversions: (none)

Ignored functions: (none)
*/

func RegisterConversions(s *conversion.Scheme) {
	registerConversions(s)
}

func registerConversions(s *conversion.Scheme) {
	s.Register((*fbusermodel.FbExternalUser)(nil), (*fbusering.FbExternalUser)(nil), func(arg, out interface{}) error {
		Convert_fbusermodel_FbExternalUser_fbusering_FbExternalUser(arg.(*fbusermodel.FbExternalUser), out.(*fbusering.FbExternalUser))
		return nil
	})
	s.Register(([]*fbusermodel.FbExternalUser)(nil), (*[]*fbusering.FbExternalUser)(nil), func(arg, out interface{}) error {
		out0 := Convert_fbusermodel_FbExternalUsers_fbusering_FbExternalUsers(arg.([]*fbusermodel.FbExternalUser))
		*out.(*[]*fbusering.FbExternalUser) = out0
		return nil
	})
	s.Register((*fbusering.FbExternalUser)(nil), (*fbusermodel.FbExternalUser)(nil), func(arg, out interface{}) error {
		Convert_fbusering_FbExternalUser_fbusermodel_FbExternalUser(arg.(*fbusering.FbExternalUser), out.(*fbusermodel.FbExternalUser))
		return nil
	})
	s.Register(([]*fbusering.FbExternalUser)(nil), (*[]*fbusermodel.FbExternalUser)(nil), func(arg, out interface{}) error {
		out0 := Convert_fbusering_FbExternalUsers_fbusermodel_FbExternalUsers(arg.([]*fbusering.FbExternalUser))
		*out.(*[]*fbusermodel.FbExternalUser) = out0
		return nil
	})
	s.Register((*fbusering.CreateFbExternalUserArgs)(nil), (*fbusering.FbExternalUser)(nil), func(arg, out interface{}) error {
		Apply_fbusering_CreateFbExternalUserArgs_fbusering_FbExternalUser(arg.(*fbusering.CreateFbExternalUserArgs), out.(*fbusering.FbExternalUser))
		return nil
	})
	s.Register((*fbusermodel.FbExternalUserInfo)(nil), (*fbusering.FbExternalUserInfo)(nil), func(arg, out interface{}) error {
		Convert_fbusermodel_FbExternalUserInfo_fbusering_FbExternalUserInfo(arg.(*fbusermodel.FbExternalUserInfo), out.(*fbusering.FbExternalUserInfo))
		return nil
	})
	s.Register(([]*fbusermodel.FbExternalUserInfo)(nil), (*[]*fbusering.FbExternalUserInfo)(nil), func(arg, out interface{}) error {
		out0 := Convert_fbusermodel_FbExternalUserInfoes_fbusering_FbExternalUserInfoes(arg.([]*fbusermodel.FbExternalUserInfo))
		*out.(*[]*fbusering.FbExternalUserInfo) = out0
		return nil
	})
	s.Register((*fbusering.FbExternalUserInfo)(nil), (*fbusermodel.FbExternalUserInfo)(nil), func(arg, out interface{}) error {
		Convert_fbusering_FbExternalUserInfo_fbusermodel_FbExternalUserInfo(arg.(*fbusering.FbExternalUserInfo), out.(*fbusermodel.FbExternalUserInfo))
		return nil
	})
	s.Register(([]*fbusering.FbExternalUserInfo)(nil), (*[]*fbusermodel.FbExternalUserInfo)(nil), func(arg, out interface{}) error {
		out0 := Convert_fbusering_FbExternalUserInfoes_fbusermodel_FbExternalUserInfoes(arg.([]*fbusering.FbExternalUserInfo))
		*out.(*[]*fbusermodel.FbExternalUserInfo) = out0
		return nil
	})
	s.Register((*fbusermodel.FbExternalUserInternal)(nil), (*fbusering.FbExternalUserInternal)(nil), func(arg, out interface{}) error {
		Convert_fbusermodel_FbExternalUserInternal_fbusering_FbExternalUserInternal(arg.(*fbusermodel.FbExternalUserInternal), out.(*fbusering.FbExternalUserInternal))
		return nil
	})
	s.Register(([]*fbusermodel.FbExternalUserInternal)(nil), (*[]*fbusering.FbExternalUserInternal)(nil), func(arg, out interface{}) error {
		out0 := Convert_fbusermodel_FbExternalUserInternals_fbusering_FbExternalUserInternals(arg.([]*fbusermodel.FbExternalUserInternal))
		*out.(*[]*fbusering.FbExternalUserInternal) = out0
		return nil
	})
	s.Register((*fbusering.FbExternalUserInternal)(nil), (*fbusermodel.FbExternalUserInternal)(nil), func(arg, out interface{}) error {
		Convert_fbusering_FbExternalUserInternal_fbusermodel_FbExternalUserInternal(arg.(*fbusering.FbExternalUserInternal), out.(*fbusermodel.FbExternalUserInternal))
		return nil
	})
	s.Register(([]*fbusering.FbExternalUserInternal)(nil), (*[]*fbusermodel.FbExternalUserInternal)(nil), func(arg, out interface{}) error {
		out0 := Convert_fbusering_FbExternalUserInternals_fbusermodel_FbExternalUserInternals(arg.([]*fbusering.FbExternalUserInternal))
		*out.(*[]*fbusermodel.FbExternalUserInternal) = out0
		return nil
	})
	s.Register((*fbusering.CreateFbExternalUserInternalArgs)(nil), (*fbusering.FbExternalUserInternal)(nil), func(arg, out interface{}) error {
		Apply_fbusering_CreateFbExternalUserInternalArgs_fbusering_FbExternalUserInternal(arg.(*fbusering.CreateFbExternalUserInternalArgs), out.(*fbusering.FbExternalUserInternal))
		return nil
	})
	s.Register((*fbusermodel.FbExternalUserShopCustomer)(nil), (*fbusering.FbExternalUserShopCustomer)(nil), func(arg, out interface{}) error {
		Convert_fbusermodel_FbExternalUserShopCustomer_fbusering_FbExternalUserShopCustomer(arg.(*fbusermodel.FbExternalUserShopCustomer), out.(*fbusering.FbExternalUserShopCustomer))
		return nil
	})
	s.Register(([]*fbusermodel.FbExternalUserShopCustomer)(nil), (*[]*fbusering.FbExternalUserShopCustomer)(nil), func(arg, out interface{}) error {
		out0 := Convert_fbusermodel_FbExternalUserShopCustomers_fbusering_FbExternalUserShopCustomers(arg.([]*fbusermodel.FbExternalUserShopCustomer))
		*out.(*[]*fbusering.FbExternalUserShopCustomer) = out0
		return nil
	})
	s.Register((*fbusering.FbExternalUserShopCustomer)(nil), (*fbusermodel.FbExternalUserShopCustomer)(nil), func(arg, out interface{}) error {
		Convert_fbusering_FbExternalUserShopCustomer_fbusermodel_FbExternalUserShopCustomer(arg.(*fbusering.FbExternalUserShopCustomer), out.(*fbusermodel.FbExternalUserShopCustomer))
		return nil
	})
	s.Register(([]*fbusering.FbExternalUserShopCustomer)(nil), (*[]*fbusermodel.FbExternalUserShopCustomer)(nil), func(arg, out interface{}) error {
		out0 := Convert_fbusering_FbExternalUserShopCustomers_fbusermodel_FbExternalUserShopCustomers(arg.([]*fbusering.FbExternalUserShopCustomer))
		*out.(*[]*fbusermodel.FbExternalUserShopCustomer) = out0
		return nil
	})
}

//-- convert o.o/api/fabo/fbusering.FbExternalUser --//

func Convert_fbusermodel_FbExternalUser_fbusering_FbExternalUser(arg *fbusermodel.FbExternalUser, out *fbusering.FbExternalUser) *fbusering.FbExternalUser {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &fbusering.FbExternalUser{}
	}
	convert_fbusermodel_FbExternalUser_fbusering_FbExternalUser(arg, out)
	return out
}

func convert_fbusermodel_FbExternalUser_fbusering_FbExternalUser(arg *fbusermodel.FbExternalUser, out *fbusering.FbExternalUser) {
	out.ExternalID = arg.ExternalID // simple assign
	out.ShopID = 0                  // zero value
	out.ExternalInfo = Convert_fbusermodel_FbExternalUserInfo_fbusering_FbExternalUserInfo(arg.ExternalInfo, nil)
	out.Status = arg.Status       // simple assign
	out.CreatedAt = arg.CreatedAt // simple assign
	out.UpdatedAt = arg.UpdatedAt // simple assign
}

func Convert_fbusermodel_FbExternalUsers_fbusering_FbExternalUsers(args []*fbusermodel.FbExternalUser) (outs []*fbusering.FbExternalUser) {
	if args == nil {
		return nil
	}
	tmps := make([]fbusering.FbExternalUser, len(args))
	outs = make([]*fbusering.FbExternalUser, len(args))
	for i := range tmps {
		outs[i] = Convert_fbusermodel_FbExternalUser_fbusering_FbExternalUser(args[i], &tmps[i])
	}
	return outs
}

func Convert_fbusering_FbExternalUser_fbusermodel_FbExternalUser(arg *fbusering.FbExternalUser, out *fbusermodel.FbExternalUser) *fbusermodel.FbExternalUser {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &fbusermodel.FbExternalUser{}
	}
	convert_fbusering_FbExternalUser_fbusermodel_FbExternalUser(arg, out)
	return out
}

func convert_fbusering_FbExternalUser_fbusermodel_FbExternalUser(arg *fbusering.FbExternalUser, out *fbusermodel.FbExternalUser) {
	out.ExternalID = arg.ExternalID // simple assign
	out.ExternalInfo = Convert_fbusering_FbExternalUserInfo_fbusermodel_FbExternalUserInfo(arg.ExternalInfo, nil)
	out.Status = arg.Status       // simple assign
	out.CreatedAt = arg.CreatedAt // simple assign
	out.UpdatedAt = arg.UpdatedAt // simple assign
}

func Convert_fbusering_FbExternalUsers_fbusermodel_FbExternalUsers(args []*fbusering.FbExternalUser) (outs []*fbusermodel.FbExternalUser) {
	if args == nil {
		return nil
	}
	tmps := make([]fbusermodel.FbExternalUser, len(args))
	outs = make([]*fbusermodel.FbExternalUser, len(args))
	for i := range tmps {
		outs[i] = Convert_fbusering_FbExternalUser_fbusermodel_FbExternalUser(args[i], &tmps[i])
	}
	return outs
}

func Apply_fbusering_CreateFbExternalUserArgs_fbusering_FbExternalUser(arg *fbusering.CreateFbExternalUserArgs, out *fbusering.FbExternalUser) *fbusering.FbExternalUser {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &fbusering.FbExternalUser{}
	}
	apply_fbusering_CreateFbExternalUserArgs_fbusering_FbExternalUser(arg, out)
	return out
}

func apply_fbusering_CreateFbExternalUserArgs_fbusering_FbExternalUser(arg *fbusering.CreateFbExternalUserArgs, out *fbusering.FbExternalUser) {
	out.ExternalID = arg.ExternalID     // simple assign
	out.ShopID = 0                      // zero value
	out.ExternalInfo = arg.ExternalInfo // simple assign
	out.Status = arg.Status             // simple assign
	out.CreatedAt = time.Time{}         // zero value
	out.UpdatedAt = time.Time{}         // zero value
}

//-- convert o.o/api/fabo/fbusering.FbExternalUserInfo --//

func Convert_fbusermodel_FbExternalUserInfo_fbusering_FbExternalUserInfo(arg *fbusermodel.FbExternalUserInfo, out *fbusering.FbExternalUserInfo) *fbusering.FbExternalUserInfo {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &fbusering.FbExternalUserInfo{}
	}
	convert_fbusermodel_FbExternalUserInfo_fbusering_FbExternalUserInfo(arg, out)
	return out
}

func convert_fbusermodel_FbExternalUserInfo_fbusering_FbExternalUserInfo(arg *fbusermodel.FbExternalUserInfo, out *fbusering.FbExternalUserInfo) {
	out.Name = arg.Name           // simple assign
	out.FirstName = arg.FirstName // simple assign
	out.LastName = arg.LastName   // simple assign
	out.ShortName = arg.ShortName // simple assign
	out.ImageURL = arg.ImageURL   // simple assign
}

func Convert_fbusermodel_FbExternalUserInfoes_fbusering_FbExternalUserInfoes(args []*fbusermodel.FbExternalUserInfo) (outs []*fbusering.FbExternalUserInfo) {
	if args == nil {
		return nil
	}
	tmps := make([]fbusering.FbExternalUserInfo, len(args))
	outs = make([]*fbusering.FbExternalUserInfo, len(args))
	for i := range tmps {
		outs[i] = Convert_fbusermodel_FbExternalUserInfo_fbusering_FbExternalUserInfo(args[i], &tmps[i])
	}
	return outs
}

func Convert_fbusering_FbExternalUserInfo_fbusermodel_FbExternalUserInfo(arg *fbusering.FbExternalUserInfo, out *fbusermodel.FbExternalUserInfo) *fbusermodel.FbExternalUserInfo {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &fbusermodel.FbExternalUserInfo{}
	}
	convert_fbusering_FbExternalUserInfo_fbusermodel_FbExternalUserInfo(arg, out)
	return out
}

func convert_fbusering_FbExternalUserInfo_fbusermodel_FbExternalUserInfo(arg *fbusering.FbExternalUserInfo, out *fbusermodel.FbExternalUserInfo) {
	out.Name = arg.Name           // simple assign
	out.FirstName = arg.FirstName // simple assign
	out.LastName = arg.LastName   // simple assign
	out.ShortName = arg.ShortName // simple assign
	out.ImageURL = arg.ImageURL   // simple assign
}

func Convert_fbusering_FbExternalUserInfoes_fbusermodel_FbExternalUserInfoes(args []*fbusering.FbExternalUserInfo) (outs []*fbusermodel.FbExternalUserInfo) {
	if args == nil {
		return nil
	}
	tmps := make([]fbusermodel.FbExternalUserInfo, len(args))
	outs = make([]*fbusermodel.FbExternalUserInfo, len(args))
	for i := range tmps {
		outs[i] = Convert_fbusering_FbExternalUserInfo_fbusermodel_FbExternalUserInfo(args[i], &tmps[i])
	}
	return outs
}

//-- convert o.o/api/fabo/fbusering.FbExternalUserInternal --//

func Convert_fbusermodel_FbExternalUserInternal_fbusering_FbExternalUserInternal(arg *fbusermodel.FbExternalUserInternal, out *fbusering.FbExternalUserInternal) *fbusering.FbExternalUserInternal {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &fbusering.FbExternalUserInternal{}
	}
	convert_fbusermodel_FbExternalUserInternal_fbusering_FbExternalUserInternal(arg, out)
	return out
}

func convert_fbusermodel_FbExternalUserInternal_fbusering_FbExternalUserInternal(arg *fbusermodel.FbExternalUserInternal, out *fbusering.FbExternalUserInternal) {
	out.ExternalID = arg.ExternalID // simple assign
	out.Token = arg.Token           // simple assign
	out.ExpiresIn = arg.ExpiresIn   // simple assign
	out.UpdatedAt = arg.UpdatedAt   // simple assign
}

func Convert_fbusermodel_FbExternalUserInternals_fbusering_FbExternalUserInternals(args []*fbusermodel.FbExternalUserInternal) (outs []*fbusering.FbExternalUserInternal) {
	if args == nil {
		return nil
	}
	tmps := make([]fbusering.FbExternalUserInternal, len(args))
	outs = make([]*fbusering.FbExternalUserInternal, len(args))
	for i := range tmps {
		outs[i] = Convert_fbusermodel_FbExternalUserInternal_fbusering_FbExternalUserInternal(args[i], &tmps[i])
	}
	return outs
}

func Convert_fbusering_FbExternalUserInternal_fbusermodel_FbExternalUserInternal(arg *fbusering.FbExternalUserInternal, out *fbusermodel.FbExternalUserInternal) *fbusermodel.FbExternalUserInternal {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &fbusermodel.FbExternalUserInternal{}
	}
	convert_fbusering_FbExternalUserInternal_fbusermodel_FbExternalUserInternal(arg, out)
	return out
}

func convert_fbusering_FbExternalUserInternal_fbusermodel_FbExternalUserInternal(arg *fbusering.FbExternalUserInternal, out *fbusermodel.FbExternalUserInternal) {
	out.ExternalID = arg.ExternalID // simple assign
	out.Token = arg.Token           // simple assign
	out.ExpiresIn = arg.ExpiresIn   // simple assign
	out.UpdatedAt = arg.UpdatedAt   // simple assign
}

func Convert_fbusering_FbExternalUserInternals_fbusermodel_FbExternalUserInternals(args []*fbusering.FbExternalUserInternal) (outs []*fbusermodel.FbExternalUserInternal) {
	if args == nil {
		return nil
	}
	tmps := make([]fbusermodel.FbExternalUserInternal, len(args))
	outs = make([]*fbusermodel.FbExternalUserInternal, len(args))
	for i := range tmps {
		outs[i] = Convert_fbusering_FbExternalUserInternal_fbusermodel_FbExternalUserInternal(args[i], &tmps[i])
	}
	return outs
}

func Apply_fbusering_CreateFbExternalUserInternalArgs_fbusering_FbExternalUserInternal(arg *fbusering.CreateFbExternalUserInternalArgs, out *fbusering.FbExternalUserInternal) *fbusering.FbExternalUserInternal {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &fbusering.FbExternalUserInternal{}
	}
	apply_fbusering_CreateFbExternalUserInternalArgs_fbusering_FbExternalUserInternal(arg, out)
	return out
}

func apply_fbusering_CreateFbExternalUserInternalArgs_fbusering_FbExternalUserInternal(arg *fbusering.CreateFbExternalUserInternalArgs, out *fbusering.FbExternalUserInternal) {
	out.ExternalID = arg.ExternalID // simple assign
	out.Token = arg.Token           // simple assign
	out.ExpiresIn = arg.ExpiresIn   // simple assign
	out.UpdatedAt = time.Time{}     // zero value
}

//-- convert o.o/api/fabo/fbusering.FbExternalUserShopCustomer --//

func Convert_fbusermodel_FbExternalUserShopCustomer_fbusering_FbExternalUserShopCustomer(arg *fbusermodel.FbExternalUserShopCustomer, out *fbusering.FbExternalUserShopCustomer) *fbusering.FbExternalUserShopCustomer {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &fbusering.FbExternalUserShopCustomer{}
	}
	convert_fbusermodel_FbExternalUserShopCustomer_fbusering_FbExternalUserShopCustomer(arg, out)
	return out
}

func convert_fbusermodel_FbExternalUserShopCustomer_fbusering_FbExternalUserShopCustomer(arg *fbusermodel.FbExternalUserShopCustomer, out *fbusering.FbExternalUserShopCustomer) {
	out.CreatedAt = arg.CreatedAt               // simple assign
	out.UpdatedAt = arg.UpdatedAt               // simple assign
	out.ShopID = arg.ShopID                     // simple assign
	out.FbExternalUserID = arg.FbExternalUserID // simple assign
	out.CustomerID = arg.CustomerID             // simple assign
	out.Status = arg.Status                     // simple assign
}

func Convert_fbusermodel_FbExternalUserShopCustomers_fbusering_FbExternalUserShopCustomers(args []*fbusermodel.FbExternalUserShopCustomer) (outs []*fbusering.FbExternalUserShopCustomer) {
	if args == nil {
		return nil
	}
	tmps := make([]fbusering.FbExternalUserShopCustomer, len(args))
	outs = make([]*fbusering.FbExternalUserShopCustomer, len(args))
	for i := range tmps {
		outs[i] = Convert_fbusermodel_FbExternalUserShopCustomer_fbusering_FbExternalUserShopCustomer(args[i], &tmps[i])
	}
	return outs
}

func Convert_fbusering_FbExternalUserShopCustomer_fbusermodel_FbExternalUserShopCustomer(arg *fbusering.FbExternalUserShopCustomer, out *fbusermodel.FbExternalUserShopCustomer) *fbusermodel.FbExternalUserShopCustomer {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &fbusermodel.FbExternalUserShopCustomer{}
	}
	convert_fbusering_FbExternalUserShopCustomer_fbusermodel_FbExternalUserShopCustomer(arg, out)
	return out
}

func convert_fbusering_FbExternalUserShopCustomer_fbusermodel_FbExternalUserShopCustomer(arg *fbusering.FbExternalUserShopCustomer, out *fbusermodel.FbExternalUserShopCustomer) {
	out.CreatedAt = arg.CreatedAt               // simple assign
	out.UpdatedAt = arg.UpdatedAt               // simple assign
	out.ShopID = arg.ShopID                     // simple assign
	out.FbExternalUserID = arg.FbExternalUserID // simple assign
	out.CustomerID = arg.CustomerID             // simple assign
	out.Status = arg.Status                     // simple assign
}

func Convert_fbusering_FbExternalUserShopCustomers_fbusermodel_FbExternalUserShopCustomers(args []*fbusering.FbExternalUserShopCustomer) (outs []*fbusermodel.FbExternalUserShopCustomer) {
	if args == nil {
		return nil
	}
	tmps := make([]fbusermodel.FbExternalUserShopCustomer, len(args))
	outs = make([]*fbusermodel.FbExternalUserShopCustomer, len(args))
	for i := range tmps {
		outs[i] = Convert_fbusering_FbExternalUserShopCustomer_fbusermodel_FbExternalUserShopCustomer(args[i], &tmps[i])
	}
	return outs
}
