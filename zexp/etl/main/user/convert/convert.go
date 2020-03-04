package convert

import (
	identitymodel "etop.vn/backend/com/main/identity/model"
	"etop.vn/backend/zexp/etl/main/user/model"
)

// +gen:convert: etop.vn/backend/zexp/etl/main/user/model->etop.vn/backend/com/main/identity/model

func ConvertUser(in *identitymodel.User, out *model.User) {
	convert_identitymodel_User_usermodel_User(in, out)
	out.FullName = in.FullName
	out.ShortName = in.ShortName
	out.Email = in.Email
	out.Phone = in.Phone
}
