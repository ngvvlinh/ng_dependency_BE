package convert

import (
	identitymodel "o.o/backend/com/main/identity/model"
	"o.o/backend/zexp/etl/main/accountuser/model"
)

// +gen:convert: o.o/backend/zexp/etl/main/accountuser/model -> o.o/backend/com/main/identity/model

func ConvertAccountUser(in *identitymodel.AccountUser, out *model.AccountUser) {
	convert_identitymodel_AccountUser_accountusermodel_AccountUser(in, out)
	for _, role := range in.Permission.Roles {
		out.Roles = append(out.Roles, role)
	}

	for _, permission := range in.Permission.Permissions {
		out.Permissions = append(out.Permissions, permission)
	}
}
