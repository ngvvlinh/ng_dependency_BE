package convert

import (
	identitymodel "o.o/backend/com/main/identity/model"
	"o.o/backend/zexp/etl/main/user/model"
)

// +gen:convert: o.o/backend/zexp/etl/main/user/model -> o.o/backend/com/main/identity/model

func ConvertUser(in *identitymodel.User, out *model.User) {
	convert_identitymodel_User_usermodel_User(in, out)
	out.FullName = in.FullName
	out.ShortName = in.ShortName
	out.Email = in.Email
	out.Phone = in.Phone
}
