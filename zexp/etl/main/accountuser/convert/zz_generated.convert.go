// +build !generator

// Code generated by generator convert. DO NOT EDIT.

package convert

import (
	time "time"

	identitymodel "etop.vn/backend/com/main/identity/model"
	conversion "etop.vn/backend/pkg/common/conversion"
	accountusermodel "etop.vn/backend/zexp/etl/main/accountuser/model"
)

/*
Custom conversions:
    ConvertAccountUser    // in use

Ignored functions: (none)
*/

func RegisterConversions(s *conversion.Scheme) {
	registerConversions(s)
}

func registerConversions(s *conversion.Scheme) {
	s.Register((*accountusermodel.AccountUser)(nil), (*identitymodel.AccountUser)(nil), func(arg, out interface{}) error {
		Convert_accountusermodel_AccountUser_identitymodel_AccountUser(arg.(*accountusermodel.AccountUser), out.(*identitymodel.AccountUser))
		return nil
	})
	s.Register(([]*accountusermodel.AccountUser)(nil), (*[]*identitymodel.AccountUser)(nil), func(arg, out interface{}) error {
		out0 := Convert_accountusermodel_AccountUsers_identitymodel_AccountUsers(arg.([]*accountusermodel.AccountUser))
		*out.(*[]*identitymodel.AccountUser) = out0
		return nil
	})
	s.Register((*identitymodel.AccountUser)(nil), (*accountusermodel.AccountUser)(nil), func(arg, out interface{}) error {
		Convert_identitymodel_AccountUser_accountusermodel_AccountUser(arg.(*identitymodel.AccountUser), out.(*accountusermodel.AccountUser))
		return nil
	})
	s.Register(([]*identitymodel.AccountUser)(nil), (*[]*accountusermodel.AccountUser)(nil), func(arg, out interface{}) error {
		out0 := Convert_identitymodel_AccountUsers_accountusermodel_AccountUsers(arg.([]*identitymodel.AccountUser))
		*out.(*[]*accountusermodel.AccountUser) = out0
		return nil
	})
	s.Register((*accountusermodel.Permission)(nil), (*identitymodel.Permission)(nil), func(arg, out interface{}) error {
		Convert_accountusermodel_Permission_identitymodel_Permission(arg.(*accountusermodel.Permission), out.(*identitymodel.Permission))
		return nil
	})
	s.Register(([]*accountusermodel.Permission)(nil), (*[]*identitymodel.Permission)(nil), func(arg, out interface{}) error {
		out0 := Convert_accountusermodel_Permissions_identitymodel_Permissions(arg.([]*accountusermodel.Permission))
		*out.(*[]*identitymodel.Permission) = out0
		return nil
	})
	s.Register((*identitymodel.Permission)(nil), (*accountusermodel.Permission)(nil), func(arg, out interface{}) error {
		Convert_identitymodel_Permission_accountusermodel_Permission(arg.(*identitymodel.Permission), out.(*accountusermodel.Permission))
		return nil
	})
	s.Register(([]*identitymodel.Permission)(nil), (*[]*accountusermodel.Permission)(nil), func(arg, out interface{}) error {
		out0 := Convert_identitymodel_Permissions_accountusermodel_Permissions(arg.([]*identitymodel.Permission))
		*out.(*[]*accountusermodel.Permission) = out0
		return nil
	})
}

//-- convert etop.vn/backend/com/main/identity/model.AccountUser --//

func Convert_accountusermodel_AccountUser_identitymodel_AccountUser(arg *accountusermodel.AccountUser, out *identitymodel.AccountUser) *identitymodel.AccountUser {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &identitymodel.AccountUser{}
	}
	convert_accountusermodel_AccountUser_identitymodel_AccountUser(arg, out)
	return out
}

func convert_accountusermodel_AccountUser_identitymodel_AccountUser(arg *accountusermodel.AccountUser, out *identitymodel.AccountUser) {
	out.AccountID = arg.AccountID                       // simple assign
	out.UserID = arg.UserID                             // simple assign
	out.Status = arg.Status                             // simple assign
	out.ResponseStatus = arg.ResponseStatus             // simple assign
	out.CreatedAt = arg.CreatedAt                       // simple assign
	out.UpdatedAt = arg.UpdatedAt                       // simple assign
	out.DeletedAt = time.Time{}                         // zero value
	out.Permission = identitymodel.Permission{}         // types do not match
	out.FullName = arg.FullName                         // simple assign
	out.ShortName = arg.ShortName                       // simple assign
	out.Position = arg.Position                         // simple assign
	out.InvitationSentAt = arg.InvitationSentAt         // simple assign
	out.InvitationSentBy = arg.InvitationSentBy         // simple assign
	out.InvitationAcceptedAt = arg.InvitationAcceptedAt // simple assign
	out.InvitationRejectedAt = arg.InvitationRejectedAt // simple assign
	out.DisabledAt = arg.DisabledAt                     // simple assign
	out.DisabledBy = arg.DisabledBy                     // simple assign
	out.DisableReason = arg.DisableReason               // simple assign
	out.Rid = arg.Rid                                   // simple assign
}

