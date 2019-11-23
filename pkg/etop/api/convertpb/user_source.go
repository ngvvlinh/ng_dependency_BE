package convertpb

import (
	"etop.vn/api/pb/etop/etc/user_source"
	"etop.vn/backend/pkg/etop/model"
)

func PbUserSource(s model.UserSource) user_source.UserSource {
	_s := string(s)
	return user_source.UserSource(user_source.UserSource_value[_s])
}

func PbPtrUserSource(s model.UserSource) *user_source.UserSource {
	res := PbUserSource(s)
	return &res
}

func UserSourceToModel(s *user_source.UserSource) model.UserSource {
	if s == nil || *s == 0 {
		return ""
	}
	return model.UserSource(s.String())
}
