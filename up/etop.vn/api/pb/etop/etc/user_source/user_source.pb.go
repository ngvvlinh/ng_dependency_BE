package user_source

import (
	"etop.vn/common/jsonx"
)

type UserSource int

const (
	UserSource_unknown        UserSource = 0
	UserSource_psx            UserSource = 1
	UserSource_etop           UserSource = 2
	UserSource_topship        UserSource = 3
	UserSource_ts_app_android UserSource = 4
	UserSource_ts_app_ios     UserSource = 5
	UserSource_ts_app_web     UserSource = 6
	UserSource_partner        UserSource = 7
)

var UserSource_name = map[int]string{
	0: "unknown",
	1: "psx",
	2: "etop",
	3: "topship",
	4: "ts_app_android",
	5: "ts_app_ios",
	6: "ts_app_web",
	7: "partner",
}

var UserSource_value = map[string]int{
	"unknown":        0,
	"psx":            1,
	"etop":           2,
	"topship":        3,
	"ts_app_android": 4,
	"ts_app_ios":     5,
	"ts_app_web":     6,
	"partner":        7,
}

func (x UserSource) Enum() *UserSource {
	p := new(UserSource)
	*p = x
	return p
}

func (x UserSource) String() string {
	return jsonx.EnumName(UserSource_name, int(x))
}

func (x *UserSource) UnmarshalJSON(data []byte) error {
	value, err := jsonx.UnmarshalJSONEnum(UserSource_value, data, "UserSource")
	if err != nil {
		return err
	}
	*x = UserSource(value)
	return nil
}
