package convertpb

import (
	"o.o/api/etelecom/usersetting"
	etelecomtypes "o.o/api/top/int/etelecom/types"
)

func Convert_usersetting_UserSetting_api_UserSetting(in *usersetting.UserSetting) *etelecomtypes.EtelecomUserSetting {
	if in == nil {
		return nil
	}
	return &etelecomtypes.EtelecomUserSetting{
		ID:                  in.ID,
		ExtensionChargeType: in.ExtensionChargeType,
		CreatedAt:           in.CreatedAt,
		UpdatedAt:           in.UpdatedAt,
	}
}

func Convert_usersetting_UserSettings_api_UserSettings(ins []*usersetting.UserSetting) (out []*etelecomtypes.EtelecomUserSetting) {
	if ins == nil {
		return nil
	}
	for _, in := range ins {
		out = append(out, Convert_usersetting_UserSetting_api_UserSetting(in))
	}
	return
}
