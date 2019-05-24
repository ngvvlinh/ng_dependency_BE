package user_source

import "etop.vn/backend/pkg/etop/model"

func PbUserSource(s model.UserSource) UserSource {
	_s := string(s)
	return UserSource(UserSource_value[_s])
}

func PbPtrUserSource(s model.UserSource) *UserSource {
	res := PbUserSource(s)
	return &res
}

func (s *UserSource) ToModel() model.UserSource {
	if s == nil || *s == 0 {
		return ""
	}
	return model.UserSource(s.String())
}