func Convert_accountusermodel_AccountUsers_identitymodel_AccountUsers(args []*accountusermodel.AccountUser) (outs []*identitymodel.AccountUser) {
	if args == nil {
		return nil
	}
	tmps := make([]identitymodel.AccountUser, len(args))
	outs = make([]*identitymodel.AccountUser, len(args))
	for i := range tmps {
		outs[i] = Convert_accountusermodel_AccountUser_identitymodel_AccountUser(args[i], &tmps[i])
	}
	return outs
}

func Convert_identitymodel_AccountUser_accountusermodel_AccountUser(arg *identitymodel.AccountUser, out *accountusermodel.AccountUser) *accountusermodel.AccountUser {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &accountusermodel.AccountUser{}
	}
	ConvertAccountUser(arg, out)
	return out
}

func convert_identitymodel_AccountUser_accountusermodel_AccountUser(arg *identitymodel.AccountUser, out *accountusermodel.AccountUser) {
	out.AccountID = arg.AccountID                       // simple assign
	out.UserID = arg.UserID                             // simple assign
	out.Status = arg.Status                             // simple assign
	out.ResponseStatus = arg.ResponseStatus             // simple assign
	out.CreatedAt = arg.CreatedAt                       // simple assign
	out.UpdatedAt = arg.UpdatedAt                       // simple assign
	out.Permission = accountusermodel.Permission{}      // types do not match
	out.FullName = arg.FullName                         // simple assign
	out.ShortName = arg.ShortName                       // simple assign
	out.Position = arg.Position                         // simple assign
	out.InvitationSentAt = arg.InvitationSentAt         // simple assign
	out.InvitationSentBy = arg.InvitationSentBy         // simple assign
	out.InvitationAcceptedAt = arg.InvitationAcceptedAt // simple assign
	out.InvitationRejectedAt = arg.InvitationRejectedAt // simple assign
	out.DisabledAt = arg.DisabledAt                     // simple assign
	out.DisabledBy = arg.DisabledBy                     // simple assign
	out.DisableReason = arg.DisableReason               // simple assign
	out.Rid = arg.Rid                                   // simple assign
}

func Convert_identitymodel_AccountUsers_accountusermodel_AccountUsers(args []*identitymodel.AccountUser) (outs []*accountusermodel.AccountUser) {
	if args == nil {
		return nil
	}
	tmps := make([]accountusermodel.AccountUser, len(args))
	outs = make([]*accountusermodel.AccountUser, len(args))
	for i := range tmps {
		outs[i] = Convert_identitymodel_AccountUser_accountusermodel_AccountUser(args[i], &tmps[i])
	}
	return outs
}

//-- convert etop.vn/backend/com/main/identity/model.Permission --//

func Convert_accountusermodel_Permission_identitymodel_Permission(arg *accountusermodel.Permission, out *identitymodel.Permission) *identitymodel.Permission {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &identitymodel.Permission{}
	}
	convert_accountusermodel_Permission_identitymodel_Permission(arg, out)
	return out
}

func convert_accountusermodel_Permission_identitymodel_Permission(arg *accountusermodel.Permission, out *identitymodel.Permission) {
	out.Roles = arg.Roles             // simple assign
	out.Permissions = arg.Permissions // simple assign
}

func Convert_accountusermodel_Permissions_identitymodel_Permissions(args []*accountusermodel.Permission) (outs []*identitymodel.Permission) {
	if args == nil {
		return nil
	}
	tmps := make([]identitymodel.Permission, len(args))
	outs = make([]*identitymodel.Permission, len(args))
	for i := range tmps {
		outs[i] = Convert_accountusermodel_Permission_identitymodel_Permission(args[i], &tmps[i])
	}
	return outs
}

func Convert_identitymodel_Permission_accountusermodel_Permission(arg *identitymodel.Permission, out *accountusermodel.Permission) *accountusermodel.Permission {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &accountusermodel.Permission{}
	}
	convert_identitymodel_Permission_accountusermodel_Permission(arg, out)
	return out
}

func convert_identitymodel_Permission_accountusermodel_Permission(arg *identitymodel.Permission, out *accountusermodel.Permission) {
	out.Roles = arg.Roles             // simple assign
	out.Permissions = arg.Permissions // simple assign
}

func Convert_identitymodel_Permissions_accountusermodel_Permissions(args []*identitymodel.Permission) (outs []*accountusermodel.Permission) {
	if args == nil {
		return nil
	}
	tmps := make([]accountusermodel.Permission, len(args))
	outs = make([]*accountusermodel.Permission, len(args))
	for i := range tmps {
		outs[i] = Convert_identitymodel_Permission_accountusermodel_Permission(args[i], &tmps[i])
	}
	return outs
}
